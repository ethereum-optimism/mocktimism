package proposer

// ProposerConfig defines the configuration for the Proposer service
type ProposerConfig struct {
	CommandFlags CommandFlags
	Ports        Ports
}

// CommandFlags contains the command-line arguments for the Proposer service
type CommandFlags struct {
	L1EthRPC          string
	RollupRPC         string
	PollInterval      string
	NumConfirmations  int
	Mnemonic          string
	L2OutputHDPath    string
	L2OOAddress       string
	PProfEnabled      bool
	MetricsEnabled    bool
	AllowNonFinalized bool
	RPCEnableAdmin    bool
}

// Ports defines the port mappings for the Proposer
type Ports struct {
	PProfPort   int
	MetricsPort int
	RPCPort     int
}
