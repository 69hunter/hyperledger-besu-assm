package localstorage

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/69hunter/hyperledger-besu-assm/core"
)

type Adapter struct {
	AllNodesInfo []core.NodeInfo
	configToml   core.ConfigToml
	genesis      core.Genesis
}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) AddNodeInfo(nodeInfo core.NodeInfo) error {
	a.AllNodesInfo = append(a.AllNodesInfo, nodeInfo)
	return nil
}

func (a *Adapter) SetConfigToml(config core.ConfigToml) {
	a.configToml = config
}

func (a *Adapter) GetConfigTomlInPlainText() string {
	return fmt.Sprintf("data-path=\"%s\"\n", a.configToml.DataPath) +
		fmt.Sprintf("bootnodes=[%s]\n", `"`+strings.Join(a.configToml.BootNodes, `","`)+`"`) +
		fmt.Sprintf("rpc-http-enabled=%t\n", a.configToml.RpcHttpEnabled) +
		fmt.Sprintf("rpc-http-api=[%s]\n", `"`+strings.Join(a.configToml.RpcHttpApi, `","`)+`"`) +
		fmt.Sprintf("rpc-http-port=%d\n", a.configToml.RpcHttpPort) +
		fmt.Sprintf("p2p-host=\"%s\"\n", a.configToml.P2pHost) +
		fmt.Sprintf("p2p-port=%d\n", a.configToml.P2pPort) +
		fmt.Sprintf("host-allowlist=[%s]\n", `"`+strings.Join(a.configToml.HostAllowlist, `","`)+`"`) +
		fmt.Sprintf("rpc-http-cors-origins=[%s]\n", `"`+strings.Join(a.configToml.RpcHttpCorsOrigins, `","`)+`"`) +
		fmt.Sprintf("genesis-file=\"%s\"\n", a.configToml.GenesisFile) +
		fmt.Sprintf("min-gas-price=\"%d\"\n", a.configToml.MinGasPrice)
}

func (a *Adapter) GetGenesisInJson() string {
	content, _ := json.MarshalIndent(a.genesis, "", " ")
	return string(content)
}

func (a *Adapter) SetGenesis(genesis core.Genesis) {
	a.genesis = genesis
}
