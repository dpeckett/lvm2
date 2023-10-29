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

	"github.com/dpeckett/args"
)

type Client struct {
	lvmPath string
}

// Construct a new lvm2 client.
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		lvmPath: "/sbin/lvm",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Display attributes of a physical volume/s.
func (c *Client) ListPhysicalVolumes(ctx context.Context, opts *ListPVOptions) ([]PhysicalVolume, error) {
	cmdArgs := []string{"pvs", "--reportformat=json", "--binary", "--options=pv_all,vg_name"}
	if opts != nil {
		cmdArgs = append(cmdArgs, args.Marshal(opts)...)
	}

	reportJSON, err := c.run(ctx, cmdArgs...)
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
	cmdArgs := []string{"pvcreate", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Change physical volume attributes.
func (c *Client) UpdatePhysicalVolume(ctx context.Context, opts UpdatePVOptions) error {
	cmdArgs := []string{"pvchange", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Remove a physical volume from a device.
func (c *Client) RemovePhysicalVolume(ctx context.Context, opts RemovePVOptions) error {
	cmdArgs := []string{"pvremove", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Check / repair physical volume metadata.
func (c *Client) CheckPhysicalVolume(ctx context.Context, opts CheckPVOptions) error {
	cmdArgs := []string{"pvck", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Move extents from one physical volume to another.
func (c *Client) MovePhysicalExtents(ctx context.Context, opts MovePEOptions) error {
	cmdArgs := []string{"pvmove", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Resize a physical volume.
func (c *Client) ResizePhysicalVolume(ctx context.Context, opts ResizePVOptions) error {
	cmdArgs := []string{"pvresize", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Display volume group/s information.
func (c *Client) ListVolumeGroups(ctx context.Context, opts *ListVGOptions) ([]VolumeGroup, error) {
	cmdArgs := []string{"vgs", "--reportformat=json", "--binary", "--options=vg_all"}
	if opts != nil {
		cmdArgs = append(cmdArgs, args.Marshal(opts)...)
	}

	reportJSON, err := c.run(ctx, cmdArgs...)
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
	cmdArgs := []string{"vgcreate", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Change volume group attributes.
func (c *Client) UpdateVolumeGroup(ctx context.Context, opts UpdateVGOptions) error {
	cmdArgs := []string{"vgchange", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Remove a volume group.
func (c *Client) RemoveVolumeGroup(ctx context.Context, opts RemoveVGOptions) error {
	cmdArgs := []string{"vgremove", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Check / repair volume group metadata.
func (c *Client) CheckVolumeGroup(ctx context.Context, opts CheckVGOptions) error {
	cmdArgs := []string{"vgck", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Unregister a volume group from the system.
func (c *Client) ExportVolumeGroup(ctx context.Context, opts ExportVGOptions) error {
	cmdArgs := []string{"vgexport", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Register a volume group with the system.
func (c *Client) ImportVolumeGroup(ctx context.Context, opts ImportVGOptions) error {
	cmdArgs := []string{"vgimport", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Import a volume group from cloned physical volumes.
func (c *Client) ImportVolumeGroupFromCloned(ctx context.Context, opts ImportVGFromClonedOptions) error {
	cmdArgs := []string{"vgimportclone", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Merge volume groups.
func (c *Client) MergeVolumeGroups(ctx context.Context, opts MergeVGOptions) error {
	cmdArgs := []string{"vgmerge", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Add physical volumes to a volume group.
func (c *Client) ExtendVolumeGroup(ctx context.Context, opts ExtendVGOptions) error {
	cmdArgs := []string{"vgextend", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Remove physical volumes from a volume group.
func (c *Client) ReduceVolumeGroup(ctx context.Context, opts ReduceVGOptions) error {
	cmdArgs := []string{"vgreduce", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Rename a volume group.
func (c *Client) RenameVolumeGroup(ctx context.Context, opts RenameVGOptions) error {
	cmdArgs := []string{"vgrename", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Move physical volumes between volume groups.
func (c *Client) MovePhysicalVolumes(ctx context.Context, opts MovePVOptions) error {
	cmdArgs := []string{"vgsplit", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Create device files for active logical volumes in the volume group.
func (c *Client) MakeVolumeGroupDeviceNodes(ctx context.Context, opts MakeVGDeviceNodesOptions) error {
	cmdArgs := []string{"vgmknodes", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Display logical volume/s information.
func (c *Client) ListLogicalVolumes(ctx context.Context, opts *ListLVOptions) ([]LogicalVolume, error) {
	cmdArgs := []string{"lvs", "--reportformat=json", "--binary", "--options=lv_all,seg_all,vg_name"}
	if opts != nil {
		cmdArgs = append(cmdArgs, args.Marshal(opts)...)
	}

	reportJSON, err := c.run(ctx, cmdArgs...)
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
	cmdArgs := []string{"lvcreate", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Change logical volume attributes.
func (c *Client) UpdateLogicalVolume(ctx context.Context, opts UpdateLVOptions) error {
	cmdArgs := []string{"lvchange", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Remove a logical volume.
func (c *Client) RemoveLogicalVolume(ctx context.Context, opts RemoveLVOptions) error {
	cmdArgs := []string{"lvremove", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Change logical volume layout.
func (c *Client) ConvertLogicalVolumeLayout(ctx context.Context, opts ConvertLVLayoutOptions) error {
	cmdArgs := []string{"lvconvert", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Add space to a logical volume.
func (c *Client) ExtendLogicalVolume(ctx context.Context, opts ExtendLVOptions) error {
	cmdArgs := []string{"lvextend", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Reduce the size of a logical volume.
func (c *Client) ReduceLogicalVolume(ctx context.Context, opts ReduceLVOptions) error {
	cmdArgs := []string{"lvreduce", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

// Rename a logical volume.
func (c *Client) RenameLogicalVolume(ctx context.Context, opts RenameLVOptions) error {
	cmdArgs := []string{"lvrename", "--yes"}
	cmdArgs = append(cmdArgs, args.Marshal(opts)...)

	_, err := c.run(ctx, cmdArgs...)
	return err
}

func (c *Client) run(ctx context.Context, cmdArgs ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, c.lvmPath, cmdArgs...)

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w: %s", err, errOut.String())
	}

	return out.Bytes(), nil
}
