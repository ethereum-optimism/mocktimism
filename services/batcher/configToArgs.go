package batcher

import (
	"fmt"
)

func BuildCommandArgs(config BatcherConfig) []string {
	args := []string{
		"--l1=" + config.CommandFlags.L1,
		"--l2=" + config.CommandFlags.L2,
		"--l2.jwt-secret=" + config.CommandFlags.JWTSecret,
	}

	if config.CommandFlags.SequencerEnabled {
		args = append(args, "--sequencer.enabled")
	}

	args = append(args, "--sequencer.l1-confs="+fmt.Sprintf("%d", config.CommandFlags.SequencerL1Confs))
	args = append(args, "--verifier.l1-confs="+fmt.Sprintf("%d", config.CommandFlags.VerifierL1Confs))
	args = append(args, "--p2p.sequencer.key="+config.CommandFlags.P2PSequencerKey)
	args = append(args, "--rollup.config="+config.CommandFlags.RollupConfig)
	args = append(args, "--rpc.addr="+config.CommandFlags.RPCAddr)
	args = append(args, "--rpc.port="+config.CommandFlags.RPCPort)
	args = append(args, "--p2p.listen.ip="+config.CommandFlags.P2PListenIP)
	args = append(args, "--p2p.listen.tcp="+fmt.Sprintf("%d", config.CommandFlags.P2PListenTCP))
	args = append(args, "--p2p.listen.udp="+fmt.Sprintf("%d", config.CommandFlags.P2PListenUDP))
	args = append(args, "--p2p.scoring.peers="+config.CommandFlags.P2PScoringPeers)

	if config.CommandFlags.P2PBanPeers {
		args = append(args, "--p2p.ban.peers=true")
	} else {
		args = append(args, "--p2p.ban.peers=false")
	}

	args = append(args, "--snapshotlog.file="+config.CommandFlags.SnapshotLogFile)
	args = append(args, "--p2p.priv.path="+config.CommandFlags.P2PPrivPath)

	if config.CommandFlags.MetricsEnabled {
		args = append(args, "--metrics.enabled")
	}

	args = append(args, "--metrics.addr="+config.CommandFlags.MetricsAddr)
	args = append(args, "--metrics.port="+fmt.Sprintf("%d", config.CommandFlags.MetricsPort))

	if config.CommandFlags.PProfEnabled {
		args = append(args, "--pprof.enabled")
	}

	if config.CommandFlags.RPCEnableAdmin {
		args = append(args, "--rpc.enable-admin")
	}

	return args
}
