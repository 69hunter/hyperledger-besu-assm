package localstorage

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/69hunter/hyperledger-besu-assm/core"
)

const downloadURL string = "https://hyperledger.jfrog.io/artifactory/besu-binaries/besu/24.1.0/besu-24.1.0.tar.gz"

const encodeFilePath string = "/tmp/toEncode.json"
const extraDataFilePath string = "/tmp/extraData.txt"

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

func (a *Adapter) GetConfigToml() *core.ConfigToml {
	return &a.configToml
}

func (a *Adapter) GetConfigTomlInPlainText() string {
	return fmt.Sprintf(`data-path="%s"\n`, a.configToml.DataPath) +
		fmt.Sprintf(`bootnodes=[%s]\n`, `"`+strings.Join(a.configToml.BootNodes, `","`)+`"`) +
		fmt.Sprintf(`rpc-http-enabled=%t\n`, a.configToml.RpcHttpEnabled) +
		fmt.Sprintf(`rpc-http-api=[%s]\n`, `"`+strings.Join(a.configToml.RpcHttpApi, `","`)+`"`) +
		fmt.Sprintf(`rpc-http-port=%d\n`, a.configToml.RpcHttpPort) +
		fmt.Sprintf(`p2p-host="%s"\n`, a.configToml.P2pHost) +
		fmt.Sprintf(`p2p-port=%d\n`, a.configToml.P2pPort) +
		fmt.Sprintf(`host-allowlist=[%s]\n`, `"`+strings.Join(a.configToml.HostAllowlist, `","`)+`"`) +
		fmt.Sprintf(`rpc-http-cors-origins=[%s]\n`, `"`+strings.Join(a.configToml.RpcHttpCorsOrigins, `","`)+`"`) +
		fmt.Sprintf(`genesis-file="%s"\n`, a.configToml.GenesisFile) +
		fmt.Sprintf(`min-gas-price="%d"\n`, a.configToml.MinGasPrice)
}

func (a *Adapter) GetGenesisInJson() string {
	content, _ := json.MarshalIndent(a.genesis, "", " ")
	return string(content)
}

func (a *Adapter) SetGenesis(genesis core.Genesis) {
	a.genesis = genesis
}

func (a *Adapter) GetBesu() error {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("could not download besu err=%w", err)
	}

	defer resp.Body.Close()

	// read tar stream
	uncompressedStream, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("could not uncompress tar err=%w", err)
	}

	// untar release package
	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("ExtractTarGz: Next() failed err=%w", err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir("/tmp/"+header.Name, 0755); err != nil && !os.IsExist(err) {
				return fmt.Errorf("ExtractTarGz: Mkdir() failed err=, %w", err)
			}
		case tar.TypeReg:
			outFile, err := os.Create("/tmp/" + header.Name)
			if err != nil {
				return fmt.Errorf("ExtractTarGz: Create() failed: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("ExtractTarGz: Copy() failed: %w", err)
			}

			err = outFile.Chmod(fs.ModePerm)
			if err != nil {
				log.Println("could not set permissions on hyperledger besu binary")
			}

			_ = outFile.Close()

		default:
			return fmt.Errorf(
				"ExtractTarGz: unknown type: %s in %s",
				string(header.Typeflag),
				header.Name)
		}
	}

	return nil
}

func (a *Adapter) RunBesuCmd(args []string) error {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("/tmp/besu-24.1.0/bin/besu", args...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return fmt.Errorf("could not run besu command err=%w", err)
	}

	return nil
}

func (a *Adapter) createToEncodeJson() error {
	validators := []string{}
	for _, nodeInfo := range a.AllNodesInfo {
		validators = append(validators, fmt.Sprint(nodeInfo.ValidatorAddress[2:]))
	}

	_, err := os.Create(encodeFilePath)
	if err != nil {
		return fmt.Errorf("could not create toEncode.json file err=%w", err)
	}

	content, err := json.Marshal(validators)
	if err != nil {
		return err
	}
	err = os.WriteFile(encodeFilePath, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (a *Adapter) PopulateGenesisExtraData() error {
	if err := a.createToEncodeJson(); err != nil {
		return err
	}

	cmd := []string{"rlp", "encode", fmt.Sprintf("--from=%s", encodeFilePath), fmt.Sprintf("--to=%s", extraDataFilePath), "--type=IBFT_EXTRA_DATA"}
	if err := a.RunBesuCmd(cmd); err != nil {
		return err
	}

	content, err := os.ReadFile(extraDataFilePath)
	if err != nil {
		return err
	}
	extraData := string(content)
	a.genesis.ExtraData = extraData
	return nil
}
