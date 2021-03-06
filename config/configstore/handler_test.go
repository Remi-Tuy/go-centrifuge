// +build unit

package configstore

import (
	"context"
	"testing"

	"github.com/centrifuge/go-centrifuge/testingutils/commons"

	"github.com/centrifuge/go-centrifuge/protobufs/gen/go/config"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestGrpcHandler_GetConfigNoConfig(t *testing.T) {
	idService := &testingcommons.MockIDService{}
	repo, _, err := getRandomStorage()
	assert.Nil(t, err)
	repo.RegisterConfig(&NodeConfig{})
	svc := DefaultService(repo, idService)
	h := GRPCHandler(svc)
	readCfg, err := h.GetConfig(context.Background(), nil)
	assert.NotNil(t, err)
	assert.Nil(t, readCfg)
}

func TestGrpcHandler_GetConfig(t *testing.T) {
	idService := &testingcommons.MockIDService{}
	repo, _, err := getRandomStorage()
	assert.Nil(t, err)
	repo.RegisterConfig(&NodeConfig{})
	svc := DefaultService(repo, idService)
	h := GRPCHandler(svc)
	nodeCfg := NewNodeConfig(cfg)
	_, err = h.CreateConfig(context.Background(), nodeCfg.CreateProtobuf())
	assert.Nil(t, err)
	readCfg, err := h.GetConfig(context.Background(), nil)
	assert.Nil(t, err)
	assert.NotNil(t, readCfg)
}

func TestGrpcHandler_GetTenantNoConfig(t *testing.T) {
	idService := &testingcommons.MockIDService{}
	repo, _, err := getRandomStorage()
	assert.Nil(t, err)
	repo.RegisterTenant(&TenantConfig{})
	svc := DefaultService(repo, idService)
	h := GRPCHandler(svc)
	readCfg, err := h.GetTenant(context.Background(), &configpb.GetTenantRequest{Identifier: "0x123456789"})
	assert.NotNil(t, err)
	assert.Nil(t, readCfg)
}

func TestGrpcHandler_GetTenant(t *testing.T) {
	idService := &testingcommons.MockIDService{}
	repo, _, err := getRandomStorage()
	assert.Nil(t, err)
	repo.RegisterTenant(&TenantConfig{})
	svc := DefaultService(repo, idService)
	h := GRPCHandler(svc)
	tenantCfg, err := NewTenantConfig("main", cfg)
	assert.Nil(t, err)
	_, err = h.CreateTenant(context.Background(), tenantCfg.CreateProtobuf())
	assert.Nil(t, err)
	tid, err := tenantCfg.GetIdentityID()
	assert.Nil(t, err)
	readCfg, err := h.GetTenant(context.Background(), &configpb.GetTenantRequest{Identifier: hexutil.Encode(tid)})
	assert.Nil(t, err)
	assert.NotNil(t, readCfg)
}

func TestGrpcHandler_GetAllTenants(t *testing.T) {
	idService := &testingcommons.MockIDService{}
	repo, _, err := getRandomStorage()
	assert.Nil(t, err)
	repo.RegisterTenant(&TenantConfig{})
	svc := DefaultService(repo, idService)
	h := GRPCHandler(svc)
	tenantCfg1, err := NewTenantConfig("main", cfg)
	tenantCfg2, err := NewTenantConfig("main", cfg)
	tc := tenantCfg2.(*TenantConfig)
	tc.IdentityID = []byte("0x123456789")
	_, err = h.CreateTenant(context.Background(), tenantCfg1.CreateProtobuf())
	assert.Nil(t, err)
	_, err = h.CreateTenant(context.Background(), tc.CreateProtobuf())
	assert.Nil(t, err)

	resp, err := h.GetAllTenants(context.Background(), nil)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(resp.Data))
}

func TestGrpcHandler_CreateConfig(t *testing.T) {
	idService := &testingcommons.MockIDService{}
	repo, _, err := getRandomStorage()
	assert.Nil(t, err)
	repo.RegisterConfig(&NodeConfig{})
	svc := DefaultService(repo, idService)
	h := GRPCHandler(svc)
	nodeCfg := NewNodeConfig(cfg)
	_, err = h.CreateConfig(context.Background(), nodeCfg.CreateProtobuf())
	assert.Nil(t, err)

	// Already exists
	_, err = h.CreateConfig(context.Background(), nodeCfg.CreateProtobuf())
	assert.NotNil(t, err)
}

func TestGrpcHandler_CreateTenant(t *testing.T) {
	idService := &testingcommons.MockIDService{}
	repo, _, err := getRandomStorage()
	assert.Nil(t, err)
	repo.RegisterTenant(&TenantConfig{})
	svc := DefaultService(repo, idService)
	h := GRPCHandler(svc)
	nodeCfg, err := NewTenantConfig("main", cfg)
	assert.Nil(t, err)
	_, err = h.CreateTenant(context.Background(), nodeCfg.CreateProtobuf())
	assert.Nil(t, err)

	// Already exists
	_, err = h.CreateTenant(context.Background(), nodeCfg.CreateProtobuf())
	assert.NotNil(t, err)
}

func TestGrpcHandler_GenerateTenant(t *testing.T) {
	s := MockService{}
	t1, _ := NewTenantConfig(cfg.GetEthereumDefaultAccountName(), cfg)
	s.On("GenerateTenant").Return(t1, nil)
	h := GRPCHandler(s)
	tc, err := h.GenerateTenant(context.Background(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, tc)
}

func TestGrpcHandler_UpdateConfig(t *testing.T) {
	idService := &testingcommons.MockIDService{}
	repo, _, err := getRandomStorage()
	assert.Nil(t, err)
	repo.RegisterConfig(&NodeConfig{})
	svc := DefaultService(repo, idService)
	h := GRPCHandler(svc)
	nodeCfg := NewNodeConfig(cfg)

	// Config doesn't exist
	_, err = h.UpdateConfig(context.Background(), nodeCfg.CreateProtobuf())
	assert.NotNil(t, err)

	_, err = h.CreateConfig(context.Background(), nodeCfg.CreateProtobuf())
	assert.Nil(t, err)
	n := nodeCfg.(*NodeConfig)
	n.NetworkString = "other"
	_, err = h.UpdateConfig(context.Background(), n.CreateProtobuf())
	assert.Nil(t, err)

	readCfg, err := h.GetConfig(context.Background(), nil)
	assert.Nil(t, err)
	assert.Equal(t, n.GetNetworkString(), readCfg.Network)
}

func TestGrpcHandler_UpdateTenant(t *testing.T) {
	idService := &testingcommons.MockIDService{}
	repo, _, err := getRandomStorage()
	assert.Nil(t, err)
	repo.RegisterTenant(&TenantConfig{})
	svc := DefaultService(repo, idService)
	h := GRPCHandler(svc)
	nodeCfg, err := NewTenantConfig("main", cfg)
	assert.Nil(t, err)

	tid, err := nodeCfg.GetIdentityID()
	assert.Nil(t, err)

	tc := nodeCfg.(*TenantConfig)

	// Config doesn't exist
	_, err = h.UpdateTenant(context.Background(), &configpb.UpdateTenantRequest{Identifier: hexutil.Encode(tid), Data: nodeCfg.CreateProtobuf()})
	assert.NotNil(t, err)

	_, err = h.CreateTenant(context.Background(), nodeCfg.CreateProtobuf())
	assert.Nil(t, err)
	tc.EthereumDefaultAccountName = "other"
	_, err = h.UpdateTenant(context.Background(), &configpb.UpdateTenantRequest{Identifier: hexutil.Encode(tid), Data: tc.CreateProtobuf()})
	assert.Nil(t, err)

	readCfg, err := h.GetTenant(context.Background(), &configpb.GetTenantRequest{Identifier: hexutil.Encode(tid)})
	assert.Nil(t, err)
	assert.Equal(t, tc.EthereumDefaultAccountName, readCfg.EthDefaultAccountName)
}

func TestGrpcHandler_DeleteConfig(t *testing.T) {
	idService := &testingcommons.MockIDService{}
	repo, _, err := getRandomStorage()
	assert.Nil(t, err)
	repo.RegisterConfig(&NodeConfig{})
	svc := DefaultService(repo, idService)
	h := GRPCHandler(svc)

	//No error when no config
	_, err = h.DeleteConfig(context.Background(), nil)
	assert.Nil(t, err)

	nodeCfg := NewNodeConfig(cfg)
	_, err = h.CreateConfig(context.Background(), nodeCfg.CreateProtobuf())
	assert.Nil(t, err)
	_, err = h.DeleteConfig(context.Background(), nil)
	assert.Nil(t, err)

	readCfg, err := h.GetConfig(context.Background(), nil)
	assert.NotNil(t, err)
	assert.Nil(t, readCfg)
}

func TestGrpcHandler_DeleteTenant(t *testing.T) {
	idService := &testingcommons.MockIDService{}
	repo, _, err := getRandomStorage()
	assert.Nil(t, err)
	repo.RegisterTenant(&TenantConfig{})
	svc := DefaultService(repo, idService)
	h := GRPCHandler(svc)

	//No error when no config
	_, err = h.DeleteTenant(context.Background(), &configpb.GetTenantRequest{Identifier: "0x12345678"})
	assert.Nil(t, err)

	nodeCfg, err := NewTenantConfig("main", cfg)
	assert.Nil(t, err)
	_, err = h.CreateTenant(context.Background(), nodeCfg.CreateProtobuf())
	assert.Nil(t, err)
	tid, err := nodeCfg.GetIdentityID()
	assert.Nil(t, err)
	_, err = h.DeleteTenant(context.Background(), &configpb.GetTenantRequest{Identifier: hexutil.Encode(tid)})
	assert.Nil(t, err)

	readCfg, err := h.GetTenant(context.Background(), &configpb.GetTenantRequest{Identifier: hexutil.Encode(tid)})
	assert.NotNil(t, err)
	assert.Nil(t, readCfg)
}
