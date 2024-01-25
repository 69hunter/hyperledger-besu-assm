package app

import (
	"fmt"

	"github.com/69hunter/hyperledger-besu-assm/core"
	"github.com/69hunter/hyperledger-besu-assm/localstorage"
	"github.com/69hunter/hyperledger-besu-assm/s3storage"
)

type Adapter struct {
	core         *core.Core
	s3Storage    *s3storage.Adapter
	localStorage *localstorage.Adapter
}

func (a *Adapter) LambdaHandler(request core.Core) (string, error) {
	if err := a.localStorage.AddNodeInfo(request.NodeInfo); err != nil {
		return "", fmt.Errorf("could not save node info in local storage: %w", err)
	}

	if request.Setup.Total == len(a.localStorage.AllNodesInfo) {
		// fetch hyperledger besu binary
		if err := a.localStorage.GetBesu(); err != nil {
			return "", fmt.Errorf("could not set local hyperledger besu binary err=%w", err)
		}

		s3Api, err := s3storage.NewAdapter(request.Setup.AWSRegion, request.Setup.S3BucketName)
		if err != nil {
			return "", fmt.Errorf("could not create new S3 adapter err=%w", err)
		}
		a.s3Storage = s3Api

		// Create & upload config.toml file
		bootnodes := []string{}
		for _, nodeInfo := range a.localStorage.AllNodesInfo {
			bootnodes = append(bootnodes, fmt.Sprintf("enode://%s@%s:%d", nodeInfo.NodePublicKey, nodeInfo.NodeHost, nodeInfo.NodePort))
		}

		request.ConfigToml.BootNodes = bootnodes
		a.localStorage.SetConfigToml(request.ConfigToml)

		if err := a.s3Storage.WriteData("config.toml", a.localStorage.GetConfigTomlInPlainText()); err != nil {
			return "", fmt.Errorf("could not write config.toml to S3 err=%w", err)
		}

		// Create & upload genesis.json file
		a.localStorage.SetGenesis(request.Genesis)
		if err := a.localStorage.PopulateGenesisExtraData(); err != nil {
			return "", fmt.Errorf("could not create genesis extra data err=%w", err)
		}

		if err := a.s3Storage.WriteData("genesis.json", a.localStorage.GetGenesisInJson()); err != nil {
			return "", fmt.Errorf("could not write genesis.json to S3 err=%w", err)
		}

		return "config.toml & genesis.json files successfully created and uploaded to S3", nil
	}

	return "Node information successfully saved", nil
}

func NewAdapter(core *core.Core, localStorage *localstorage.Adapter) *Adapter {
	return &Adapter{
		core:         core,
		localStorage: localStorage,
	}
}
