package opnode

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"

	"github.com/ethereum/go-ethereum/log"
)

var (
	SERVICE_TYPE = "op-node"
	DEPENDS_ON   = []string{
		"l1",
		"l2",
	}
)

type OpNodeService struct {
	id     string
	config OpNodeConfig
	cmd    *exec.Cmd
	logger log.Logger
}

func NewOpNodeService(id string, logger log.Logger, cfg OpNodeConfig) (*OpNodeService, error) {
	return &OpNodeService{
		id:     id,
		config: cfg,
		logger: logger,
	}, nil
}

func (s *OpNodeService) Hostname() string {
	return s.config.Host
}

// TODO get rid of this from service interface since some services have many ports
func (s *OpNodeService) Port() int {
	return int(s.config.Ports.RPCPort)
}

func (s *OpNodeService) ServiceType() string {
	return SERVICE_TYPE
}

func (s *OpNodeService) ID() string {
	return s.id
}

func (s *OpNodeService) Config() interface{} {
	return s.config
}

func (s *OpNodeService) HealthCheck() (bool, error) {
	// TODO hit healthz for healthcheck
	return true, nil
}

func (s *OpNodeService) Start(ctx context.Context) error {
	args := s.buildCommandArgs()
	s.cmd = exec.CommandContext(ctx, "op-node", args...)

	stdout, _ := s.cmd.StdoutPipe()
	stderr, _ := s.cmd.StderrPipe()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			s.logger.Info(scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			s.logger.Error(scanner.Text())
		}
	}()

	if err := s.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start OpNode service: %w", err)
	}

	s.logger.Info("message", "OpNode service started", "id", s.id)
	return nil
}

func (s *OpNodeService) Stop() error {
	if s.cmd != nil && s.cmd.Process != nil {
		if err := s.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to stop OpNode service: %w", err)
		}
		s.logger.Info("message", "OpNode service stopped", "id", s.id)
	}
	return nil
}

func (s *OpNodeService) buildCommandArgs() []string {
	args := []string{
		"--l1=" + s.config.CommandFlags.L1,
		"--l2=" + s.config.CommandFlags.L2,
		"--l2.jwt-secret=" + s.config.CommandFlags.JWTSecret,
	}

	if s.config.CommandFlags.SequencerEnabled {
		args = append(args, "--sequencer.enabled")
	}

	args = append(args, "--sequencer.l1-confs="+fmt.Sprintf("%d", s.config.CommandFlags.SequencerL1Confs))
	args = append(args, "--verifier.l1-confs="+fmt.Sprintf("%d", s.config.CommandFlags.VerifierL1Confs))
	args = append(args, "--p2p.sequencer.key="+s.config.CommandFlags.P2PSequencerKey)
	args = append(args, "--rollup.config="+s.config.CommandFlags.RollupConfig)
	args = append(args, "--rpc.addr="+s.config.CommandFlags.RPCAddr)
	args = append(args, "--rpc.port="+s.config.CommandFlags.RPCPort)
	args = append(args, "--p2p.listen.ip="+s.config.CommandFlags.P2PListenIP)
	args = append(args, "--p2p.listen.tcp="+fmt.Sprintf("%d", s.config.CommandFlags.P2PListenTCP))
	args = append(args, "--p2p.listen.udp="+fmt.Sprintf("%d", s.config.CommandFlags.P2PListenUDP))
	args = append(args, "--p2p.scoring.peers="+s.config.CommandFlags.P2PScoringPeers)

	if s.config.CommandFlags.P2PBanPeers {
		args = append(args, "--p2p.ban.peers=true")
	} else {
		args = append(args, "--p2p.ban.peers=false")
	}

	args = append(args, "--snapshotlog.file="+s.config.CommandFlags.SnapshotLogFile)
	args = append(args, "--p2p.priv.path="+s.config.CommandFlags.P2PPrivPath)

	if s.config.CommandFlags.MetricsEnabled {
		args = append(args, "--metrics.enabled")
	}

	args = append(args, "--metrics.addr="+s.config.CommandFlags.MetricsAddr)
	args = append(args, "--metrics.port="+fmt.Sprintf("%d", s.config.CommandFlags.MetricsPort))

	if s.config.CommandFlags.PProfEnabled {
		args = append(args, "--pprof.enabled")
	}

	if s.config.CommandFlags.RPCEnableAdmin {
		args = append(args, "--rpc.enable-admin")
	}

	return args
}
