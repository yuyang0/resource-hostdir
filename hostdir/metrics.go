package hostdir

import (
	"context"

	"github.com/mitchellh/mapstructure"
	plugintypes "github.com/projecteru2/core/resource/plugins/types"
)

// GetMetricsDescription .
func (p Plugin) GetMetricsDescription(context.Context) (*plugintypes.GetMetricsDescriptionResponse, error) {
	resp := &plugintypes.GetMetricsDescriptionResponse{}
	return resp, mapstructure.Decode([]map[string]any{
		// {
		// 	"name":   "hostdir_used",
		// 	"help":   "node used host directories.",
		// 	"type":   "gauge",
		// 	"labels": []string{"podname", "nodename"},
		// },
	}, resp)
}

// GetMetrics .
func (p Plugin) GetMetrics(ctx context.Context, podname, nodename string) (*plugintypes.GetMetricsResponse, error) { //nolint
	// safeNodename := strings.ReplaceAll(nodename, ".", "_")
	metrics := []map[string]any{
		// {
		// 	"name":   "hostdir_used",
		// 	"labels": []string{podname, nodename},
		// 	"value":  fmt.Sprintf("%+v", 0),
		// 	"key":    fmt.Sprintf("core.node.%s.hostdir.used", safeNodename),
		// },
	}

	resp := &plugintypes.GetMetricsResponse{}
	return resp, mapstructure.Decode(metrics, resp)
}
