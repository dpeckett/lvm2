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

func New() *Client {
	return &Client{}
}

// Display attributes of a physical volume/s.
func (c *Client) PVList(ctx context.Context, opts *PVListOptions) ([]PhysicalVolume, error) {
	args := []string{"pvs", "--reportformat=json", "--options=pv_all,vg_name"}
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
func (c *Client) PVCreate(ctx context.Context, opts PVCreateOptions) error {
	args := []string{"pvcreate", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Change physical volume attributes.
func (c *Client) PVChange(ctx context.Context, opts PVChangeOptions) error {
	args := []string{"pvchange", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Remove a physical volume from a device.
func (c *Client) PVRemove(ctx context.Context, opts PVRemoveOptions) error {
	args := []string{"pvremove", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Check / repair physical volume metadata.
func (c *Client) PVCheck(ctx context.Context, opts PVCheckOptions) error {
	args := []string{"pvck", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Move extents from one physical volume to another.
func (c *Client) PVMove(ctx context.Context, opts PVMoveOptions) error {
	args := []string{"pvmove", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Resize a physical volume.
func (c *Client) PVResize(ctx context.Context, opts PVResizeOptions) error {
	args := []string{"pvresize", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Display volume group/s information.
func (c *Client) VGList(ctx context.Context, opts *VGListOptions) ([]VolumeGroup, error) {
	args := []string{"vgs", "--reportformat=json", "--options=vg_all"}
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
func (c *Client) VGCreate(ctx context.Context, opts VGCreateOptions) error {
	args := []string{"vgcreate", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Change volume group attributes.
func (c *Client) VGChange(ctx context.Context, opts VGChangeOptions) error {
	args := []string{"vgchange", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Remove a volume group.
func (c *Client) VGRemove(ctx context.Context, opts VGRemoveOptions) error {
	args := []string{"vgremove", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Check / repair volume group metadata.
func (c *Client) VGCheck(ctx context.Context, opts VGCheckOptions) error {
	args := []string{"vgck", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Unregister a volume group from the system.
func (c *Client) VGExport(ctx context.Context, opts VGExportOptions) error {
	args := []string{"vgexport", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Register a volume group with the system.
func (c *Client) VGImport(ctx context.Context, opts VGImportOptions) error {
	args := []string{"vgimport", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Import a volume group from cloned physical volumes.
func (c *Client) VGImportClone(ctx context.Context, opts VGImportCloneOptions) error {
	args := []string{"vgimportclone", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Add physical volumes to a volume group.
func (c *Client) VGExtend(ctx context.Context, opts VGExtendOptions) error {
	args := []string{"vgextend", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Merge volume groups.
func (c *Client) VGMerge(ctx context.Context, opts VGMergeOptions) error {
	args := []string{"vgmerge", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Remove physical volumes from a volume group.
func (c *Client) VGReduce(ctx context.Context, opts VGReduceOptions) error {
	args := []string{"vgreduce", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Rename a volume group.
func (c *Client) VGRename(ctx context.Context, opts VGRenameOptions) error {
	args := []string{"vgrename", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Move physical volumes between volume groups.
func (c *Client) VGSplit(ctx context.Context, opts VGSplitOptions) error {
	args := []string{"vgsplit", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Create device files for active logical volumes in the volume group.
func (c *Client) VGMknodes(ctx context.Context, opts VGMknodesOptions) error {
	args := []string{"vgmknodes", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Create a new logical volume in a volume group.
func (c *Client) LVCreate(ctx context.Context, opts LVCreateOptions) error {
	args := []string{"lvcreate", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Display logical volume/s information.
func (c *Client) LVList(ctx context.Context, opts *LVListOptions) ([]LogicalVolume, error) {
	args := []string{"lvs", "--reportformat=json", "--options=lv_all,seg_all,vg_name"}
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

// Change logical volume attributes.
func (c *Client) LVChange(ctx context.Context, opts LVChangeOptions) error {
	args := []string{"lvchange", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Remove a logical volume.
func (c *Client) LVRemove(ctx context.Context, opts LVRemoveOptions) error {
	args := []string{"lvremove", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Change logical volume layout.
func (c *Client) LVConvert(ctx context.Context, opts LVConvertOptions) error {
	args := []string{"lvconvert", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Add space to a logical volume.
func (c *Client) LVExtend(ctx context.Context, opts LVExtendOptions) error {
	args := []string{"lvextend", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Reduce the size of a logical volume.
func (c *Client) LVReduce(ctx context.Context, opts LVReduceOptions) error {
	args := []string{"lvreduce", "--yes"}
	args = append(args, MarshalArgs(opts)...)

	_, err := c.run(ctx, args...)
	return err
}

// Rename a logical volume.
func (c *Client) LVRename(ctx context.Context, opts LVRenameOptions) error {
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
