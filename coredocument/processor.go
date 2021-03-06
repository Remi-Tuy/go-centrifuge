package coredocument

import (
	"context"
	"time"

	"github.com/centrifuge/go-centrifuge/contextutil"

	"github.com/centrifuge/centrifuge-protobufs/gen/go/coredocument"
	"github.com/centrifuge/centrifuge-protobufs/gen/go/p2p"
	"github.com/centrifuge/go-centrifuge/anchors"
	"github.com/centrifuge/go-centrifuge/crypto/secp256k1"
	"github.com/centrifuge/go-centrifuge/documents"
	"github.com/centrifuge/go-centrifuge/errors"
	"github.com/centrifuge/go-centrifuge/identity"
	"github.com/centrifuge/go-centrifuge/utils"
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
	Send(ctx context.Context, coreDocument *coredocumentpb.CoreDocument, recipient identity.CentID) (err error)
	PrepareForSignatureRequests(ctx context.Context, model documents.Model) error
	RequestSignatures(ctx context.Context, model documents.Model) error
	PrepareForAnchoring(model documents.Model) error
	AnchorDocument(ctx context.Context, model documents.Model) error
	SendDocument(ctx context.Context, model documents.Model) error
}

// client defines the methods for p2pclient
// we redefined it here so that we can avoid cyclic dependencies with p2p
type client interface {
	GetSignaturesForDocument(ctx context.Context, identityService identity.Service, doc *coredocumentpb.CoreDocument) error
	SendAnchoredDocument(ctx context.Context, id identity.Identity, in *p2ppb.AnchorDocumentRequest) (*p2ppb.AnchorDocumentResponse, error)
}

// defaultProcessor implements Processor interface
type defaultProcessor struct {
	identityService  identity.Service
	p2pClient        client
	anchorRepository anchors.AnchorRepository
	// TODO [multi-tenancy] replace this with config service
	config Config
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
func (dp defaultProcessor) Send(ctx context.Context, coreDocument *coredocumentpb.CoreDocument, recipient identity.CentID) (err error) {
	if coreDocument == nil {
		return errors.New("passed coreDoc is nil")
	}
	log.Infof("sending coredocument %x to recipient %x", coreDocument.DocumentIdentifier, recipient)
	id, err := dp.identityService.LookupIdentityForID(recipient)
	if err != nil {
		return errors.New("error fetching receiver identity: %v", err)
	}

	c, _ := context.WithTimeout(ctx, dp.config.GetP2PConnectionTimeout())
	resp, err := dp.p2pClient.SendAnchoredDocument(c, id, &p2ppb.AnchorDocumentRequest{Document: coreDocument})
	if err != nil || !resp.Accepted {
		return errors.New("failed to send document to the node: %v", err)
	}

	log.Infof("Done opening connection against recipient [%x]\n", recipient)

	return nil
}

// PrepareForSignatureRequests gets the core document from the model, and adds the node's own signature
func (dp defaultProcessor) PrepareForSignatureRequests(ctx context.Context, model documents.Model) error {
	cd, err := model.PackCoreDocument()
	if err != nil {
		return errors.New("failed to pack core document: %v", err)
	}

	// calculate the signing root
	err = CalculateSigningRoot(cd)
	if err != nil {
		return errors.New("failed to calculate signing root: %v", err)
	}

	self, err := contextutil.Self(ctx)
	if err != nil {
		return err
	}

	sig := identity.Sign(self, identity.KeyPurposeSigning, cd.SigningRoot)
	cd.Signatures = append(cd.Signatures, sig)

	err = model.UnpackCoreDocument(cd)
	if err != nil {
		return errors.New("failed to unpack the core document: %v", err)
	}

	return nil
}

// RequestSignatures gets the core document from the model, validates pre signature requirements,
// collects signatures, and validates the signatures,
func (dp defaultProcessor) RequestSignatures(ctx context.Context, model documents.Model) error {
	cd, err := model.PackCoreDocument()
	if err != nil {
		return errors.New("failed to pack core document: %v", err)
	}

	self, err := contextutil.Self(ctx)
	if err != nil {
		return err
	}

	idKeys, ok := self.Keys[identity.KeyPurposeSigning]
	if !ok {
		return errors.New("missing keys for signing")
	}

	psv := PreSignatureRequestValidator(self.ID[:], idKeys.PrivateKey, idKeys.PublicKey)
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
func (dp defaultProcessor) AnchorDocument(ctx context.Context, model documents.Model) error {
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

	self, err := contextutil.Self(ctx)
	if err != nil {
		return err
	}

	// generate message authentication code for the anchor call
	mac, err := secp256k1.SignEthereum(anchors.GenerateCommitHash(anchorID, centID, rootHash), self.Keys[identity.KeyPurposeEthMsgAuth].PrivateKey)
	if err != nil {
		return errors.New("failed to generate ethereum MAC: %v", err)
	}

	log.Infof("Anchoring document with identifiers: [document: %#x, current: %#x, next: %#x], rootHash: %#x", cd.DocumentIdentifier, cd.CurrentVersion, cd.NextVersion, cd.DocumentRoot)
	confirmations, err := dp.anchorRepository.CommitAnchor(ctx, anchorID, rootHash, centID, [][anchors.DocumentProofLength]byte{utils.RandomByte32()}, mac)
	if err != nil {
		return errors.New("failed to commit anchor: %v", err)
	}

	<-confirmations
	log.Infof("Anchored document with identifiers: [document: %#x, current: %#x, next: %#x], rootHash: %#x", cd.DocumentIdentifier, cd.CurrentVersion, cd.NextVersion, cd.DocumentRoot)
	return nil
}

// SendDocument does post anchor validations and sends the document to collaborators
func (dp defaultProcessor) SendDocument(ctx context.Context, model documents.Model) error {
	cd, err := model.PackCoreDocument()
	if err != nil {
		return errors.New("failed to pack core document: %v", err)
	}

	av := PostAnchoredValidator(dp.identityService, dp.anchorRepository)
	err = av.Validate(nil, model)
	if err != nil {
		return errors.New("post anchor validations failed: %v", err)
	}

	self, err := contextutil.Self(ctx)
	if err != nil {
		return err
	}

	extCollaborators, err := GetExternalCollaborators(self.ID, cd)
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
