package coredocument

import (
	"context"
	"time"

	context2 "github.com/centrifuge/go-centrifuge/context"

	"github.com/centrifuge/centrifuge-protobufs/gen/go/coredocument"
	"github.com/centrifuge/centrifuge-protobufs/gen/go/p2p"
	"github.com/centrifuge/go-centrifuge/anchors"
	"github.com/centrifuge/go-centrifuge/centerrors"
	"github.com/centrifuge/go-centrifuge/crypto/secp256k1"
	"github.com/centrifuge/go-centrifuge/documents"
	"github.com/centrifuge/go-centrifuge/errors"
	"github.com/centrifuge/go-centrifuge/identity"
	"github.com/centrifuge/go-centrifuge/utils"
	"github.com/centrifuge/go-centrifuge/version"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("coredocument")

// Config defines required methods required for the coredocument package.
type Config interface {
	GetNetworkID() uint32
	GetIdentityID() ([]byte, error)
	GetP2PConnectionTimeout() time.Duration
}

// Processor identifies an implementation, which can do a bunch of things with a CoreDocument.
// E.g. send, anchor, etc.
type Processor interface {
	Send(ctx *context2.Header, coreDocument *coredocumentpb.CoreDocument, recipient identity.CentID) (err error)
	PrepareForSignatureRequests(ctx *context2.Header, model documents.Model) error
	RequestSignatures(ctx *context2.Header, model documents.Model) error
	PrepareForAnchoring(model documents.Model) error
	AnchorDocument(ctx *context2.Header, model documents.Model) error
	SendDocument(ctx *context2.Header, model documents.Model) error
}

// client defines the methods for p2pclient
// we redefined it here so that we can avoid cyclic dependencies with p2p
type client interface {
	OpenClient(id identity.Identity) (p2ppb.P2PServiceClient, error)
	GetSignaturesForDocument(ctx *context2.Header, identityService identity.Service, doc *coredocumentpb.CoreDocument) error
}

// defaultProcessor implements Processor interface
type defaultProcessor struct {
	identityService  identity.Service
	p2pClient        client
	anchorRepository anchors.AnchorRepository
	config           Config
}

// DefaultProcessor returns the default implementation of CoreDocument Processor
func DefaultProcessor(idService identity.Service, p2pClient client, repository anchors.AnchorRepository, config Config) Processor {
	return defaultProcessor{
		identityService:  idService,
		p2pClient:        p2pClient,
		anchorRepository: repository,
		config:           config,
	}
}

// Send sends the given defaultProcessor to the given recipient on the P2P layer
func (dp defaultProcessor) Send(ctx *context2.Header, coreDocument *coredocumentpb.CoreDocument, recipient identity.CentID) (err error) {
	if coreDocument == nil {
		return centerrors.NilError(coreDocument)
	}
	log.Infof("sending coredocument %x to recipient %x", coreDocument.DocumentIdentifier, recipient)
	id, err := dp.identityService.LookupIdentityForID(recipient)
	if err != nil {
		return centerrors.Wrap(err, "error fetching receiver identity")
	}
	client, err := dp.p2pClient.OpenClient(id)
	if err != nil {
		return errors.New("failed to open client: %v", err)
	}

	log.Infof("Done opening connection against recipient [%x]\n", recipient)
	idConfig := ctx.Self()
	centIDBytes := idConfig.ID[:]
	p2pheader := &p2ppb.CentrifugeHeader{
		SenderCentrifugeId: centIDBytes,
		CentNodeVersion:    version.GetVersion().String(),
		NetworkIdentifier:  dp.config.GetNetworkID(),
	}

	c, _ := context.WithTimeout(ctx.Context(), dp.config.GetP2PConnectionTimeout())
	resp, err := client.SendAnchoredDocument(c, &p2ppb.AnchorDocumentRequest{Document: coreDocument, Header: p2pheader})
	if err != nil || !resp.Accepted {
		return centerrors.Wrap(err, "failed to send document to the node")
	}

	return nil
}

// PrepareForSignatureRequests gets the core document from the model, and adds the node's own signature
func (dp defaultProcessor) PrepareForSignatureRequests(ctx *context2.Header, model documents.Model) error {
	cd, err := model.PackCoreDocument()
	if err != nil {
		return errors.New("failed to pack core document: %v", err)
	}

	// calculate the signing root
	err = CalculateSigningRoot(cd)
	if err != nil {
		return errors.New("failed to calculate signing root: %v", err)
	}

	sig := identity.Sign(ctx.Self(), identity.KeyPurposeSigning, cd.SigningRoot)
	cd.Signatures = append(cd.Signatures, sig)

	err = model.UnpackCoreDocument(cd)
	if err != nil {
		return errors.New("failed to unpack the core document: %v", err)
	}

	return nil
}

// RequestSignatures gets the core document from the model, validates pre signature requirements,
// collects signatures, and validates the signatures,
func (dp defaultProcessor) RequestSignatures(ctx *context2.Header, model documents.Model) error {
	cd, err := model.PackCoreDocument()
	if err != nil {
		return errors.New("failed to pack core document: %v", err)
	}

	idKeys, ok := ctx.Self().Keys[identity.KeyPurposeSigning]
	if !ok {
		return errors.New("missing keys for signing")
	}

	psv := PreSignatureRequestValidator(ctx.Self().ID[:], idKeys.PrivateKey, idKeys.PublicKey)
	err = psv.Validate(nil, model)
	if err != nil {
		return errors.New("failed to validate model for signature request: %v", err)
	}

	err = dp.p2pClient.GetSignaturesForDocument(ctx, dp.identityService, cd)
	if err != nil {
		return errors.New("failed to collect signatures from the collaborators: %v", err)
	}

	err = model.UnpackCoreDocument(cd)
	if err != nil {
		return errors.New("failed to unpack core document: %v", err)
	}

	return nil
}

// PrepareForAnchoring validates the signatures and generates the document root
func (dp defaultProcessor) PrepareForAnchoring(model documents.Model) error {
	cd, err := model.PackCoreDocument()
	if err != nil {
		return errors.New("failed to pack core document: %v", err)
	}

	psv := PostSignatureRequestValidator(dp.identityService)
	err = psv.Validate(nil, model)
	if err != nil {
		return errors.New("failed to validate signatures: %v", err)
	}

	err = CalculateDocumentRoot(cd)
	if err != nil {
		return errors.New("failed to generate document root: %v", err)
	}

	err = model.UnpackCoreDocument(cd)
	if err != nil {
		return errors.New("failed to unpack core document: %v", err)
	}

	return nil
}

// AnchorDocument validates the model, and anchors the document
func (dp defaultProcessor) AnchorDocument(ctx *context2.Header, model documents.Model) error {
	cd, err := model.PackCoreDocument()
	if err != nil {
		return errors.New("failed to pack core document: %v", err)
	}

	pav := PreAnchorValidator(dp.identityService)
	err = pav.Validate(nil, model)
	if err != nil {
		return errors.New("pre anchor validation failed: %v", err)
	}

	rootHash, err := anchors.ToDocumentRoot(cd.DocumentRoot)
	if err != nil {
		return errors.New("failed to get document root: %v", err)
	}

	id, err := dp.config.GetIdentityID()
	if err != nil {
		return errors.New("failed to get self cent ID: %v", err)
	}

	centID, err := identity.ToCentID(id)
	if err != nil {
		return errors.New("centID invalid: %v", err)
	}

	anchorID, err := anchors.ToAnchorID(cd.CurrentVersion)
	if err != nil {
		return errors.New("failed to get anchor ID: %v", err)
	}

	// generate message authentication code for the anchor call
	mac, err := secp256k1.SignEthereum(anchors.GenerateCommitHash(anchorID, centID, rootHash), ctx.Self().Keys[identity.KeyPurposeEthMsgAuth].PrivateKey)
	if err != nil {
		return errors.New("failed to generate ethereum MAC: %v", err)
	}

	log.Infof("Anchoring document with identifiers: [document: %#x, current: %#x, next: %#x], rootHash: %#x", cd.DocumentIdentifier, cd.CurrentVersion, cd.NextVersion, cd.DocumentRoot)
	confirmations, err := dp.anchorRepository.CommitAnchor(anchorID, rootHash, centID, [][anchors.DocumentProofLength]byte{utils.RandomByte32()}, mac)
	if err != nil {
		return errors.New("failed to commit anchor: %v", err)
	}

	<-confirmations
	log.Infof("Anchored document with identifiers: [document: %#x, current: %#x, next: %#x], rootHash: %#x", cd.DocumentIdentifier, cd.CurrentVersion, cd.NextVersion, cd.DocumentRoot)
	return nil
}

// SendDocument does post anchor validations and sends the document to collaborators
func (dp defaultProcessor) SendDocument(ctx *context2.Header, model documents.Model) error {
	cd, err := model.PackCoreDocument()
	if err != nil {
		return errors.New("failed to pack core document: %v", err)
	}

	av := PostAnchoredValidator(dp.identityService, dp.anchorRepository)
	err = av.Validate(nil, model)
	if err != nil {
		return errors.New("post anchor validations failed: %v", err)
	}

	extCollaborators, err := GetExternalCollaborators(ctx.Self().ID, cd)
	if err != nil {
		return errors.New("get external collaborators failed: %v", err)
	}

	for _, c := range extCollaborators {
		cID, erri := identity.ToCentID(c)
		if erri != nil {
			err = errors.AppendError(err, erri)
			continue
		}

		erri = dp.Send(ctx, cd, cID)
		if erri != nil {
			err = errors.AppendError(err, erri)
		}
	}

	return err
}
