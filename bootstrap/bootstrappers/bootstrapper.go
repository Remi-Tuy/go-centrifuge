package bootstrappers

import (
	"github.com/centrifuge/go-centrifuge/anchors"
	"github.com/centrifuge/go-centrifuge/api"
	"github.com/centrifuge/go-centrifuge/bootstrap"
	"github.com/centrifuge/go-centrifuge/config"
	"github.com/centrifuge/go-centrifuge/config/configstore"
	"github.com/centrifuge/go-centrifuge/coredocument"
	"github.com/centrifuge/go-centrifuge/documents"
	"github.com/centrifuge/go-centrifuge/documents/genericdoc"
	"github.com/centrifuge/go-centrifuge/documents/invoice"
	"github.com/centrifuge/go-centrifuge/documents/purchaseorder"
	"github.com/centrifuge/go-centrifuge/ethereum"
	"github.com/centrifuge/go-centrifuge/identity/ethid"
	"github.com/centrifuge/go-centrifuge/nft"
	"github.com/centrifuge/go-centrifuge/node"
	"github.com/centrifuge/go-centrifuge/p2p"
	"github.com/centrifuge/go-centrifuge/queue"
	"github.com/centrifuge/go-centrifuge/storage/leveldb"
	"github.com/centrifuge/go-centrifuge/transactions"
	"github.com/centrifuge/go-centrifuge/version"
	log2 "github.com/ipfs/go-log"
)

var log = log2.Logger("context")

// MainBootstrapper holds all the bootstrapper implementations
type MainBootstrapper struct {
	Bootstrappers []bootstrap.Bootstrapper
}

// PopulateBaseBootstrappers adds all the bootstrapper implementations to MainBootstrapper
func (m *MainBootstrapper) PopulateBaseBootstrappers() {
	m.Bootstrappers = []bootstrap.Bootstrapper{
		&version.Bootstrapper{},
		&config.Bootstrapper{},
		&leveldb.Bootstrapper{},
		transactions.Bootstrapper{},
		ethereum.Bootstrapper{},
		&queue.Bootstrapper{},
		&ethid.Bootstrapper{},
		&configstore.Bootstrapper{},
		&anchors.Bootstrapper{},
		documents.Bootstrapper{},
		p2p.Bootstrapper{},
		api.Bootstrapper{},
		&genericdoc.Bootstrapper{},
		&invoice.Bootstrapper{},
		&purchaseorder.Bootstrapper{},
		coredocument.Bootstrapper{},
		documents.PostBootstrapper{},
		&nft.Bootstrapper{},
	}
}

// PopulateCommandBootstrappers adds all the bootstrapper implementations relevant for one off commands
func (m *MainBootstrapper) PopulateCommandBootstrappers() {
	m.Bootstrappers = []bootstrap.Bootstrapper{
		&version.Bootstrapper{},
		&config.Bootstrapper{},
		&leveldb.Bootstrapper{},
		transactions.Bootstrapper{},
		ethereum.Bootstrapper{},
		&queue.Bootstrapper{},
		&anchors.Bootstrapper{},
		&ethid.Bootstrapper{},
	}
}

// PopulateRunBootstrappers adds blocking Node bootstrapper at the end.
// Note: Node bootstrapper must be the last bootstrapper to be invoked as it won't return until node is shutdown
func (m *MainBootstrapper) PopulateRunBootstrappers() {
	m.PopulateBaseBootstrappers()
	m.Bootstrappers = append(m.Bootstrappers, &node.Bootstrapper{})
}

// Bootstrap runs all the loaded bootstrapper implementations.
func (m *MainBootstrapper) Bootstrap(context map[string]interface{}) error {
	for _, b := range m.Bootstrappers {
		err := b.Bootstrap(context)
		if err != nil {
			log.Error("Error encountered while bootstrapping", err)
			return err
		}
	}
	return nil
}
