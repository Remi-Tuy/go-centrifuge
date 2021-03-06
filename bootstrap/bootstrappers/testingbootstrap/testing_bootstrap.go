// +build integration

package testingbootstrap

import (
	"github.com/centrifuge/go-centrifuge/anchors"
	"github.com/centrifuge/go-centrifuge/bootstrap"
	"github.com/centrifuge/go-centrifuge/bootstrap/bootstrappers/testlogging"
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
	"github.com/centrifuge/go-centrifuge/p2p"
	"github.com/centrifuge/go-centrifuge/queue"
	"github.com/centrifuge/go-centrifuge/storage/leveldb"
	"github.com/centrifuge/go-centrifuge/testingutils"
	"github.com/centrifuge/go-centrifuge/transactions"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("context")

var bootstappers = []bootstrap.TestBootstrapper{
	&testlogging.TestLoggingBootstrapper{},
	&config.Bootstrapper{},
	&leveldb.Bootstrapper{},
	transactions.Bootstrapper{},
	ethereum.Bootstrapper{},
	&queue.Bootstrapper{},
	&ethid.Bootstrapper{},
	&configstore.Bootstrapper{},
	anchors.Bootstrapper{},
	documents.Bootstrapper{},
	p2p.Bootstrapper{},
	&genericdoc.Bootstrapper{},
	&invoice.Bootstrapper{},
	&purchaseorder.Bootstrapper{},
	coredocument.Bootstrapper{},
	documents.PostBootstrapper{},
	&nft.Bootstrapper{},
	&queue.Starter{},
}

func TestFunctionalEthereumBootstrap() map[string]interface{} {
	cm := testingutils.BuildIntegrationTestingContext()
	for _, b := range bootstappers {
		err := b.TestBootstrap(cm)
		if err != nil {
			log.Error("Error encountered while bootstrapping", err)
			panic(err)
		}
	}

	return cm
}
func TestFunctionalEthereumTearDown() {
	for _, b := range bootstappers {
		err := b.TestTearDown()
		if err != nil {
			log.Error("Error encountered while bootstrapping", err)
			panic(err)
		}
	}
}
