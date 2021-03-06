package purchaseorder

import (
	"github.com/centrifuge/go-centrifuge/config"
	"github.com/centrifuge/go-centrifuge/contextutil"

	"github.com/centrifuge/centrifuge-protobufs/documenttypes"
	"github.com/centrifuge/go-centrifuge/centerrors"
	"github.com/centrifuge/go-centrifuge/documents"
	"github.com/centrifuge/go-centrifuge/errors"
	clientpurchaseorderpb "github.com/centrifuge/go-centrifuge/protobufs/gen/go/purchaseorder"
	"github.com/ethereum/go-ethereum/common/hexutil"
	logging "github.com/ipfs/go-log"
	"golang.org/x/net/context"
)

var apiLog = logging.Logger("purchaseorder-api")

// grpcHandler handles all the purchase order document related actions
// anchoring, sending, finding stored purchase order document
type grpcHandler struct {
	service Service
	config  config.Service
}

// GRPCHandler returns an implementation of the purchaseorder DocumentServiceServer
func GRPCHandler(config config.Service, registry *documents.ServiceRegistry) (clientpurchaseorderpb.DocumentServiceServer, error) {
	srv, err := registry.LocateService(documenttypes.PurchaseOrderDataTypeUrl)
	if err != nil {
		return nil, errors.New("failed to fetch purchase order service")
	}

	return grpcHandler{
		service: srv.(Service),
		config:  config,
	}, nil
}

// Create validates the purchase order, persists it to DB, and anchors it the chain
func (h grpcHandler) Create(ctx context.Context, req *clientpurchaseorderpb.PurchaseOrderCreatePayload) (*clientpurchaseorderpb.PurchaseOrderResponse, error) {
	apiLog.Debugf("Create request %v", req)
	ctxh, err := contextutil.CentContext(ctx, h.config)
	if err != nil {
		apiLog.Error(err)
		return nil, err
	}

	doc, err := h.service.DeriveFromCreatePayload(ctxh, req)
	if err != nil {
		apiLog.Error(err)
		return nil, centerrors.Wrap(err, "could not derive create payload")
	}

	// validate, persist, and anchor
	doc, txID, err := h.service.Create(ctxh, doc)
	if err != nil {
		apiLog.Error(err)
		return nil, centerrors.Wrap(err, "could not create document")
	}

	resp, err := h.service.DerivePurchaseOrderResponse(doc)
	if err != nil {
		apiLog.Error(err)
		return nil, centerrors.Wrap(err, "could not derive response")
	}

	resp.Header.TransactionId = txID.String()
	return resp, nil
}

// Update handles the document update and anchoring
func (h grpcHandler) Update(ctx context.Context, payload *clientpurchaseorderpb.PurchaseOrderUpdatePayload) (*clientpurchaseorderpb.PurchaseOrderResponse, error) {
	apiLog.Debugf("Update request %v", payload)
	ctxHeader, err := contextutil.CentContext(ctx, h.config)
	if err != nil {
		apiLog.Error(err)
		return nil, err
	}

	doc, err := h.service.DeriveFromUpdatePayload(ctxHeader, payload)
	if err != nil {
		apiLog.Error(err)
		return nil, centerrors.Wrap(err, "could not derive update payload")
	}

	doc, txID, err := h.service.Update(ctxHeader, doc)
	if err != nil {
		apiLog.Error(err)
		return nil, centerrors.Wrap(err, "could not update document")
	}

	resp, err := h.service.DerivePurchaseOrderResponse(doc)
	if err != nil {
		apiLog.Error(err)
		return nil, centerrors.Wrap(err, "could not derive response")
	}

	resp.Header.TransactionId = txID.String()
	return resp, nil
}

// GetVersion returns the requested version of a purchase order
func (h grpcHandler) GetVersion(ctx context.Context, req *clientpurchaseorderpb.GetVersionRequest) (*clientpurchaseorderpb.PurchaseOrderResponse, error) {
	apiLog.Debugf("GetVersion request %v", req)
	ctxHeader, err := contextutil.CentContext(ctx, h.config)
	if err != nil {
		apiLog.Error(err)
		return nil, err
	}

	identifier, err := hexutil.Decode(req.Identifier)
	if err != nil {
		apiLog.Error(err)
		return nil, centerrors.Wrap(err, "identifier is invalid")
	}

	version, err := hexutil.Decode(req.Version)
	if err != nil {
		apiLog.Error(err)
		return nil, centerrors.Wrap(err, "version is invalid")
	}

	model, err := h.service.GetVersion(ctxHeader, identifier, version)
	if err != nil {
		apiLog.Error(err)
		return nil, centerrors.Wrap(err, "document not found")
	}

	resp, err := h.service.DerivePurchaseOrderResponse(model)
	if err != nil {
		apiLog.Error(err)
		return nil, centerrors.Wrap(err, "could not derive response")
	}

	return resp, nil
}

// Get returns the purchase order the latest version of the document with given identifier
func (h grpcHandler) Get(ctx context.Context, getRequest *clientpurchaseorderpb.GetRequest) (*clientpurchaseorderpb.PurchaseOrderResponse, error) {
	apiLog.Debugf("Get request %v", getRequest)
	ctxHeader, err := contextutil.CentContext(ctx, h.config)
	if err != nil {
		apiLog.Error(err)
		return nil, err
	}

	identifier, err := hexutil.Decode(getRequest.Identifier)
	if err != nil {
		apiLog.Error(err)
		return nil, centerrors.Wrap(err, "identifier is an invalid hex string")
	}

	model, err := h.service.GetCurrentVersion(ctxHeader, identifier)
	if err != nil {
		apiLog.Error(err)
		return nil, centerrors.Wrap(err, "document not found")
	}

	resp, err := h.service.DerivePurchaseOrderResponse(model)
	if err != nil {
		apiLog.Error(err)
		return nil, centerrors.Wrap(err, "could not derive response")
	}

	return resp, nil
}
