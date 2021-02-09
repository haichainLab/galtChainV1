package main

import (
	_ "net/http/pprof"

	"github.com/samoslab/galtcoin/src/skycoin"
	"github.com/samoslab/galtcoin/src/util/logging"
	"github.com/samoslab/galtcoin/src/visor"
)

var (
	// Version of the node. Can be set by -ldflags
	Version = "1.24.1"
	// Commit ID. Can be set by -ldflags
	Commit = ""
	// Branch name. Can be set by -ldflags
	Branch = ""
	// ConfigMode (possible values are "", "STANDALONE_CLIENT").
	// This is used to change the default configuration.
	// Can be set by -ldflags
	ConfigMode = ""

	logger = logging.MustGetLogger("main")

	// GenesisSignatureStr hex string of genesis signature
	GenesisSignatureStr = "33820d76a0fb96fad643a2feb5ca6fce5b26851a9fb139187b272129a8c0ebce5bd3906754542b6c3e5da60806fb4678be96575d7bb5e221656bbdc0e332676100"
	// GenesisAddressStr genesis address string
	GenesisAddressStr = "DudoieHAEgeyr9N6CMQ3x9UB5knqScAcaM"
	// BlockchainPubkeyStr pubic key string
	BlockchainPubkeyStr = "028b216327e836adfcaa8508cb7321571d00b8fce94994b41e1f4fa62907650707"
	// BlockchainSeckeyStr empty private key string
	BlockchainSeckeyStr = ""

	// GenesisTimestamp genesis block create unix time
	GenesisTimestamp uint64 = 1541124558
	// GenesisCoinVolume represents the coin capacity
	GenesisCoinVolume uint64 = 1000000000000000

	// DefaultConnections the default trust node addresses
	DefaultConnections = []string{
		"193.112.246.101:6868",
		"150.109.100.158:6868",
		"150.109.62.225:6868",
		"47.52.211.167:6868",
	}
)

func main() {
	// get node config
	nodeConfig := skycoin.NewNodeConfig(ConfigMode, skycoin.NodeParameters{
		GenesisSignatureStr: GenesisSignatureStr,
		GenesisAddressStr:   GenesisAddressStr,
		GenesisCoinVolume:   GenesisCoinVolume,
		GenesisTimestamp:    GenesisTimestamp,
		BlockchainPubkeyStr: BlockchainPubkeyStr,
		BlockchainSeckeyStr: BlockchainSeckeyStr,
		DefaultConnections:  DefaultConnections,
		PeerListURL:         "http://download.galtcoin.io/peers.txt",
		Port:                6868,
		WebInterfacePort:    6869,
		DataDirectory:       "$HOME/.galtcoin",
		ProfileCPUFile:      "galtcoin.prof",
	})

	// create a new fiber coin instance
	coin := skycoin.NewCoin(
		skycoin.Config{
			Node: *nodeConfig,
			Build: visor.BuildInfo{
				Version: Version,
				Commit:  Commit,
				Branch:  Branch,
			},
		},
		logger,
	)

	// parse config values
	coin.ParseConfig()

	// run fiber coin node
	coin.Run()
}
