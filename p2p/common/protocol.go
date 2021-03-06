package p2pcommon

import (
	"context"
	"fmt"
	"strings"

	"github.com/centrifuge/go-centrifuge/crypto"

	"github.com/centrifuge/centrifuge-protobufs/gen/go/p2p"
	"github.com/centrifuge/go-centrifuge/contextutil"
	"github.com/centrifuge/go-centrifuge/errors"
	"github.com/centrifuge/go-centrifuge/identity"
	"github.com/centrifuge/go-centrifuge/protobufs/gen/go/protocol"
	"github.com/centrifuge/go-centrifuge/version"
	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p-protocol"
)

// MessageType holds the protocol message type
type MessageType string

const (
	// CentrifugeProtocol is the centrifuge wire protocol
	CentrifugeProtocol protocol.ID = "/centrifuge/0.0.1"

	// MessageTypeError defines any protocol error
	MessageTypeError MessageType = "MessageTypeError"
	// MessageTypeInvalid defines invalid protocol type
	MessageTypeInvalid MessageType = "MessageTypeInvalid"
	// MessageTypeRequestSignature defines RequestSignature type
	MessageTypeRequestSignature MessageType = "MessageTypeRequestSignature"
	// MessageTypeRequestSignatureRep defines RequestSignature response type
	MessageTypeRequestSignatureRep MessageType = "MessageTypeRequestSignatureRep"
	// MessageTypeSendAnchoredDoc defines SendAnchored type
	MessageTypeSendAnchoredDoc MessageType = "MessageTypeSendAnchoredDoc"
	// MessageTypeSendAnchoredDocRep defines SendAnchored response type
	MessageTypeSendAnchoredDocRep MessageType = "MessageTypeSendAnchoredDocRep"
)

// Equals compares if string is of a particular MessageType
func (mt MessageType) Equals(mt2 string) bool {
	return mt.String() == mt2
}

// String representation
func (mt MessageType) String() string {
	return string(mt)
}

// MessageTypeFromString Resolves MessageType out of string
func MessageTypeFromString(mt string) MessageType {
	var found MessageType
	if MessageTypeError.Equals(mt) {
		found = MessageTypeError
	} else if MessageTypeInvalid.Equals(mt) {
		found = MessageTypeInvalid
	} else if MessageTypeRequestSignature.Equals(mt) {
		found = MessageTypeRequestSignature
	} else if MessageTypeRequestSignatureRep.Equals(mt) {
		found = MessageTypeRequestSignatureRep
	} else if MessageTypeSendAnchoredDoc.Equals(mt) {
		found = MessageTypeSendAnchoredDoc
	} else if MessageTypeSendAnchoredDocRep.Equals(mt) {
		found = MessageTypeSendAnchoredDocRep
	}
	return found
}

// ProtocolForCID creates the protocol string for the given CID
func ProtocolForCID(CID identity.CentID) protocol.ID {
	return protocol.ID(fmt.Sprintf("%s/%s", CentrifugeProtocol, CID.String()))
}

// ExtractCID extracts CID from a protocol string
func ExtractCID(id protocol.ID) (identity.CentID, error) {
	parts := strings.Split(string(id), "/")
	cidHexStr := parts[len(parts)-1]
	return identity.CentIDFromString(cidHexStr)
}

// ResolveDataEnvelope unwraps Content Envelope out of p2pEnvelope
func ResolveDataEnvelope(mes proto.Message) (*p2ppb.Envelope, error) {
	recv, ok := mes.(*protocolpb.P2PEnvelope)
	if !ok {
		return nil, errors.New("cannot cast proto.Message to protocolpb.P2PEnvelope: %v", recv)
	}
	recvEnvelope := new(p2ppb.Envelope)
	err := proto.Unmarshal(recv.Body, recvEnvelope)
	if err != nil {
		return nil, err
	}

	// Validate at least not-nil fields
	if recvEnvelope.Header == nil {
		return nil, errors.New("Header field is empty")
	}

	return recvEnvelope, nil
}

// PrepareP2PEnvelope wraps content message into p2p envelope
func PrepareP2PEnvelope(ctx context.Context, networkID uint32, messageType MessageType, mes proto.Message) (*protocolpb.P2PEnvelope, error) {
	self, err := contextutil.Self(ctx)
	if err != nil {
		return nil, err
	}

	centIDBytes := self.ID[:]
	p2pheader := &p2ppb.Header{
		SenderId:          centIDBytes,
		NodeVersion:       version.GetVersion().String(),
		NetworkIdentifier: networkID,
		Type:              messageType.String(),
	}

	body, err := proto.Marshal(mes)
	if err != nil {
		return nil, err
	}

	envelope := &p2ppb.Envelope{
		Header: p2pheader,
		Body:   body,
	}

	signRequest, err := proto.Marshal(envelope)
	if err != nil {
		return nil, err
	}

	signKeys := self.Keys[identity.KeyPurposeSigning]
	envelope.Header.Signature = crypto.Sign(self.ID[:], signKeys.PrivateKey, signKeys.PublicKey, signRequest)

	marshalledRequest, err := proto.Marshal(envelope)
	if err != nil {
		return nil, err
	}

	return &protocolpb.P2PEnvelope{Body: marshalledRequest}, nil
}
