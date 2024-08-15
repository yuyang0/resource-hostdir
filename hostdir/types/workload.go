package types

import (
	"encoding/json"

	"github.com/cockroachdb/errors"
	resourcetypes "github.com/projecteru2/core/resource/types"
)

// WorkloadResource indicate hostdir workload resource
type WorkloadResource struct {
	Volumes VolumeBindings `json:"volumes" mapstructure:"volumes"`
}

func NewWorkloadResoure() *WorkloadResource {
	return &WorkloadResource{
		Volumes: VolumeBindings{},
	}
}

func (w *WorkloadResource) AsRawParams() resourcetypes.RawParams {
	return resourcetypes.RawParams{
		"volumes": w.Volumes,
	}
}

func (w *WorkloadResource) Size() int64 {
	return w.Volumes.TotalSize()
}

func (w *WorkloadResource) DeepCopy() *WorkloadResource {
	ans := &WorkloadResource{
		Volumes: VolumeBindings{},
	}
	for _, vb := range w.Volumes {
		ans.Volumes = append(ans.Volumes, vb.DeepCopy())
	}
	return ans
}

// ParseFromRawParams .
func (w *WorkloadResource) Parse(rawParams resourcetypes.RawParams) (err error) {
	// Have to use json because volume plan use customize marshal
	body, err := json.Marshal(rawParams)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, w)
}

// WorkloadResourceRaw includes all possible fields passed by eru-core for editing workload
// for request calculation
type WorkloadResourceRequest struct {
	Volumes VolumeBindings `json:"volumes" mapstructure:"volumes"`
}

func (w *WorkloadResourceRequest) DeepCopy() *WorkloadResourceRequest {
	ans := &WorkloadResourceRequest{}
	for _, vb := range w.Volumes {
		newVB := *vb
		ans.Volumes = append(ans.Volumes, &newVB)
	}
	return ans
}

// Validate .
func (w *WorkloadResourceRequest) Validate() error {
	return w.Volumes.Validate()
}

// Parse .
func (w *WorkloadResourceRequest) Parse(rawParams resourcetypes.RawParams) (err error) {
	if w.Volumes, err = NewVolumeBindings(rawParams.OneOfStringSlice("volumes", "volume-request", "volumes-request")); err != nil {
		return errors.Wrap(err, "failed to parse workload resource request")
	}
	return nil
}
