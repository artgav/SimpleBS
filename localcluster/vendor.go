// localcluster/vendor.go - Core Business Logic
package localcluster

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Vendor struct {
	Folder string
}

func NewVendor(folder string) *Vendor {
	os.MkdirAll(folder, 0755)
	return &Vendor{Folder: folder}
}

func (v *Vendor) volumePath(name string) string {
	return filepath.Join(v.Folder, name)
}

func (v *Vendor) metaPath(name string) string {
	return v.volumePath(name) + ".meta.json"
}

func (v *Vendor) CreateVolume(name string, size int64) error {
	path := v.volumePath(name)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("volume %s already exists", name)
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := f.Truncate(size); err != nil {
		return err
	}
	return v.writeMeta(name, map[string]interface{}{
		"status": "available",
		"size":   size,
	})
}

func (v *Vendor) ResizeVolume(name string, newSize int64) error {
	path := v.volumePath(name)
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := f.Truncate(newSize); err != nil {
		return err
	}
	meta, _ := v.readMeta(name)
	meta["size"] = newSize
	return v.writeMeta(name, meta)
}

func (v *Vendor) DeleteVolume(name string) error {
	if err := os.Remove(v.volumePath(name)); err != nil {
		return err
	}
	_ = os.Remove(v.metaPath(name))
	return nil
}

func (v *Vendor) GetVolumeInfo(name string) (map[string]interface{}, error) {
	path := v.volumePath(name)
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	meta, _ := v.readMeta(name)
	meta["size"] = info.Size()
	return meta, nil
}

func (v *Vendor) CloneVolume(src, dst string) error {
	srcPath := v.volumePath(src)
	dstPath := v.volumePath(dst)
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	meta, _ := v.readMeta(src)
	return v.writeMeta(dst, meta)
}

func (v *Vendor) AttachVolume(name, serverIP string) error {
	meta, _ := v.readMeta(name)
	meta["status"] = "attached"
	meta["server_ip"] = serverIP
	return v.writeMeta(name, meta)
}

func (v *Vendor) DetachVolume(name string) error {
	meta, _ := v.readMeta(name)
	meta["status"] = "available"
	delete(meta, "server_ip")
	return v.writeMeta(name, meta)
}

func (v *Vendor) readMeta(name string) (map[string]interface{}, error) {
	f, err := os.Open(v.metaPath(name))
	if err != nil {
		return map[string]interface{}{}, nil
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	meta := map[string]interface{}{}
	err = decoder.Decode(&meta)
	return meta, err
}

func (v *Vendor) writeMeta(name string, meta map[string]interface{}) error {
	f, err := os.Create(v.metaPath(name))
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(meta)
}
