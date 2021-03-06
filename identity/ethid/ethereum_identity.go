package ethid

import (
	"context"
	"fmt"
	"math/big"

	"github.com/centrifuge/go-centrifuge/contextutil"

	"github.com/centrifuge/go-centrifuge/crypto"
	"github.com/centrifuge/go-centrifuge/identity"

	"bytes"

	"time"

	"github.com/centrifuge/centrifuge-protobufs/gen/go/coredocument"
	"github.com/centrifuge/go-centrifuge/crypto/ed25519"
	"github.com/centrifuge/go-centrifuge/errors"
	"github.com/centrifuge/go-centrifuge/ethereum"
	"github.com/centrifuge/go-centrifuge/queue"
	"github.com/centrifuge/go-centrifuge/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("identity")

type factory interface {
	CreateIdentity(opts *bind.TransactOpts, centID *big.Int) (*types.Transaction, error)
}

type registry interface {
	GetIdentityByCentrifugeId(opts *bind.CallOpts, bigInt *big.Int) (common.Address, error)
}

type contract interface {
	AddKey(opts *bind.TransactOpts, _key [32]byte, _kPurpose *big.Int) (*types.Transaction, error)

	GetKeysByPurpose(opts *bind.CallOpts, _purpose *big.Int) ([][32]byte, error)

	GetKey(opts *bind.CallOpts, _key [32]byte) (struct {
		Key       [32]byte
		Purposes  []*big.Int
		RevokedAt *big.Int
	}, error)

	FilterKeyAdded(opts *bind.FilterOpts, key [][32]byte, purpose []*big.Int) (*EthereumIdentityContractKeyAddedIterator, error)
}

// EthereumIdentityKey holds the identity related details
type EthereumIdentityKey struct {
	Key       [32]byte
	Purposes  []*big.Int
	RevokedAt *big.Int
}

// GetKey returns the public key
func (idk *EthereumIdentityKey) GetKey() [32]byte {
	return idk.Key
}

// GetPurposes returns the purposes intended for the key
func (idk *EthereumIdentityKey) GetPurposes() []*big.Int {
	return idk.Purposes
}

// GetRevokedAt returns the block at which the identity is revoked
func (idk *EthereumIdentityKey) GetRevokedAt() *big.Int {
	return idk.RevokedAt
}

// String prints the peerID extracted from the key
func (idk *EthereumIdentityKey) String() string {
	peerID, _ := ed25519.PublicKeyToP2PKey(idk.Key)
	return fmt.Sprintf("%s", peerID.Pretty())
}

type ethereumIdentity struct {
	centID           identity.CentID
	contract         contract
	contractProvider func(address common.Address, backend bind.ContractBackend) (contract, error)
	registryContract registry
	// TODO [multi-tenancy] replace this with config service
	config           identity.Config
	gethClientFinder func() ethereum.Client
	queue            *queue.Server
}

func newEthereumIdentity(id identity.CentID, registryContract registry, config identity.Config,
	queue *queue.Server,
	gethClientFinder func() ethereum.Client,
	contractProvider func(address common.Address, backend bind.ContractBackend) (contract, error)) *ethereumIdentity {
	return &ethereumIdentity{centID: id, registryContract: registryContract, config: config, gethClientFinder: gethClientFinder, contractProvider: contractProvider, queue: queue}
}

// CentrifugeID sets the CentID to the Identity
func (id *ethereumIdentity) SetCentrifugeID(centID identity.CentID) {
	id.centID = centID
}

// String returns CentrifugeID
func (id *ethereumIdentity) String() string {
	return fmt.Sprintf("CentrifugeID [%s]", id.centID)
}

// CentrifugeID returns the CentrifugeID
func (id *ethereumIdentity) CentID() identity.CentID {
	return id.centID
}

// LastKeyForPurpose returns the latest key for given purpose
func (id *ethereumIdentity) LastKeyForPurpose(keyPurpose int) (key []byte, err error) {
	idKeys, err := id.fetchKeysByPurpose(keyPurpose)
	if err != nil {
		return []byte{}, err
	}

	if len(idKeys) == 0 {
		return []byte{}, errors.New("no key found for type [%d] in ID [%s]", keyPurpose, id.centID)
	}

	return idKeys[len(idKeys)-1].Key[:32], nil
}

// FetchKey fetches the Key from the chain
func (id *ethereumIdentity) FetchKey(key []byte) (identity.Key, error) {
	contract, err := id.getContract()
	if err != nil {
		return nil, err
	}
	// Ignoring cancelFunc as code will block until response or timeout is triggered
	opts, _ := id.gethClientFinder().GetGethCallOpts()
	key32, _ := utils.SliceToByte32(key)
	keyInstance, err := contract.GetKey(opts, key32)
	if err != nil {
		return nil, err
	}

	return &EthereumIdentityKey{
		Key:       keyInstance.Key,
		Purposes:  keyInstance.Purposes,
		RevokedAt: keyInstance.RevokedAt,
	}, nil

}

// CurrentP2PKey returns the latest P2P key
func (id *ethereumIdentity) CurrentP2PKey() (ret string, err error) {
	key, err := id.LastKeyForPurpose(identity.KeyPurposeP2P)
	if err != nil {
		return ret, err
	}
	key32, _ := utils.SliceToByte32(key)
	p2pID, err := ed25519.PublicKeyToP2PKey(key32)
	if err != nil {
		return ret, err
	}

	return p2pID.Pretty(), nil
}

func (id *ethereumIdentity) findContract() (exists bool, err error) {
	if id.contract != nil {
		return true, nil
	}

	client := id.gethClientFinder()
	// Ignoring cancelFunc as code will block until response or timeout is triggered
	opts, _ := client.GetGethCallOpts()
	idAddress, err := id.registryContract.GetIdentityByCentrifugeId(opts, id.centID.BigInt())
	if err != nil {
		return false, err
	}
	if utils.IsEmptyByteSlice(idAddress.Bytes()) {
		return false, errors.New("Identity not found by address provided")
	}

	idContract, err := id.contractProvider(idAddress, client.GetEthClient())
	if err == bind.ErrNoCode {
		return false, err
	}
	if err != nil {
		log.Errorf("Failed to instantiate the identity contract: %v", err)
		return false, err
	}
	id.contract = idContract
	return true, nil
}

func (id *ethereumIdentity) getContract() (contract contract, err error) {
	if id.contract == nil {
		_, err := id.findContract()
		if err != nil {
			return nil, err
		}
		return id.contract, nil
	}
	return id.contract, nil
}

// AddKeyToIdentity adds key to the purpose on chain
func (id *ethereumIdentity) AddKeyToIdentity(ctx context.Context, keyPurpose int, key []byte) (confirmations chan *identity.WatchIdentity, err error) {
	if utils.IsEmptyByteSlice(key) || len(key) > 32 {
		log.Errorf("Can't add key to identity: empty or invalid length(>32) key for [id: %s]: %x", id, key)
		return confirmations, errors.New("Can't add key to identity: Invalid key")
	}

	ethIdentityContract, err := id.getContract()
	if err != nil {
		log.Errorf("Failed to set up event listener for identity [id: %s]: %v", id, err)
		return confirmations, err
	}

	conn := id.gethClientFinder()
	opts, err := ethereum.GetClient().GetTxOpts(id.config.GetEthereumDefaultAccountName())
	if err != nil {
		return confirmations, err
	}

	h, err := conn.GetEthClient().HeaderByNumber(ctx, nil)
	if err != nil {
		return confirmations, err
	}

	var keyFixed [32]byte
	copy(keyFixed[:], key)
	confirmations, err = id.setUpKeyRegisteredEventListener(id.config, id, keyPurpose, keyFixed, h.Number.Uint64())
	if err != nil {
		wError := errors.New("%v", err)
		log.Errorf("Failed to set up event listener for identity [id: %s]: %v", id, wError)
		return confirmations, wError
	}

	err = sendKeyRegistrationTransaction(ethIdentityContract, opts, id, keyPurpose, key)
	if err != nil {
		wError := errors.New("%v", err)
		log.Errorf("Failed to create transaction for identity [id: %s]: %v", id, wError)
		return confirmations, wError
	}
	return confirmations, nil
}

// fetchKeysByPurpose fetches keys from chain matching purpose
func (id *ethereumIdentity) fetchKeysByPurpose(keyPurpose int) ([]EthereumIdentityKey, error) {
	contract, err := id.getContract()
	if err != nil {
		return nil, err
	}
	// Ignoring cancelFunc as code will block until response or timeout is triggered
	opts, _ := id.gethClientFinder().GetGethCallOpts()
	bigInt := big.NewInt(int64(keyPurpose))
	keys, err := contract.GetKeysByPurpose(opts, bigInt)
	if err != nil {
		return nil, err
	}
	log.Infof("fetched keys: %d %x\n", keyPurpose, keys)

	var ids []EthereumIdentityKey
	for _, key := range keys {
		ids = append(ids, EthereumIdentityKey{Key: key})
	}

	return ids, nil
}

// sendRegistrationTransaction sends the actual transaction to add a Key on Ethereum registry contract
func sendKeyRegistrationTransaction(identityContract contract, opts *bind.TransactOpts, identity *ethereumIdentity, keyPurpose int, key []byte) error {

	//preparation of data in specific types for the call to Ethereum
	bigInt := big.NewInt(int64(keyPurpose))
	bKey, err := utils.SliceToByte32(key)
	if err != nil {
		return err
	}

	tx, err := ethereum.GetClient().SubmitTransactionWithRetries(identityContract.AddKey, opts, bKey, bigInt)
	if err != nil {
		log.Infof("Failed to send key [%v:%x] to add to CentrifugeID [%x]: %v", keyPurpose, bKey, identity.CentID, err)
		return err
	}

	log.Infof("Sent off key [%v:%x] to add to CentrifugeID [%s]. Ethereum transaction hash [%x]", keyPurpose, bKey, identity, tx.Hash())
	return nil
}

// sendIdentityCreationTransaction sends the actual transaction to create identity on Ethereum registry contract
func sendIdentityCreationTransaction(identityFactory factory, opts *bind.TransactOpts, identityToBeCreated identity.Identity) error {
	//preparation of data in specific types for the call to Ethereum
	tx, err := ethereum.GetClient().SubmitTransactionWithRetries(identityFactory.CreateIdentity, opts, identityToBeCreated.CentID().BigInt())
	if err != nil {
		log.Infof("Failed to send identity for creation [CentrifugeID: %s] : %v", identityToBeCreated, err)
		return err
	}

	log.Infof("Sent off identity creation [%s]. Ethereum transaction hash [%x] and Nonce [%v] and Check [%v]", identityToBeCreated, tx.Hash(), tx.Nonce(), tx.CheckNonce())
	log.Infof("Transfer pending: 0x%x\n", tx.Hash())
	return err
}

// setUpKeyRegisteredEventListener listens for Identity creation
func (id *ethereumIdentity) setUpKeyRegisteredEventListener(config identity.Config, i identity.Identity, keyPurpose int, key [32]byte, bh uint64) (confirmations chan *identity.WatchIdentity, err error) {
	confirmations = make(chan *identity.WatchIdentity)
	centID := i.CentID()
	if err != nil {
		return nil, err
	}
	asyncRes, err := id.queue.EnqueueJob(keyRegistrationConfirmationTaskName,
		map[string]interface{}{
			centIDParam:            centID,
			keyParam:               key,
			keyPurposeParam:        keyPurpose,
			queue.BlockHeightParam: bh,
		})
	if err != nil {
		return nil, err
	}
	go waitAndRouteKeyRegistrationEvent(config.GetEthereumContextWaitTimeout(), asyncRes, confirmations, i)
	return confirmations, nil
}

// setUpRegistrationEventListener sets up the listened for the "IdentityCreated" event to notify the upstream code about successful mining/creation
// of the identity.
func (ids *ethereumIdentityService) setUpRegistrationEventListener(config identity.Config, identityToBeCreated identity.Identity, blockHeight uint64) (confirmations chan *identity.WatchIdentity, err error) {
	confirmations = make(chan *identity.WatchIdentity)
	centID := identityToBeCreated.CentID()
	if err != nil {
		return nil, err
	}

	asyncRes, err := ids.queue.EnqueueJob(idRegistrationConfirmationTaskName, map[string]interface{}{centIDParam: centID, queue.BlockHeightParam: blockHeight})
	if err != nil {
		return nil, err
	}
	go waitAndRouteIdentityRegistrationEvent(config.GetEthereumContextWaitTimeout(), asyncRes, confirmations, identityToBeCreated)
	return confirmations, nil
}

// waitAndRouteKeyRegistrationEvent notifies the confirmations channel whenever the key has been added to the identity and has been noted as Ethereum event
func waitAndRouteKeyRegistrationEvent(timeout time.Duration, asyncRes queue.TaskResult, confirmations chan<- *identity.WatchIdentity, pushThisIdentity identity.Identity) {
	_, err := asyncRes.Get(timeout)
	confirmations <- &identity.WatchIdentity{Identity: pushThisIdentity, Error: err}
}

// waitAndRouteIdentityRegistrationEvent notifies the confirmations channel whenever the identity creation is being noted as Ethereum event
func waitAndRouteIdentityRegistrationEvent(timeout time.Duration, asyncRes queue.TaskResult, confirmations chan<- *identity.WatchIdentity, pushThisIdentity identity.Identity) {
	_, err := asyncRes.Get(timeout)
	confirmations <- &identity.WatchIdentity{pushThisIdentity, err}
}

// ethereumIdentityService implements `Service`
type ethereumIdentityService struct {
	// TODO [multi-tenancy] replace this with config service
	config           identity.Config
	factoryContract  factory
	registryContract registry
	gethClientFinder func() ethereum.Client
	contractProvider func(address common.Address, backend bind.ContractBackend) (contract, error)
	queue            *queue.Server
}

// NewEthereumIdentityService creates a new NewEthereumIdentityService given the config and the smart contracts
func NewEthereumIdentityService(config identity.Config, factoryContract factory, registryContract registry,
	queue *queue.Server,
	gethClientFinder func() ethereum.Client,
	contractProvider func(address common.Address, backend bind.ContractBackend) (contract, error)) identity.Service {
	return &ethereumIdentityService{config: config, factoryContract: factoryContract, registryContract: registryContract, gethClientFinder: gethClientFinder, contractProvider: contractProvider, queue: queue}
}

// CheckIdentityExists checks if the identity represented by id actually exists on ethereum
func (ids *ethereumIdentityService) CheckIdentityExists(centrifugeID identity.CentID) (exists bool, err error) {
	client := ids.gethClientFinder()
	// Ignoring cancelFunc as code will block until response or timeout is triggered
	opts, _ := client.GetGethCallOpts()
	idAddress, err := ids.registryContract.GetIdentityByCentrifugeId(opts, centrifugeID.BigInt())
	if err != nil {
		return false, err
	}
	if utils.IsEmptyByteSlice(idAddress.Bytes()) {
		return false, errors.New("Identity not found by address provided")
	}

	_, err = NewEthereumIdentityContract(idAddress, client.GetEthClient())
	if err == bind.ErrNoCode {
		return false, err
	}
	if err != nil {
		log.Errorf("Failed to instantiate the identity contract: %v", err)
		return false, err
	}
	return true, nil
}

// CreateIdentity creates an identity representing the id on ethereum
func (ids *ethereumIdentityService) CreateIdentity(ctx context.Context, centrifugeID identity.CentID) (id identity.Identity, confirmations chan *identity.WatchIdentity, err error) {
	log.Infof("Creating Identity [%x]", centrifugeID)
	tc, err := contextutil.Tenant(ctx)
	if err != nil {
		return nil, confirmations, err
	}

	id = newEthereumIdentity(centrifugeID, ids.registryContract, ids.config, ids.queue, ids.gethClientFinder, ids.contractProvider)
	conn := ids.gethClientFinder()
	opts, err := conn.GetTxOpts(tc.GetEthereumDefaultAccountName())
	if err != nil {
		return nil, confirmations, err
	}

	h, err := conn.GetEthClient().HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, confirmations, err
	}

	confirmations, err = ids.setUpRegistrationEventListener(ids.config, id, h.Number.Uint64())
	if err != nil {
		wError := errors.New("%v", err)
		log.Infof("Failed to set up event listener for identity [mockID: %s]: %v", id, wError)
		return nil, confirmations, wError
	}

	err = sendIdentityCreationTransaction(ids.factoryContract, opts, id)
	if err != nil {
		wError := errors.New("%v", err)
		log.Infof("Failed to create transaction for identity [mockID: %s]: %v", id, wError)
		return nil, confirmations, wError
	}
	return id, confirmations, nil
}

// GetIdentityAddress gets the address of the ethereum identity contract for the given CentID
func (ids *ethereumIdentityService) GetIdentityAddress(centID identity.CentID) (common.Address, error) {
	// Ignoring cancelFunc as code will block until response or timeout is triggered
	opts, _ := ethereum.GetClient().GetGethCallOpts()
	address, err := ids.registryContract.GetIdentityByCentrifugeId(opts, centID.BigInt())
	if err != nil {
		return common.Address{}, err
	}

	if utils.IsEmptyAddress(address) {
		return common.Address{}, errors.New("No address found for centID")
	}
	return address, nil
}

// LookupIdentityForID looks up if the identity for given CentID exists on ethereum
func (ids *ethereumIdentityService) LookupIdentityForID(centrifugeID identity.CentID) (identity.Identity, error) {
	exists, err := ids.CheckIdentityExists(centrifugeID)
	if !exists {
		return nil, errors.New("identity [%s] does not exist with err [%v]", centrifugeID, err)
	}

	if err != nil {
		return nil, err
	}
	return newEthereumIdentity(centrifugeID, ids.registryContract, ids.config, ids.queue, ids.gethClientFinder, ids.contractProvider), nil
}

// GetClientP2PURL returns the p2p url associated with the centID
func (ids *ethereumIdentityService) GetClientP2PURL(centID identity.CentID) (url string, err error) {
	target, err := ids.LookupIdentityForID(centID)
	if err != nil {
		return url, errors.New("error fetching receiver identity: %v", err)
	}

	p2pKey, err := target.CurrentP2PKey()
	if err != nil {
		return url, errors.New("error fetching p2p key: %v", err)
	}

	return fmt.Sprintf("/ipfs/%s", p2pKey), nil
}

// GetClientsP2PURLs returns p2p urls associated with each centIDs
// will error out at first failure
func (ids *ethereumIdentityService) GetClientsP2PURLs(centIDs []identity.CentID) ([]string, error) {
	var p2pURLs []string
	for _, id := range centIDs {
		url, err := ids.GetClientP2PURL(id)
		if err != nil {
			return nil, err
		}

		p2pURLs = append(p2pURLs, url)
	}

	return p2pURLs, nil
}

// GetIdentityKey returns the key for provided identity
func (ids *ethereumIdentityService) GetIdentityKey(identity identity.CentID, pubKey []byte) (keyInfo identity.Key, err error) {
	id, err := ids.LookupIdentityForID(identity)
	if err != nil {
		return keyInfo, err
	}

	key, err := id.FetchKey(pubKey)
	if err != nil {
		return keyInfo, err
	}

	if utils.IsEmptyByte32(key.GetKey()) {
		return keyInfo, errors.New("key not found for identity: %x", identity)
	}

	return key, nil
}

// ValidateKey checks if a given key is valid for the given centrifugeID.
func (ids *ethereumIdentityService) ValidateKey(centID identity.CentID, key []byte, purpose int) error {
	idKey, err := ids.GetIdentityKey(centID, key)
	if err != nil {
		return err
	}

	if !bytes.Equal(key, utils.Byte32ToSlice(idKey.GetKey())) {
		return errors.New("[Key: %x] Key doesn't match", idKey.GetKey())
	}

	if !utils.ContainsBigIntInSlice(big.NewInt(int64(purpose)), idKey.GetPurposes()) {
		return errors.New("[Key: %x] Key doesn't have purpose [%d]", idKey.GetKey(), purpose)
	}

	if idKey.GetRevokedAt().Cmp(big.NewInt(0)) != 0 {
		return errors.New("[Key: %x] Key is currently revoked since block [%d]", idKey.GetKey(), idKey.GetRevokedAt())
	}

	return nil
}

// AddKeyFromConfig adds a key previously generated and indexed in the configuration file to the identity specified in such config file
func (ids *ethereumIdentityService) AddKeyFromConfig(config identity.Config, purpose int) error {
	identityConfig, err := identity.GetIdentityConfig(config)
	if err != nil {
		return err
	}

	id, err := ids.LookupIdentityForID(identityConfig.ID)
	if err != nil {
		return err
	}

	ctx, cancel := ethereum.DefaultWaitForTransactionMiningContext(config.GetEthereumContextWaitTimeout())
	defer cancel()
	confirmations, err := id.AddKeyToIdentity(ctx, purpose, identityConfig.Keys[purpose].PublicKey)
	if err != nil {
		return err
	}
	watchAddedToIdentity := <-confirmations

	lastKey, errLocal := watchAddedToIdentity.Identity.LastKeyForPurpose(purpose)
	if errLocal != nil {
		return err
	}

	log.Infof("Key [%v] with type [%d] Added to Identity [%s]", lastKey, purpose, watchAddedToIdentity.Identity)

	return nil
}

// ValidateSignature validates a signature on a message based on identity data
func (ids *ethereumIdentityService) ValidateSignature(signature *coredocumentpb.Signature, message []byte) error {
	centID, err := identity.ToCentID(signature.EntityId)
	if err != nil {
		return err
	}

	err = ids.ValidateKey(centID, signature.PublicKey, identity.KeyPurposeSigning)
	if err != nil {
		return err
	}

	return crypto.VerifySignature(signature.PublicKey, message, signature.Signature)
}
