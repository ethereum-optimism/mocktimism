package batcher

// OpBatcherConfig defines the configuration for the OpBatcher service
type OpBatcherConfig struct {
	CommandFlags CommandFlags
	Ports        Ports
}

// CommandFlags contains the command-line arguments for the OpBatcher service
type CommandFlags struct {
	L1EthRPC           string
	L2EthRPC           string
	RollupRPC          string
	MaxChannelDuration int
	SubSafetyMargin    int
	PollInterval       string
	NumConfirmations   int
	Mnemonic           string
	SequencerHDPath    string
	PProfEnabled       bool
	MetricsEnabled     bool
	RPCEnableAdmin     bool
	BatchType          int
}

// Ports defines the port mappings for the OpBatcher
type Ports struct {
	PProfPort   int
	MetricsPort int
	RPCPort     int
}
