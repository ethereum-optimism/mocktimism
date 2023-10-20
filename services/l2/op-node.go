package opnode

import (
	"github.com/ethereum-optimism/mocktimism/service-discovery"
	"github.com/urfave/cli/v2"
)

type RollupNodeService struct {
	ver        string
	ctx        *cli.Context
	serviceCfg servicediscovery.ServiceConfig
}

func (r *RollupNodeService) Hostname() string {
	return "myRollupNodeHostname"
}

func (r *RollupNodeService) Port() int {
	return 8080
}

func (r *RollupNodeService) ServiceType() string {
	return "_rollupnode._tcp"
}

func (r *RollupNodeService) ID() string {
	return "rollupNode123"
}

func (r *RollupNodeService) Config() servicediscovery.ServiceConfig {
	return r.serviceCfg
}

func (r *RollupNodeService) Start() error {
	return rollupNodeMainFunc(r.ver, r.ctx)
}

func NewRollupNodeService(ver string, ctx *cli.Context, config servicediscovery.ServiceConfig) *RollupNodeService {
	return &RollupNodeService{
		ver:        ver,
		ctx:        ctx,
		serviceCfg: config,
	}
}

func rollupNodeMainFunc(ver string, ctx *cli.Context) error {
	return nil
}
