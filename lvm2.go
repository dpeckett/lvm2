/* SPDX-License-Identifier: Apache-2.0
 *
 * Copyright 2023 Damian Peckett <damian@pecke.tt>.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package lvm2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

// Display attributes of a physical volume/s.
func (c *Client) ListPhysicalVolumes(ctx context.Context, opts *ListPVOptions) ([]PhysicalVolume, error) {
	args := []string{"pvs", "--reportformat=json", "--binary", "--options=pv_all,vg_name"}
	if opts != nil {
		args = append(args, MarshalArgs(opts)...)
	}

	reportJSON, err := c.run(ctx, args...)
	if err != nil {
		return nil, err
	}

	var report struct {
		Report []struct {
			PV []PhysicalVolume `json:"pv"`
		} `json:"report"`
	}
	if err := json.Unmarshal(reportJSON, &report); err != nil {
		return nil, fmt.Errorf("failed to parse lvm output: %w", err)
	}

	if len(report.Report) > 0 && len(report.Report[0].PV) > 0 {
		return report.Report[0].PV, nil
	}

	return nil, nil
}

// Create a new physical volume on a device.
func (c *Client) CreatePhysicalVolume(ctx context.Context, opts CreatePVOptions) error {
	args := []string{"pvcreate", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Change physical volume attributes.
func (c *Client) UpdatePhysicalVolume(ctx context.Context, opts UpdatePVOptions) error {
	args := []string{"pvchange", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Remove a physical volume from a device.
func (c *Client) RemovePhysicalVolume(ctx context.Context, opts RemovePVOptions) error {
	args := []string{"pvremove", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Check / repair physical volume metadata.
func (c *Client) CheckPhysicalVolume(ctx context.Context, opts CheckPVOptions) error {
	args := []string{"pvck", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Move extents from one physical volume to another.
func (c *Client) MovePhysicalExtents(ctx context.Context, opts MovePEOptions) error {
	args := []string{"pvmove", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Resize a physical volume.
func (c *Client) ResizePhysicalVolume(ctx context.Context, opts ResizePVOptions) error {
	args := []string{"pvresize", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Display volume group/s information.
func (c *Client) ListVolumeGroups(ctx context.Context, opts *ListVGOptions) ([]VolumeGroup, error) {
	args := []string{"vgs", "--reportformat=json", "--binary", "--options=vg_all"}
	if opts != nil {
		args = append(args, MarshalArgs(opts)...)
	}

	reportJSON, err := c.run(ctx, args...)
	if err != nil {
		return nil, err
	}

	var report struct {
		Report []struct {
			VG []VolumeGroup `json:"vg"`
		} `json:"report"`
	}
	if err := json.Unmarshal(reportJSON, &report); err != nil {
		return nil, fmt.Errorf("failed to parse lvm output: %w", err)
	}

	if len(report.Report) > 0 && len(report.Report[0].VG) > 0 {
		return report.Report[0].VG, nil
	}

	return nil, nil
}

// Create a new volume group.
func (c *Client) CreateVolumeGroup(ctx context.Context, opts CreateVGOptions) error {
	args := []string{"vgcreate", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Change volume group attributes.
func (c *Client) UpdateVolumeGroup(ctx context.Context, opts UpdateVGOptions) error {
	args := []string{"vgchange", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Remove a volume group.
func (c *Client) RemoveVolumeGroup(ctx context.Context, opts RemoveVGOptions) error {
	args := []string{"vgremove", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Check / repair volume group metadata.
func (c *Client) CheckVolumeGroup(ctx context.Context, opts CheckVGOptions) error {
	args := []string{"vgck", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Unregister a volume group from the system.
func (c *Client) ExportVolumeGroup(ctx context.Context, opts ExportVGOptions) error {
	args := []string{"vgexport", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Register a volume group with the system.
func (c *Client) ImportVolumeGroup(ctx context.Context, opts ImportVGOptions) error {
	args := []string{"vgimport", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Import a volume group from cloned physical volumes.
func (c *Client) ImportVolumeGroupFromCloned(ctx context.Context, opts ImportVGFromClonedOptions) error {
	args := []string{"vgimportclone", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Merge volume groups.
func (c *Client) MergeVolumeGroups(ctx context.Context, opts MergeVGOptions) error {
	args := []string{"vgmerge", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Add physical volumes to a volume group.
func (c *Client) ExtendVolumeGroup(ctx context.Context, opts ExtendVGOptions) error {
	args := []string{"vgextend", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Remove physical volumes from a volume group.
func (c *Client) ReduceVolumeGroup(ctx context.Context, opts ReduceVGOptions) error {
	args := []string{"vgreduce", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Rename a volume group.
func (c *Client) RenameVolumeGroup(ctx context.Context, opts RenameVGOptions) error {
	args := []string{"vgrename", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Move physical volumes between volume groups.
func (c *Client) MovePhysicalVolumes(ctx context.Context, opts MovePVOptions) error {
	args := []string{"vgsplit", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Create device files for active logical volumes in the volume group.
func (c *Client) MakeVolumeGroupDeviceNodes(ctx context.Context, opts MakeVGDeviceNodesOptions) error {
	args := []string{"vgmknodes", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Display logical volume/s information.
func (c *Client) ListLogicalVolumes(ctx context.Context, opts *ListLVOptions) ([]LogicalVolume, error) {
	args := []string{"lvs", "--reportformat=json", "--binary", "--options=lv_all,seg_all,vg_name"}
	if opts != nil {
		args = append(args, MarshalArgs(opts)...)
	}

	reportJSON, err := c.run(ctx, args...)
	if err != nil {
		return nil, err
	}

	var report struct {
		Report []struct {
			LV []LogicalVolume `json:"lv"`
		} `json:"report"`
	}
	if err := json.Unmarshal(reportJSON, &report); err != nil {
		return nil, fmt.Errorf("failed to parse lvm output: %w", err)
	}

	if len(report.Report) > 0 && len(report.Report[0].LV) > 0 {
		return report.Report[0].LV, nil
	}

	return nil, nil
}

// Create a new logical volume in a volume group.
func (c *Client) CreateLogicalVolume(ctx context.Context, opts CreateLVOptions) error {
	args := []string{"lvcreate", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Change logical volume attributes.
func (c *Client) UpdateLogicalVolume(ctx context.Context, opts UpdateLVOptions) error {
	args := []string{"lvchange", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Remove a logical volume.
func (c *Client) RemoveLogicalVolume(ctx context.Context, opts RemoveLVOptions) error {
	args := []string{"lvremove", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Change logical volume layout.
func (c *Client) ConvertLogicalVolumeLayout(ctx context.Context, opts ConvertLVLayoutOptions) error {
	args := []string{"lvconvert", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Add space to a logical volume.
func (c *Client) ExtendLogicalVolume(ctx context.Context, opts ExtendLVOptions) error {
	args := []string{"lvextend", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Reduce the size of a logical volume.
func (c *Client) ReduceLogicalVolume(ctx context.Context, opts ReduceLVOptions) error {
	args := []string{"lvreduce", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Rename a logical volume.
func (c *Client) RenameLogicalVolume(ctx context.Context, opts RenameLVOptions) error {
	args := []string{"lvrename", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

func (c *Client) run(ctx context.Context, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "/sbin/lvm", args...)

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w: %s", err, errOut.String())
	}

	return out.Bytes(), nil
}
