package node

// This was originally from TODO add
import (
	"context"
	"errors"
	"net"
	"strconv"

	opnode "github.com/ethereum-optimism/optimism/op-node"
	"github.com/ethereum-optimism/optimism/op-node/chaincfg"
	"github.com/ethereum-optimism/optimism/op-node/client"
	"github.com/ethereum-optimism/optimism/op-node/flags"
	"github.com/ethereum-optimism/optimism/op-node/heartbeat"
	"github.com/ethereum-optimism/optimism/op-node/metrics"
	"github.com/ethereum-optimism/optimism/op-node/node"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/sources"
	"github.com/ethereum-optimism/optimism/op-node/version"
	opservice "github.com/ethereum-optimism/optimism/op-service"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum-optimism/optimism/op-service/opio"
	oppprof "github.com/ethereum-optimism/optimism/op-service/pprof"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"

	"github.com/mdehoog/op-anvil/util"
)

func RollupNodeMain(ver string) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		log.Info("Initializing Rollup Node")
		logCfg := oplog.ReadCLIConfig(ctx)
		log := oplog.NewLogger(oplog.AppOut(ctx), logCfg)
		oplog.SetGlobalLogHandler(log.GetHandler())
		opservice.ValidateEnvVars(flags.EnvVarPrefix, flags.Flags, log)
		m := metrics.NewMetrics("default")

		cfg, err := opnode.NewConfig(ctx, log)
		if err != nil {
			log.Error("Unable to create the rollup node config", "error", err)
			return err
		}

		cl1, _, err := cfg.L1.Setup(context.Background(), log, &cfg.Rollup)
		if err != nil {
			return err
		}
		cfg.L1 = &node.PreparedL1Endpoint{
			Client:          util.NewAnvilClient(cl1),
			TrustRPC:        true,
			RPCProviderKind: sources.RPCKindBasic,
		}

		cl2, _, err := cfg.L2.Setup(context.Background(), log, &cfg.Rollup)
		if err != nil {
			return err
		}
		cfg.L2 = &PreparedL2Endpoint{
			Client: util.NewAnvilClient(cl2),
		}

		snapshotLog, err := opnode.NewSnapshotLogger(ctx)
		if err != nil {
			log.Error("Unable to create snapshot root logger", "error", err)
			return err
		}

		// Only pretty-print the banner if it is a terminal log. Other log it as key-value pairs.
		if logCfg.Format == "terminal" {
			log.Info("rollup config:\n" + cfg.Rollup.Description(chaincfg.L2ChainIDToNetworkDisplayName))
		} else {
			cfg.Rollup.LogDescription(log, chaincfg.L2ChainIDToNetworkDisplayName)
		}

		n, err := node.New(context.Background(), cfg, log, snapshotLog, ver, m)
		if err != nil {
			log.Error("Unable to create the rollup node", "error", err)
			return err
		}
		log.Info("Starting rollup node", "version", ver)

		if err := n.Start(context.Background()); err != nil {
			log.Error("Unable to start rollup node", "error", err)
			return err
		}
		defer n.Close()

		m.RecordInfo(ver)
		m.RecordUp()
		log.Info("Rollup node started")

		if cfg.Heartbeat.Enabled {
			var peerID string
			if cfg.P2P.Disabled() {
				peerID = "disabled"
			} else {
				peerID = n.P2P().Host().ID().String()
			}

			beatCtx, beatCtxCancel := context.WithCancel(context.Background())
			payload := &heartbeat.Payload{
				Version: version.Version,
				Meta:    version.Meta,
				Moniker: cfg.Heartbeat.Moniker,
				PeerID:  peerID,
				ChainID: cfg.Rollup.L2ChainID.Uint64(),
			}
			go func() {
				if err := heartbeat.Beat(beatCtx, log, cfg.Heartbeat.URL, payload); err != nil {
					log.Error("heartbeat goroutine crashed", "err", err)
				}
			}()
			defer beatCtxCancel()
		}

		if cfg.Pprof.Enabled {
			pprofCtx, pprofCancel := context.WithCancel(context.Background())
			go func() {
				log.Info("pprof server started", "addr", net.JoinHostPort(cfg.Pprof.ListenAddr, strconv.Itoa(cfg.Pprof.ListenPort)))
				if err := oppprof.ListenAndServe(pprofCtx, cfg.Pprof.ListenAddr, cfg.Pprof.ListenPort); err != nil {
					log.Error("error starting pprof", "err", err)
				}
			}()
			defer pprofCancel()
		}

		opio.BlockOnInterrupts()

		return nil
	}
}

type PreparedL2Endpoint struct {
	Client client.RPC
}

var _ node.L2EndpointSetup = (*PreparedL2Endpoint)(nil)

func (p *PreparedL2Endpoint) Setup(_ context.Context, _ log.Logger, rollupCfg *rollup.Config) (client.RPC, *sources.EngineClientConfig, error) {
	return p.Client, sources.EngineClientDefaultConfig(rollupCfg), nil
}

func (p *PreparedL2Endpoint) Check() error {
	if p.Client == nil {
		return errors.New("rpc client cannot be nil")
	}

	return nil
}
