package nft

import (
	"context"
)

// PaymentObligation handles transactions related to minting of NFTs
type PaymentObligation interface {

	// MintNFT mints an NFT
	MintNFT(ctx context.Context, documentID []byte, registryAddress, depositAddress string, proofFields []string) (*MintNFTResponse, error)
}

// MintNFTResponse holds tokenID and transaction ID.
type MintNFTResponse struct {
	TokenID       string
	TransactionID string
}
