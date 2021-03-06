package ethid

import (
	"github.com/centrifuge/go-centrifuge/config"
	"github.com/centrifuge/go-centrifuge/config/configstore"
	"github.com/centrifuge/go-centrifuge/errors"
	"github.com/centrifuge/go-centrifuge/identity"

	"github.com/centrifuge/go-centrifuge/bootstrap"
	"github.com/centrifuge/go-centrifuge/ethereum"
	"github.com/centrifuge/go-centrifuge/queue"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// Bootstrapper implements bootstrap.Bootstrapper.
type Bootstrapper struct{}

// Bootstrap initializes the IdentityFactoryContract as well as the idRegistrationConfirmationTask that depends on it.
// the idRegistrationConfirmationTask is added to be registered on the queue at queue.Bootstrapper
func (*Bootstrapper) Bootstrap(context map[string]interface{}) error {
	// we have to allow loading from file in case this is coming from create config cmd where we don't add configs to db
	cfg, err := configstore.RetrieveConfig(false, context)
	if err != nil {
		return err
	}

	if _, ok := context[ethereum.BootstrappedEthereumClient]; !ok {
		return errors.New("ethereum client hasn't been initialized")
	}
	gethClient := context[ethereum.BootstrappedEthereumClient].(ethereum.Client)

	idFactory, err := getIdentityFactoryContract(cfg.GetContractAddress(config.IdentityFactory), gethClient)
	if err != nil {
		return err
	}

	registryContract, err := getIdentityRegistryContract(cfg.GetContractAddress(config.IdentityRegistry), gethClient)
	if err != nil {
		return err
	}

	if _, ok := context[bootstrap.BootstrappedQueueServer]; !ok {
		return errors.New("queue hasn't been initialized")
	}
	queueSrv := context[bootstrap.BootstrappedQueueServer].(*queue.Server)

	context[identity.BootstrappedIDService] = NewEthereumIdentityService(cfg, idFactory, registryContract, queueSrv, ethereum.GetClient,
		func(address common.Address, backend bind.ContractBackend) (contract, error) {
			return NewEthereumIdentityContract(address, backend)
		})

	idRegTask := newIDRegistrationConfirmationTask(cfg.GetEthereumContextWaitTimeout(), &idFactory.EthereumIdentityFactoryContractFilterer, ethereum.DefaultWaitForTransactionMiningContext)
	keyRegTask := newKeyRegistrationConfirmationTask(ethereum.DefaultWaitForTransactionMiningContext, registryContract, cfg, queueSrv, ethereum.GetClient,
		func(address common.Address, backend bind.ContractBackend) (contract, error) {
			return NewEthereumIdentityContract(address, backend)
		})
	queueSrv.RegisterTaskType(idRegTask.TaskTypeName(), idRegTask)
	queueSrv.RegisterTaskType(keyRegTask.TaskTypeName(), keyRegTask)
	return nil
}

func getIdentityFactoryContract(factoryAddress common.Address, ethClient ethereum.Client) (identityFactoryContract *EthereumIdentityFactoryContract, err error) {
	return NewEthereumIdentityFactoryContract(factoryAddress, ethClient.GetEthClient())
}

func getIdentityRegistryContract(registryAddress common.Address, ethClient ethereum.Client) (identityRegistryContract *EthereumIdentityRegistryContract, err error) {
	return NewEthereumIdentityRegistryContract(registryAddress, ethClient.GetEthClient())
}
