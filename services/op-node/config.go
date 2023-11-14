// TODO generalize most of this boilerplate into a common service interface
package opnode

type OpNodeConfig struct {
	CommandFlags CommandFlags
	Host         string
	Ports        Ports
	FilePaths    FilePaths
}

type CommandFlags struct {
	L1               string
	L2               string
	JWTSecret        string
	SequencerEnabled bool
	SequencerL1Confs int
	VerifierL1Confs  int
	P2PSequencerKey  string
	RollupConfig     string
	RPCAddr          string
	RPCPort          string
	P2PListenIP      string
	P2PListenTCP     int
	P2PListenUDP     int
	P2PScoringPeers  string
	P2PBanPeers      bool
	SnapshotLogFile  string
	P2PPrivPath      string
	MetricsEnabled   bool
	MetricsAddr      string
	MetricsPort      int
	PProfEnabled     bool
	RPCEnableAdmin   bool
}

// Ports defines the port mappings for the OpNode
type Ports struct {
	RPCPort     int
	P2PPort     int
	MetricsPort int
	PProfPort   int
}

// FilePaths contains paths to various configuration files needed by OpNode
type FilePaths struct {
	SequencerKeyPath string
	NodeKeyPath      string
	JWTSecretPath    string
	RollupConfigPath string
}
