package types

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/projecteru2/core/utils"
)

// VolumeBinding format =>  pool/image:dst[:flags][:size][:read_IOPS:write_IOPS:read_bytes:write_bytes]
type VolumeBinding struct {
	Source      string
	Destination string `json:"destination" mapstructure:"destination"`
	SizeInBytes int64  `json:"size_in_bytes" mapstructure:"size_in_bytes"`
}

func (vb *VolumeBinding) GetSource() string {
	return vb.Source
}

func (vb *VolumeBinding) GetMapKey() [2]string {
	return [2]string{vb.Source, vb.Destination}
}

func (vb *VolumeBinding) DeepCopy() *VolumeBinding {
	return &VolumeBinding{
		Source:      vb.Source,
		Destination: vb.Destination,
		SizeInBytes: vb.SizeInBytes,
	}
}

// NewVolumeBinding returns pointer of VolumeBinding
func NewVolumeBinding(volume string) (_ *VolumeBinding, err error) {
	var (
		src, dst string
		size     int64
	)

	switch parts := strings.Split(volume, ":"); len(parts) {
	case 2:
		src, dst = parts[0], parts[1]
	case 3:
		src, dst = parts[0], parts[1]
		if size, err = utils.ParseRAMInHuman(parts[2]); err != nil {
			return nil, errors.Wrapf(ErrInvalidVolume, volume)
		}
	default:
		return nil, errors.Wrap(ErrInvalidVolume, volume)
	}

	vb := &VolumeBinding{
		Source:      src,
		Destination: dst,
		SizeInBytes: size,
	}

	return vb, vb.Validate()
}

// Validate return error if invalid
// Please note: we allow negative value for SizeInBytes,
// because Realloc uses negative value to descrease the size of volume.
func (vb VolumeBinding) Validate() error {
	if vb.Destination == "" {
		return errors.Wrapf(ErrInvalidVolume, "dest must be provided: %+v", vb)
	}
	if !filepath.IsAbs(vb.Destination) {
		return errors.Wrapf(ErrInvalidVolume, "dest must be absolute: %+v", vb)
	}
	if vb.Source == "" {
		return errors.Wrapf(ErrInvalidVolume, "source must be provided: %+v", vb)
	}
	if !filepath.IsAbs(vb.Source) {
		return errors.Wrapf(ErrInvalidVolume, "source must be absolute: %+v", vb)
	}
	return nil
}

// ToString returns volume string
func (vb VolumeBinding) ToString() (volume string) {
	volume = fmt.Sprintf("%s:%s:%d", vb.Source, vb.Destination, vb.SizeInBytes)
	return volume
}

type VolumeBindings []*VolumeBinding

func (vbs VolumeBindings) Equal(vbs1 VolumeBindings) bool {
	if len(vbs) != len(vbs1) {
		return false
	}
	seen := map[[2]string]*VolumeBinding{}
	for _, vb := range vbs {
		seen[vb.GetMapKey()] = vb
	}
	for _, vb1 := range vbs1 {
		vb, ok := seen[vb1.GetMapKey()]
		if !ok {
			return false
		}
		if *vb != *vb1 {
			return false
		}
	}
	return true
}

func (vbs VolumeBindings) TotalSize() int64 {
	ans := int64(0)
	for _, vb := range vbs {
		ans += vb.SizeInBytes
	}
	return ans
}

func (vbs *VolumeBindings) UnmarshalJSON(b []byte) (err error) {
	volumes := []string{}
	if err = json.Unmarshal(b, &volumes); err != nil {
		return err
	}
	*vbs, err = NewVolumeBindings(volumes)
	return
}

// MarshalJSON is used for encoding/json.Marshal
func (vbs VolumeBindings) MarshalJSON() ([]byte, error) {
	volumes := []string{}
	for _, vb := range vbs {
		volumes = append(volumes, vb.ToString())
	}
	bs, err := json.Marshal(volumes)
	return bs, err
}

// NewVolumeBindings return VolumeBindings of reference type
func NewVolumeBindings(volumes []string) (volumeBindings VolumeBindings, err error) {
	for _, vb := range volumes {
		volumeBinding, err := NewVolumeBinding(vb)
		if err != nil {
			return nil, err
		}
		volumeBindings = append(volumeBindings, volumeBinding)
	}
	return
}

// Validate .
func (vbs VolumeBindings) Validate() error {
	seenDest := map[string]bool{}
	seenSrc := map[string]bool{}
	for _, vb := range vbs {
		if err := vb.Validate(); err != nil {
			return errors.Wrapf(ErrInvalidVolumes, "invalid VolumeBinding: %s", err)
		}
		v := seenDest[vb.Destination]
		if v {
			return errors.Wrapf(ErrInvalidVolumes, "duplicated destination: %s", vb.Destination)
		}
		seenDest[vb.Destination] = true

		src := vb.GetSource()
		if v := seenSrc[src]; v {
			return errors.Wrapf(ErrInvalidVolumes, "duplicated source: %s", src)
		}
		seenSrc[src] = true
	}
	return nil
}

// MergeVolumeBindings combines two VolumeBindings
func MergeVolumeBindings(vbs1 VolumeBindings, vbs2 ...VolumeBindings) (ans VolumeBindings) {
	vbMap := map[[2]string]*VolumeBinding{}
	for _, vbs := range append(vbs2, vbs1) {
		for _, vb := range vbs {
			if binding, ok := vbMap[vb.GetMapKey()]; ok {
				binding.SizeInBytes += vb.SizeInBytes
			} else {
				vbMap[vb.GetMapKey()] = &VolumeBinding{
					Source:      vb.Source,
					Destination: vb.Destination,
					SizeInBytes: vb.SizeInBytes,
				}
			}
		}
	}

	for _, vb := range vbMap {
		if vb.SizeInBytes > 0 {
			ans = append(ans, vb)
		}
	}
	return ans
}

func RemoveEmptyVolumeBinding(vbs VolumeBindings) VolumeBindings {
	var ans VolumeBindings
	for _, vb := range vbs {
		if vb.SizeInBytes > 0 {
			ans = append(ans, vb)
		}
	}
	return ans
}
