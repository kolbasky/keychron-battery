// Code generated for package main by go-bindata DO NOT EDIT. (@generated)
// sources:
// icon.ico
package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _iconIco = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x97\x5d\x48\x54\x5b\x14\xc7\x97\xf8\x70\x5f\x84\x2b\xf7\xc2\x7d\xbd\xf3\x58\xd9\xbb\x8f\x46\x71\x26\x8d\x82\xd0\x74\x0f\x81\x14\xbd\x55\x3a\x96\x29\x0e\x44\x39\x13\x39\xb3\x67\x72\x8a\x9a\x4c\xc2\xd4\x74\x3e\xf4\x1c\x23\xc9\x26\x2b\xc1\x9a\xb2\x87\x22\xb2\x9c\x87\xb0\x34\xd0\x08\x44\x22\x51\x88\x18\x72\x70\xc7\x3a\x79\xe4\xcc\x9c\x39\xe7\xec\x33\x89\x0b\xfe\x08\x32\x7b\xfd\xf6\x5e\x7b\xed\xb5\xd6\x01\x28\x80\x02\xb0\xd9\xf0\xaf\x0d\x06\x8b\x01\xfe\x03\x80\x6d\x00\x60\x03\x80\x5d\xf0\xfb\xff\x68\xee\x62\xc8\xdb\x82\x52\xf5\x3f\x54\x74\x1c\xa1\x22\xe9\xa3\x22\x99\xf2\x4b\x8e\x65\xbf\xe4\x60\x28\x4f\x5f\xe5\xb2\xd3\x6f\x4f\x3a\xa9\xd0\x5f\xe7\x13\x8e\x9e\xf4\xee\xf9\x37\x7f\x52\xa6\xd1\x21\x47\x89\x5f\x22\x61\x2a\x92\x94\xc2\x53\x44\x45\xc2\xce\x75\x1d\x64\x4e\x6a\xcf\x94\x4f\x48\x39\xa9\x10\x71\xb6\x09\x3b\xf3\xe5\xb6\x4a\xd5\x45\x7e\xc9\x11\xa2\x12\x49\x67\x73\x51\xbe\x41\xc2\x5a\x42\xfb\xb5\x6c\x95\xea\x7d\x42\xba\xde\x67\xef\x38\xd1\x5a\x56\x64\xe9\xcc\x03\x64\x3b\x95\xc8\x74\x2e\x2e\xaa\x2d\x5a\xcd\x1a\x83\x15\x86\xec\x4c\x09\x1f\xea\xdb\x76\xef\xe0\x61\x7b\x63\x8e\x52\x2a\x92\x25\x3d\xb6\xbb\xb7\x92\x35\x04\xf6\x5a\x60\x6f\xec\x61\xa9\xc1\x2b\x94\x9a\x9e\x5b\x87\x8d\x77\x1d\x88\x1e\x63\xee\xce\x5a\xe6\xbe\xa1\xaf\xb3\x21\x62\xb8\x07\xbd\x38\xe0\x7d\x1b\xc5\xbc\xe7\x91\x8b\xf1\xd8\xda\xda\x1a\xf3\xf7\x1c\x37\xbc\x8b\x5c\xf9\x80\xb9\xa6\xc7\x46\x45\x9f\x78\x64\xff\xf1\x67\xbd\xb9\xcf\xde\x59\xcb\x5e\x4c\xc6\xe5\xdf\x5c\x8b\x36\x99\xe4\xa5\xbd\x23\x23\xee\x43\x8e\x12\xbd\x3c\xcf\xe6\x1b\xf9\x8e\xc4\x2f\x71\xf2\x85\xb4\xfa\x6d\xe2\xfb\x36\x62\xa3\x62\x4f\x37\x8f\x2f\xef\x81\x0a\x51\xd8\xa8\x6b\xda\xda\xa2\xc7\xbf\x7d\xcf\xab\xeb\x73\x24\xd1\xcd\xcd\xc7\x1a\x85\x75\x12\x6b\xaa\x19\x5b\xcd\xff\x91\xfa\xce\x3e\xce\xbf\xd3\x68\xe6\x73\x92\xad\xae\xfe\xe4\xe7\x53\x3b\xc3\x5a\x8d\xf5\xdc\x0a\xff\xd6\x5d\x8f\xae\xbf\xe1\xf1\x9b\x96\xf8\xd8\x2f\xd6\x7b\x09\x37\x7f\xb3\xee\x7f\x3d\x07\x92\xea\x3e\xb6\xd5\xfc\x53\x01\xfb\x0a\x0f\x5b\xcd\x4f\xbc\x1e\x96\xfd\x6b\x14\x6b\x62\x6f\xde\x27\xb8\xf9\x8d\xc1\x72\xb9\x8f\x58\xe5\x2f\x7c\x9d\x63\x2f\x93\x8f\x35\x7a\x95\x1c\x63\xdf\x56\x16\xb9\xf8\xcd\x57\xf7\x31\xdf\x40\x8d\xe2\xdb\x52\xfc\xaf\x0f\xb4\xe8\xfa\x8d\x8d\x5e\x36\xe5\xbb\x3a\x0e\xc8\xbd\x44\x61\x6f\x65\xfe\xe1\xac\x92\xd5\xd3\xa6\xac\xbe\xbf\x50\xac\x59\x97\x1f\x7d\xd0\x9e\x93\xdf\xe0\xb7\x33\x4f\x6f\x65\xae\x9e\xda\x67\xb5\xfe\xcc\xcc\x4f\xb1\xd1\x89\x7e\x8d\x1e\x4e\x84\xd9\x97\xc5\x59\x0d\xff\x74\x7b\x39\xbb\x18\x3e\x94\xd3\x27\xb2\xad\xd6\x5f\x1e\x53\xf8\x67\xae\x54\x30\x6f\xac\x46\x87\x4d\x52\xc8\x96\xfb\x9f\x48\x22\x9b\xcd\xc7\xd9\x90\x0e\x12\x03\x9f\x24\x6c\xa5\xff\x46\xc6\xcf\x73\xf3\xbb\x46\x5c\xea\x1c\xd7\x9e\x5d\x22\x69\x64\x66\xce\x1f\x35\x86\xf3\x87\xb2\x07\x8c\xc3\xa7\x85\xb7\x1a\x26\xd6\x00\xac\x41\xdd\xa3\x2e\x9e\x7a\x12\xb2\x3a\x7f\xa9\x35\x39\x3b\xa6\xe1\xdf\x7f\xde\xcd\x2e\xf4\x57\x99\xae\x45\x06\xb2\xac\xce\x9f\x6a\x05\xef\xd4\xb2\xce\x78\x5d\x86\x02\xd2\x61\x73\xb6\x48\x96\x90\x61\x34\x03\x7b\x63\x55\x86\xf3\x77\xbe\x42\x9f\x38\xdb\x1b\xb1\x33\xe2\xc0\x79\x17\x5c\x6c\x89\x4c\x9b\x9d\x3b\x57\x3e\x60\x4e\x9a\xbd\x0b\x13\x2e\xae\x0d\xe9\xdd\x37\x57\x2c\x0c\xbe\x3f\x0d\x62\x9d\x92\xd7\x64\xbd\xb1\x3f\x31\xa3\xef\x6f\xa5\x8f\x29\xf5\x5c\xa9\x6b\x3c\xc6\x12\x85\x8c\x01\xb0\x34\x40\x59\x0a\xe0\xff\x65\x80\xbf\xe7\x00\xfe\x4a\x00\x14\xa2\xdc\x00\x05\x28\xab\xfb\x55\xd6\x29\x7e\xd0\x27\xfa\x46\x06\xb2\x90\x89\xec\x5f\x01\x00\x00\xff\xff\x80\xde\x3b\xc4\xbe\x10\x00\x00")

func iconIcoBytes() ([]byte, error) {
	return bindataRead(
		_iconIco,
		"icon.ico",
	)
}

func iconIco() (*asset, error) {
	bytes, err := iconIcoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "icon.ico", size: 4286, mode: os.FileMode(438), modTime: time.Unix(1775410970, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"icon.ico": iconIco,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"icon.ico": &bintree{iconIco, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
