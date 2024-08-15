package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	resourcetypes "github.com/projecteru2/core/resource/types"
	"github.com/projecteru2/core/utils"
	"github.com/urfave/cli/v2"
	"github.com/yuyang0/resource-hostdir/hostdir"
)

var (
	ConfigPath      string
	EmbeddedStorage bool
)

func Serve(c *cli.Context, f func(s *hostdir.Plugin, in resourcetypes.RawParams) (interface{}, error)) error {
	cfg, err := utils.LoadConfig(ConfigPath)
	if err != nil {
		return cli.Exit(err, 128)
	}

	var t *testing.T
	if EmbeddedStorage {
		t = &testing.T{}
	}

	s, err := hostdir.NewPlugin(c.Context, cfg, t)
	if err != nil {
		return cli.Exit(err, 128)
	}

	in := resourcetypes.RawParams{}
	if err := json.NewDecoder(os.Stdin).Decode(&in); err != nil {
		fmt.Fprintf(os.Stderr, "HOSTDIR: failed decode input json: %s\n", err)
		fmt.Fprintf(os.Stderr, "HOSTDIR: input: %v\n", in)
		return cli.Exit(err, 128)
	}

	if r, err := f(s, in); err != nil {
		fmt.Fprintf(os.Stderr, "HOSTDIR: failed call function: %s\n", err)
		fmt.Fprintf(os.Stderr, "HOSTDIR: input: %v\n", in)
		return cli.Exit(err, 128)
	} else if o, err := json.Marshal(r); err != nil {
		fmt.Fprintf(os.Stderr, "HOSTDIR: failed encode return object: %s\n", err)
		fmt.Fprintf(os.Stderr, "HOSTDIR: input: %v\n", in)
		fmt.Fprintf(os.Stderr, "HOSTDIR: output: %v\n", o)
		return cli.Exit(err, 128)
	} else { //nolint
		fmt.Print(string(o))
	}
	return nil
}
