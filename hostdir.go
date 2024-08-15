package main

import (
	"context"
	"fmt"
	"os"

	"github.com/projecteru2/core/resource/plugins"
	coretypes "github.com/projecteru2/core/types"

	"github.com/yuyang0/resource-hostdir/cmd"
	"github.com/yuyang0/resource-hostdir/cmd/calculate"
	"github.com/yuyang0/resource-hostdir/cmd/hostdir"
	"github.com/yuyang0/resource-hostdir/cmd/metrics"
	"github.com/yuyang0/resource-hostdir/cmd/node"
	hostdirlib "github.com/yuyang0/resource-hostdir/hostdir"
	"github.com/yuyang0/resource-hostdir/version"

	"github.com/urfave/cli/v2"
)

func NewPlugin(ctx context.Context, config coretypes.Config) (plugins.Plugin, error) {
	p, err := hostdirlib.NewPlugin(ctx, config, nil)
	return p, err
}

func main() {
	cli.VersionPrinter = func(_ *cli.Context) {
		fmt.Print(version.String())
	}

	app := cli.NewApp()
	app.Name = version.NAME
	app.Usage = "Run eru resource hostdir plugin"
	app.Version = version.VERSION
	app.Commands = []*cli.Command{
		hostdir.Name(),
		metrics.Description(),
		metrics.GetMetrics(),

		node.AddNode(),
		node.RemoveNode(),
		node.GetNodesDeployCapacity(),
		node.SetNodeResourceCapacity(),
		node.GetNodeResourceInfo(),
		node.SetNodeResourceInfo(),
		node.SetNodeResourceUsage(),
		node.GetMostIdleNode(),
		node.FixNodeResource(),

		calculate.CalculateDeploy(),
		calculate.CalculateRealloc(),
		calculate.CalculateRemap(),
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Value:       "hostdir.yaml",
			Usage:       "config file path for plugin, in yaml",
			Destination: &cmd.ConfigPath,
			EnvVars:     []string{"ERU_RESOURCE_CONFIG_PATH"},
		},
		&cli.BoolFlag{
			Name:        "embedded-storage",
			Usage:       "active embedded storage",
			Destination: &cmd.EmbeddedStorage,
		},
	}
	_ = app.Run(os.Args)
}
