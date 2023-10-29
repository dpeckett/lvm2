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

package lvm2_test

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dpeckett/lvm2"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	err := loadNBDModule()
	require.NoError(t, err)

	c := lvm2.NewClient()

	t.Run("Physical volumes", func(t *testing.T) {
		t.Log("Creating virtual block device")

		imagePath := filepath.Join(t.TempDir(), ".qcow2")
		err = createImage(imagePath)
		require.NoError(t, err)

		devPath, err := attachNBDDevice(imagePath)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := detachNBDDevice(devPath)
			require.NoError(t, err)
		})

		t.Log("Virtual block device created", devPath)

		ctx := context.Background()

		t.Log("Creating physical volume")

		err = c.CreatePhysicalVolume(ctx, lvm2.CreatePVOptions{
			Name: devPath,
		})
		require.NoError(t, err, "failed to create PV")

		pvs, err := c.ListPhysicalVolumes(ctx, &lvm2.ListPVOptions{
			Names: []string{devPath},
		})
		require.NoError(t, err, "failed to list PVs")

		require.Len(t, pvs, 1)
		require.Equal(t, devPath, pvs[0].Name)

		t.Log("Changing physical volume UUID")

		err = c.UpdatePhysicalVolume(ctx, lvm2.UpdatePVOptions{
			Name: devPath,
			UUID: true,
		})
		require.NoError(t, err, "failed to change PV")

		t.Log("Resizing physical volume")

		err = c.ResizePhysicalVolume(ctx, lvm2.ResizePVOptions{
			Name:                  devPath,
			SetPhysicalVolumeSize: "100M",
		})
		require.NoError(t, err, "failed to resize PV")

		pvs, err = c.ListPhysicalVolumes(ctx, &lvm2.ListPVOptions{
			Names: []string{devPath},
		})
		require.NoError(t, err, "failed to list PVs")

		require.Len(t, pvs, 1)
		require.Equal(t, devPath, pvs[0].Name)
		require.Equal(t, "100.00m", pvs[0].Size)

		t.Log("Checking physical volume")

		err = c.CheckPhysicalVolume(ctx, lvm2.CheckPVOptions{
			Name: devPath,
		})
		require.NoError(t, err, "failed to check PV")

		t.Log("Removing physical volume")

		err = c.RemovePhysicalVolume(ctx, lvm2.RemovePVOptions{
			Name: devPath,
		})
		require.NoError(t, err, "failed to remove PV")

		pvs, err = c.ListPhysicalVolumes(ctx, &lvm2.ListPVOptions{
			Names: []string{devPath},
		})
		require.Contains(t, err.Error(), "Failed to find physical volume")
		require.Empty(t, pvs)
	})

	t.Run("Volume groups", func(t *testing.T) {
		t.Log("Creating virtual block devices")

		firstImagePath := filepath.Join(t.TempDir(), ".qcow2")
		err = createImage(firstImagePath)
		require.NoError(t, err)

		secondImagePath := filepath.Join(t.TempDir(), ".qcow2")
		err = createImage(secondImagePath)
		require.NoError(t, err)

		firstDevPath, err := attachNBDDevice(firstImagePath)
		require.NoError(t, err)

		secondDevPath, err := attachNBDDevice(secondImagePath)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := detachNBDDevice(firstDevPath)
			require.NoError(t, err)

			err = detachNBDDevice(secondDevPath)
			require.NoError(t, err)
		})

		t.Log("Virtual block devices created", firstDevPath, secondDevPath)

		ctx := context.Background()

		t.Log("Creating physical volumes")

		err = c.CreatePhysicalVolume(ctx, lvm2.CreatePVOptions{
			Name: firstDevPath,
		})
		require.NoError(t, err, "failed to create first PV")

		err = c.CreatePhysicalVolume(ctx, lvm2.CreatePVOptions{
			Name: secondDevPath,
		})
		require.NoError(t, err, "failed to create second PV")

		t.Log("Creating volume group")

		vgName := strings.ReplaceAll(t.Name()+"_vg", "/", "_")
		vgTag := strings.ReplaceAll(t.Name(), "/", "_")

		err = c.CreateVolumeGroup(ctx, lvm2.CreateVGOptions{
			Name:    vgName,
			PVNames: []string{firstDevPath},
			Tags:    []string{vgTag},
		})
		require.NoError(t, err, "failed to create VG")

		t.Cleanup(func() {
			_ = c.UpdateVolumeGroup(ctx, lvm2.UpdateVGOptions{
				Name:     vgName,
				Activate: lvm2.No,
			})

			_ = c.RemoveVolumeGroup(ctx, lvm2.RemoveVGOptions{
				Name: vgName,
			})
		})

		vgs, err := c.ListVolumeGroups(ctx, &lvm2.ListVGOptions{
			Names: []string{vgName},
		})
		require.NoError(t, err, "failed to list VGs")

		require.Len(t, vgs, 1)
		require.Equal(t, vgName, vgs[0].Name)
		require.Equal(t, 1, int(vgs[0].PVCount))

		t.Log("Activating volume group")

		err = c.UpdateVolumeGroup(ctx, lvm2.UpdateVGOptions{
			Name:     vgName,
			Activate: lvm2.Yes,
		})
		require.NoError(t, err, "failed to activate VG")

		t.Log("Adding second physical volume to volume group")

		err = c.ExtendVolumeGroup(ctx, lvm2.ExtendVGOptions{
			Name:    vgName,
			PVNames: []string{secondDevPath},
		})
		require.NoError(t, err, "failed to add second PV to VG")

		vgs, err = c.ListVolumeGroups(ctx, &lvm2.ListVGOptions{
			Names: []string{vgName},
		})
		require.NoError(t, err, "failed to list VGs")

		require.Len(t, vgs, 1)
		require.Equal(t, 2, int(vgs[0].PVCount))

		t.Log("Splitting volume group")

		tmpSecondVGName := strings.ReplaceAll(t.Name()+"_tmp", "/", "_")

		err = c.MovePhysicalVolumes(ctx, lvm2.MovePVOptions{
			Source:      vgName,
			Destination: tmpSecondVGName,
			PVNames:     []string{secondDevPath},
		})
		require.NoError(t, err, "failed to split VG")

		t.Log("Renaming volume group")

		secondVGName := strings.ReplaceAll(t.Name()+"_vg2", "/", "_")

		err = c.RenameVolumeGroup(ctx, lvm2.RenameVGOptions{
			From: tmpSecondVGName,
			To:   secondVGName,
		})
		require.NoError(t, err, "failed to rename VG")

		t.Cleanup(func() {
			_ = c.RemoveVolumeGroup(ctx, lvm2.RemoveVGOptions{
				Name: secondVGName,
			})
		})

		t.Log("Adding tag to volume group")

		err = c.UpdateVolumeGroup(ctx, lvm2.UpdateVGOptions{
			Name:    secondVGName,
			AddTags: []string{vgTag},
		})
		require.NoError(t, err, "failed to add tag to VG")

		vgs, err = c.ListVolumeGroups(ctx, &lvm2.ListVGOptions{
			Select: "vg_tags=" + vgTag,
		})
		require.NoError(t, err, "failed to list VGs")

		require.Len(t, vgs, 2)
		require.Equal(t, 1, int(vgs[0].PVCount))
		require.Equal(t, 1, int(vgs[1].PVCount))

		t.Log("Merging volume groups")

		err = c.MergeVolumeGroups(ctx, lvm2.MergeVGOptions{
			Destination: vgName,
			Source:      secondVGName,
		})
		require.NoError(t, err, "failed to merge VGs")

		vgs, err = c.ListVolumeGroups(ctx, &lvm2.ListVGOptions{
			Select: "vg_tags=" + vgTag,
		})
		require.NoError(t, err, "failed to list VGs")

		require.Len(t, vgs, 1)

		t.Log("Removing second physical volume from volume group")

		err = c.ReduceVolumeGroup(ctx, lvm2.ReduceVGOptions{
			Name: vgName,
			PVNames: []string{
				secondDevPath,
			},
		})
		require.NoError(t, err, "failed to remove VG")

		vgs, err = c.ListVolumeGroups(ctx, &lvm2.ListVGOptions{
			Select: "vg_tags=" + vgTag,
		})
		require.NoError(t, err, "failed to list VGs")

		require.Len(t, vgs, 1)

		t.Log("Checking volume group")

		err = c.CheckVolumeGroup(ctx, lvm2.CheckVGOptions{
			Name:           vgName,
			UpdateMetadata: true,
		})
		require.NoError(t, err, "failed to check VG")

		t.Log("Removing volume group")

		err = c.RemoveVolumeGroup(ctx, lvm2.RemoveVGOptions{
			Name: vgName,
		})
		require.NoError(t, err, "failed to remove VG")

		vgs, err = c.ListVolumeGroups(ctx, &lvm2.ListVGOptions{
			Names: []string{vgName},
		})
		require.Contains(t, err.Error(), "not found")
		require.Empty(t, vgs)
	})

	t.Run("Logical volumes", func(t *testing.T) {
		t.Log("Creating virtual block device")

		imagePath := filepath.Join(t.TempDir(), ".qcow2")
		err = createImage(imagePath)
		require.NoError(t, err)

		devPath, err := attachNBDDevice(imagePath)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := detachNBDDevice(devPath)
			require.NoError(t, err)
		})

		t.Log("Virtual block device created", devPath)

		ctx := context.Background()

		vgName := strings.ReplaceAll(t.Name()+"_vg", "/", "_")

		t.Log("Creating volume group", vgName)

		err = c.CreateVolumeGroup(ctx, lvm2.CreateVGOptions{
			Name:    vgName,
			PVNames: []string{devPath},
		})
		require.NoError(t, err, "failed to create VG")

		t.Cleanup(func() {
			err := c.UpdateVolumeGroup(ctx, lvm2.UpdateVGOptions{
				Name:     vgName,
				Activate: lvm2.No,
			})
			require.NoError(t, err)
		})

		lvName := strings.ReplaceAll(t.Name()+"_lv", "/", "_")

		t.Log("Creating logical volume", lvName)

		err = c.CreateLogicalVolume(ctx, lvm2.CreateLVOptions{
			Name:     lvName,
			VGName:   vgName,
			Size:     "100M",
			Activate: lvm2.No,
		})
		require.NoError(t, err, "failed to create LV")

		lvs, err := c.ListLogicalVolumes(ctx, &lvm2.ListLVOptions{
			Names: []string{
				fmt.Sprintf("%s/%s", vgName, lvName),
			},
		})
		require.NoError(t, err, "failed to list LVs")

		require.Len(t, lvs, 1)
		require.Equal(t, lvName, lvs[0].Name)
		require.Equal(t, "100.00m", lvs[0].Size)
		require.Empty(t, lvs[0].Active)

		t.Log("Resizing logical volume")

		err = c.ExtendLogicalVolume(ctx, lvm2.ExtendLVOptions{
			Name: fmt.Sprintf("%s/%s", vgName, lvName),
			Size: "200M",
		})
		require.NoError(t, err, "failed to extend LV")

		lvs, err = c.ListLogicalVolumes(ctx, &lvm2.ListLVOptions{
			Names: []string{
				fmt.Sprintf("%s/%s", vgName, lvName),
			},
		})
		require.NoError(t, err, "failed to list LVs")

		require.Len(t, lvs, 1)
		require.Equal(t, "200.00m", lvs[0].Size)

		err = c.ReduceLogicalVolume(ctx, lvm2.ReduceLVOptions{
			Name: fmt.Sprintf("%s/%s", vgName, lvName),
			Size: "128M",
		})
		require.NoError(t, err, "failed to reduce LV")

		lvs, err = c.ListLogicalVolumes(ctx, &lvm2.ListLVOptions{
			Names: []string{
				fmt.Sprintf("%s/%s", vgName, lvName),
			},
		})
		require.NoError(t, err, "failed to list LVs")

		require.Len(t, lvs, 1)
		require.Equal(t, "128.00m", lvs[0].Size)

		t.Log("Renaming and activating logical volume")

		err = c.RenameLogicalVolume(ctx, lvm2.RenameLVOptions{
			From: fmt.Sprintf("%s/%s", vgName, lvName),
			To:   lvName + "_2",
		})
		require.NoError(t, err, "failed to rename LV")

		lvName += "_2"

		err = c.UpdateLogicalVolume(ctx, lvm2.UpdateLVOptions{
			Name:     fmt.Sprintf("%s/%s", vgName, lvName),
			Activate: lvm2.Yes,
		})
		require.NoError(t, err, "failed to activate LV")

		lvs, err = c.ListLogicalVolumes(ctx, &lvm2.ListLVOptions{
			Names: []string{
				fmt.Sprintf("%s/%s", vgName, lvName),
			},
		})
		require.NoError(t, err, "failed to list LVs")

		require.Len(t, lvs, 1)
		require.Equal(t, lvName, lvs[0].Name)
		require.NotEmpty(t, lvs[0].Active)

		t.Log("Creating second virtual block device")

		secondImagePath := filepath.Join(t.TempDir(), ".qcow2")
		err = createImage(secondImagePath)
		require.NoError(t, err)

		secondDevPath, err := attachNBDDevice(secondImagePath)
		require.NoError(t, err)

		t.Cleanup(func() {
			err := detachNBDDevice(secondDevPath)
			require.NoError(t, err)
		})

		t.Log("Virtual block device created", secondDevPath)

		t.Log("Adding second physical volume to volume group")

		err = c.CreatePhysicalVolume(ctx, lvm2.CreatePVOptions{
			Name: secondDevPath,
		})
		require.NoError(t, err, "failed to create second PV")

		err = c.ExtendVolumeGroup(ctx, lvm2.ExtendVGOptions{
			Name:    vgName,
			PVNames: []string{secondDevPath},
		})
		require.NoError(t, err, "failed to extend VG")

		t.Log("Converting logical volume to RAID1")

		err = c.ConvertLogicalVolumeLayout(ctx, lvm2.ConvertLVLayoutOptions{
			Name:    fmt.Sprintf("%s/%s", vgName, lvName),
			Type:    "raid1",
			Mirrors: lvm2.PtrTo(1),
		})
		require.NoError(t, err, "failed to convert LV to RAID1")

		lvs, err = c.ListLogicalVolumes(ctx, &lvm2.ListLVOptions{
			Names: []string{
				fmt.Sprintf("%s/%s", vgName, lvName),
			},
		})
		require.NoError(t, err, "failed to list LVs")

		require.Len(t, lvs, 1)
		require.Equal(t, "raid1", lvs[0].Type, "expected LV to be of type RAID1")

		t.Log("Removing second physical volume from volume group")

		err = c.RemoveLogicalVolume(ctx, lvm2.RemoveLVOptions{
			Name: fmt.Sprintf("%s/%s", vgName, lvName),
		})
		require.NoError(t, err, "failed to remove logical volume")
	})
}

func loadNBDModule() error {
	cmd := exec.Command("/sbin/modprobe", "nbd", "max_part=16")
	return cmd.Run()
}

func createImage(imagePath string) error {
	cmd := exec.Command("qemu-img", "create", "-f", "qcow2", imagePath, "1G")
	return cmd.Run()
}

func attachNBDDevice(imagePath string) (string, error) {
	for i := 0; i < 16; i++ {
		devPath := fmt.Sprintf("/dev/nbd%d", i)
		cmd := exec.Command("qemu-nbd", "-c", devPath, imagePath)
		err := cmd.Run()
		if err == nil {
			return devPath, nil
		}
	}

	return "", fmt.Errorf("no free nbd device found")
}

func detachNBDDevice(devPath string) error {
	cmd := exec.Command("qemu-nbd", "-d", devPath)
	return cmd.Run()
}
