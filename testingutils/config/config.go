package testingconfig

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/centrifuge/go-centrifuge/config/configstore"
	"github.com/centrifuge/go-centrifuge/contextutil"
	"github.com/stretchr/testify/assert"

	"github.com/centrifuge/go-centrifuge/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
)

type MockConfig struct {
	config.Configuration
	mock.Mock
}

func (m *MockConfig) GetStoragePath() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockConfig) GetP2PPort() int {
	args := m.Called()
	return args.Get(0).(int)
}

func (m *MockConfig) GetP2PExternalIP() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockConfig) GetP2PConnectionTimeout() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func (m *MockConfig) GetReceiveEventNotificationEndpoint() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockConfig) GetServerPort() int {
	args := m.Called()
	return args.Get(0).(int)
}

func (m *MockConfig) GetServerAddress() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockConfig) GetNumWorkers() int {
	args := m.Called()
	return args.Get(0).(int)
}

func (m *MockConfig) GetWorkerWaitTimeMS() int {
	args := m.Called()
	return args.Get(0).(int)
}

func (m *MockConfig) GetEthereumNodeURL() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockConfig) GetEthereumContextReadWaitTimeout() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func (m *MockConfig) GetEthereumContextWaitTimeout() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func (m *MockConfig) GetEthereumIntervalRetry() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func (m *MockConfig) GetEthereumMaxRetries() int {
	args := m.Called()
	return args.Get(0).(int)
}

func (m *MockConfig) GetEthereumGasPrice() *big.Int {
	args := m.Called()
	return args.Get(0).(*big.Int)
}

func (m *MockConfig) GetEthereumGasLimit() uint64 {
	args := m.Called()
	return args.Get(0).(uint64)
}

func (m *MockConfig) GetEthereumDefaultAccountName() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockConfig) GetEthereumAccount(accountName string) (account *config.AccountConfig, err error) {
	args := m.Called()
	return args.Get(0).(*config.AccountConfig), args.Error(1)
}

func (m *MockConfig) GetTxPoolAccessEnabled() bool {
	args := m.Called()
	return args.Get(0).(bool)
}

func (m *MockConfig) GetNetworkString() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockConfig) GetNetworkKey(k string) string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockConfig) GetContractAddressString(address string) string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockConfig) GetContractAddress(contractName config.ContractName) common.Address {
	args := m.Called()
	return args.Get(0).(common.Address)
}

func (m *MockConfig) GetBootstrapPeers() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockConfig) GetNetworkID() uint32 {
	args := m.Called()
	return args.Get(0).(uint32)
}

func (m *MockConfig) GetIdentityID() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockConfig) GetSigningKeyPair() (pub, priv string) {
	args := m.Called()
	return args.Get(0).(string), args.Get(1).(string)
}

func (m *MockConfig) GetEthAuthKeyPair() (pub, priv string) {
	args := m.Called()
	return args.Get(0).(string), args.Get(1).(string)
}

func CreateTenantContext(t *testing.T, cfg config.Configuration) context.Context {
	return CreateTenantContextWithContext(t, context.Background(), cfg)
}

func CreateTenantContextWithContext(t *testing.T, ctx context.Context, cfg config.Configuration) context.Context {
	tc, err := configstore.NewTenantConfig("", cfg)
	assert.Nil(t, err)

	contextHeader, err := contextutil.NewCentrifugeContext(ctx, tc)
	assert.Nil(t, err)
	return contextHeader
}
