package core

type Core struct {
	Setup      `json:"setup,omitempty"`
	Genesis    `json:"genesis,omitempty"`
	ConfigToml `json:"config_toml,omitempty"`
	NodeInfo   `json:"node_info,omitempty"`
}

type ConfigToml struct {
	DataPath           string   `json:"data-path"`
	BootNodes          []string `json:"bootnodes"`
	RpcHttpEnabled     bool     `json:"rpc-http-enabled"`
	RpcHttpApi         []string `json:"rpc-http-api"`
	RpcHttpPort        int      `json:"rpc-http-port"`
	P2pHost            string   `json:"p2p-host"`
	P2pPort            int      `json:"p2p-port"`
	HostAllowlist      []string `json:"host-allowlist"`
	RpcHttpCorsOrigins []string `json:"rpc-http-cors-origins"`
	GenesisFile        string   `json:"genesis-file"`
	MinGasPrice        int      `json:"min-gas-price"`
}

type NodeInfo struct {
	ValidatorAddress string `json:"validator_address"`
	NodePublicKey    string `json:"node_public_key"`
	NodeHost         string `json:"node_host"`
	NodePort         int    `json:"node_port"`
}

type Setup struct {
	AWSRegion    string `json:"aws_region"`
	S3BucketName string `json:"s3_bucket_name"`
	Total        int    `json:"total"`
}

type Genesis struct {
	Config     Config   `json:"config"`
	Nonce      string   `json:"nonce"`
	Timestamp  string   `json:"timestamp"`
	ExtraData  string   `json:"extraData"`
	GasLimit   string   `json:"gasLimit"`
	Difficulty string   `json:"difficulty"`
	MixHash    string   `json:"mixHash"`
	Alloc      AllocMap `json:"alloc"`
}

type Config struct {
	ChainId     int   `json:"chainId"`
	BerlinBlock int   `json:"berlinBlock"`
	Ibft2       Ibft2 `json:"ibft2"`
}

type Ibft2 struct {
	BlockPeriodSeconds    int    `json:"blockperiodseconds"`
	EpochLength           int    `json:"epochlength"`
	RequestTimeoutSeconds int    `json:"requesttimeoutseconds"`
	BlockReward           string `json:"blockreward"`
	MiningBeneficiary     string `json:"miningbeneficiary"`
}

type AllocMap map[string]Allocation

type Allocation struct {
	Balance string `json:"balance"`
}

func (c *Core) GetCore() *Core {
	return c
}

func NewAdapter() *Core {
	return &Core{}
}
