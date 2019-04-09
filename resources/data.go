// Code generated by go-bindata.
// sources:
// ../build/configs/default_config.yaml
// ../build/configs/testing_config.yaml
// DO NOT EDIT!

package resources

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

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _goCentrifugeBuildConfigsDefault_configYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x58\x59\x73\xdc\xb6\xb2\x7e\x9f\x5f\xd1\x25\xbf\x24\x55\x97\x23\xee\xcb\x54\xa5\x6e\x71\x36\xd9\xb1\xad\x8c\x36\xeb\xc6\x2f\xd7\x20\xd8\x9c\x81\x45\x02\x34\x00\xce\xe2\x5f\x7f\x0a\x20\x47\x1e\x5b\x96\x73\x4e\x52\x27\x2f\xd1\xcb\x88\x0d\x74\xa3\x97\xaf\x3f\x2c\x2f\x60\x8e\x15\xe9\x6a\x0d\x25\x6e\xb1\x16\x6d\x83\x5c\x83\x46\xa5\x39\x6a\x20\x6b\xc2\xb8\xd2\x20\x19\x7f\xc0\xe2\x30\xa2\xc8\xb5\x64\x55\xb7\xc6\x4b\xd4\x3b\x21\x1f\x26\x20\x3b\xa5\x18\xe1\x1b\x56\xd7\x23\x6b\x8c\x71\x04\xbd\x41\x28\x07\xbb\xbc\x9f\xa9\x40\x6f\x88\x86\xd9\xa3\x05\x68\x08\xe3\xda\xd8\x1f\x1d\xa7\x4c\x46\x00\x2f\xe0\x8d\xa0\xa4\xb6\x2e\x30\xbe\x06\x2a\xb8\x96\x84\x6a\x20\x65\x29\x51\x29\x54\xc0\x11\x4b\xd0\x02\x0a\x04\x85\x1a\x76\x4c\x6f\x00\xf9\x16\xb6\x44\x32\x52\xd4\xa8\xc6\x23\x38\xea\x1b\x93\x00\xac\x9c\x40\x10\x04\xf6\x7f\xd4\x1b\x94\xd8\x35\x43\x04\xaf\xca\x09\xa4\x41\xda\x8f\x15\x42\x68\xa5\x25\x69\x57\x88\x52\xf5\xba\x0e\x9c\x9d\xb3\x36\x3c\xf7\xfc\x64\xec\x8e\xdd\xb1\x77\xae\x69\x7b\x1e\xa4\xbe\xeb\x9f\xb3\xb6\x52\xe7\x57\xcd\xed\xd5\xbe\xd8\x3d\x74\xef\x7f\xff\x7d\x5e\x75\x9f\x6f\x8b\xfd\x22\xbf\xc6\xdb\xcb\xd9\x1b\xf1\xf9\x70\x88\xa2\x74\x7b\xc5\xd7\xef\xb6\xab\xb7\x1f\xdf\xfc\xfe\x70\xf6\x07\x46\x83\xa3\xd1\x77\x55\xbc\xb8\x8c\x9b\x87\x4f\xf7\xf8\xf1\xfe\xf5\xbd\xff\x69\xd5\x79\xf1\xff\xb5\xe5\x45\xf0\xf0\xab\xf0\x6e\x83\x66\x43\x36\xab\x69\x74\x83\x11\xf7\x7a\xa3\xc7\x54\xe5\xc7\x4c\xf5\x01\x98\xf0\x91\x6b\xa6\x0f\x4b\x42\xb5\x90\x87\x09\x9c\x9d\x0d\x23\x84\xd3\x8d\x90\xd7\xd8\x0a\xc5\xbe\x19\x62\x7c\x2b\x18\xc5\x3b\xde\x12\x93\xbe\xb3\xb3\x91\xad\xce\x5b\xc2\xf8\x77\xb1\x32\x14\x11\x7e\xba\xee\xc1\xf2\xf3\x08\x4e\xc1\xd1\xfb\xf2\x02\x2e\xbb\x06\x25\xa3\xf0\x6a\x0e\xa2\xb2\x40\x39\x81\xc4\x60\xe3\xb1\x66\x91\x37\x68\x4d\x8f\x85\x81\x9a\x29\x6d\x34\xb9\x28\xf1\x29\xa6\x5a\x29\xb6\xcc\x0e\x08\x6b\xfb\xc4\x81\xa3\xa3\x7f\x58\xe8\x20\x1a\xfb\x7e\x34\xf6\x5d\x77\x1c\xfa\xdf\x16\xdb\xf3\xe7\xc1\x6b\x21\xee\x2f\xd5\x7b\x75\x9f\xdc\x16\xf4\x7d\x94\x5e\x26\xde\xdd\xd5\xcd\xeb\x68\xfe\xf1\xfd\xa7\x66\xf9\xf0\x72\xf5\x72\xb7\x5f\xbe\xbe\xcd\x0f\xe2\xee\x6e\x9e\x96\xd5\xd9\xf7\xcc\xa7\xf1\xd8\xf7\xdc\xe7\xcc\xcf\xd1\x57\xbb\xfb\x45\x50\xf9\xec\xd7\xe2\x0e\xaf\xb2\x8b\xbb\xbb\xab\xe9\xcb\x99\xbc\x7f\x53\x4c\x29\xc9\xde\x5e\xbc\xfd\x54\x35\xc5\x6c\x2d\xbb\xe2\x6c\xc8\xd1\x62\x00\xf6\x63\x25\x5e\xcd\xc1\x81\xa1\x1a\xcf\x41\x3f\x1c\x94\xdf\x10\x93\x1e\x28\xb1\xad\xc5\x01\x4b\xb8\x69\x88\xd4\x30\x1b\x10\xa5\xa0\x12\xd2\x26\x74\xcd\xb6\xc8\xbf\x4a\xe5\x53\xd4\xc1\xb3\xb0\x73\xf7\x55\x9a\xba\x45\x1a\xbb\x9e\x1b\x14\x65\x18\x91\xc8\x0f\xa2\x24\xcc\x11\x67\x6e\x32\x0b\x33\xdf\x0d\xbc\x2a\x4c\x52\xef\x07\x00\x75\xf7\x99\x9f\xcf\xc3\x70\x3a\x4d\x97\x7e\x30\x8f\x4a\xcf\xcf\x70\x9a\xfa\x24\x72\xcb\x20\x8d\xd3\x62\x1a\x16\x1e\xc5\xa5\xb7\x7c\x0e\xca\xee\x9e\x86\x79\x8a\x53\x3f\xa9\xa6\xc1\x82\xf8\x33\x37\x8b\xa2\x65\x4a\xa2\xa9\x17\x7b\xd1\xd4\x8f\xcb\x34\x5a\xce\xa6\x98\xe2\x00\xfa\xd7\x62\x4b\xfa\xa8\x4f\x20\x5a\xa0\xe4\xa4\xde\x20\x5b\x6f\xf4\x00\xa1\x17\x2f\x5e\x0c\xf9\xec\x35\x96\xf9\xd5\xf0\xed\xc0\xbd\xa1\x29\xc6\xab\x4e\x12\x38\x88\x0e\xd6\x86\x5f\x39\xa0\x94\x42\x1a\x70\xdc\x6e\x98\x02\x89\x9f\x3a\xb3\x0a\x53\xc0\x85\x06\xd5\xb5\xad\x90\x1a\x4b\x28\x90\x92\x4e\xa1\xd1\x94\x16\xfb\x66\x8a\xec\x38\x37\x1c\x69\x19\x50\x69\xa2\x4d\x03\x74\x46\x34\x86\xeb\x8e\xf7\x72\xc7\x19\x64\xbf\x10\x49\x37\x6c\x8b\xe3\xb3\xff\x19\x9c\x02\xd8\x99\xfe\xd1\x02\x4a\xf1\xbf\x56\x83\x40\x6d\xd9\xb7\x25\x92\xe9\x43\xbf\x90\xb5\xf2\x60\xe3\xc1\xf5\xa4\xff\xfc\x30\x4c\x70\x1c\xba\x21\x8c\xff\xd2\x0f\x3b\x8e\xf1\xf6\x97\xc0\x0d\xdc\x10\x1c\x67\x47\x64\x3b\xfc\x38\x05\x91\x92\xa1\x84\x28\x4e\x5d\xd7\x75\xc1\x71\xb8\x70\x08\xa7\x0c\xb9\x76\x8a\x5a\xd0\x07\xd5\xcb\x14\xca\x2d\x3a\xb5\x49\x2a\x38\x4e\x43\xf6\x4e\x6b\x5a\x14\xfc\xc8\x28\x29\x4e\x5a\xb5\x11\x7a\x10\x5a\x59\xc3\xf8\x57\x9f\xc6\x67\x42\x35\xdb\x22\x38\x8e\x81\xa6\x49\x91\xa8\xaa\xa7\x99\x00\xc7\x29\x0b\x87\x8a\xa6\x35\xf3\x05\x07\xa5\x4a\x13\x12\xa1\x1b\x74\x14\xfb\x8c\x10\xba\x59\x0c\x8e\xf3\x51\x09\x2e\x5b\xea\x6c\x84\xd2\x0a\x48\x5d\x9f\xc8\x18\xd7\x28\x2b\x42\xd1\xc8\x3f\x7c\x5d\xee\xa7\xc9\xfc\x5e\xe5\xa7\x26\x7c\x2c\x4d\x27\x71\xec\x1d\xd1\x02\xee\xb1\xb8\x31\x72\xad\xc0\xe6\x44\x42\x25\x45\x03\x1d\xd7\xb2\x53\x06\x12\x42\xb2\x35\xe3\x13\x18\x8f\xcf\x9e\xad\xa7\x69\xd9\x27\xb5\xfc\xe0\x38\x1d\x57\xa4\x42\x07\xf7\xad\x50\xf8\x01\xaa\x9a\xac\xbf\x01\xf0\x7f\xc6\xd3\xfe\x5f\xe4\xe9\xaf\x7a\xe9\xdf\x66\x6a\xcf\x0d\xc7\x5e\x14\x8e\xbd\x74\x1c\x3d\xd9\x96\x8f\x54\x7a\x53\xec\x8b\xd7\xb3\xe2\xfd\x26\xfb\xf5\x9d\x56\x57\x87\x77\x17\xe5\xed\x4a\x92\xf0\xba\xbd\xc9\x43\x5d\x6c\x55\x4c\xb8\xe7\x7d\xdc\x5d\xe4\xfe\xe7\xb3\xef\x98\x8f\xc6\x5e\x1a\x8d\xfd\x20\x79\x6e\x81\xab\xc6\xa7\x37\x8d\x5c\x30\x72\xf3\xf6\x5d\xb8\xbe\xdb\x26\xf7\x17\x9b\x76\x7d\xbd\x13\xe9\x4e\x2c\x6f\xd4\xcb\xcd\xfb\x8b\xe2\x82\x05\x24\x4f\xf7\x3f\xe6\x6a\x4b\x1a\xcf\x32\xb5\xff\x5f\xa0\xea\x1f\x30\xb5\x17\xc4\xfe\x82\x4e\xab\x34\x4e\x32\x3f\x0c\x16\x7e\x58\xe5\xee\x62\x16\xfa\x51\xe9\xa3\xe7\xe6\x6e\xea\xfb\x01\x4d\xe6\x3f\x64\xea\xc4\x4b\xdd\x79\x92\x04\x9e\x5b\x22\x4d\xf3\xa9\x9f\xe6\x24\x75\xfd\x05\x75\xb3\x65\x95\xfb\xf3\x65\x1c\x62\xe6\x26\xf4\x79\xa6\xf6\xd2\xc0\x4b\xdc\x30\xf5\xe2\x30\xad\xb0\xaa\x30\xcc\x42\x77\x19\xcc\xf3\xbc\x0c\x48\x52\xd0\xa2\x70\x69\x94\xe7\xcb\x81\xa9\xaf\x45\xab\x34\x3e\xe1\xea\x52\xac\x5b\xa2\xe9\xe6\xcf\x9d\x42\x82\xbf\x88\xee\xe3\xea\xf0\xd3\xed\x6f\xf3\xdf\x80\x4a\x34\x54\x2d\x07\x57\x0d\xc2\xad\x9d\x9f\xff\x69\x47\x93\x3e\x01\xcf\x01\x3e\xf8\x7b\xf1\xee\x96\x41\xe6\x2d\x12\x3f\xf0\xa3\x19\x96\xb3\xd0\x5b\x84\xa9\x1b\x05\x8b\x24\xf1\xd3\x94\xa4\xd9\xd2\x5f\x04\x9e\xe7\x45\x3f\xc4\xbb\x3f\x4b\xdd\xa5\x37\x27\xd5\x9c\x24\x24\x9f\x63\xe1\xcf\xbc\x24\x2a\xc3\x69\x18\xe4\x69\x94\x86\x49\xb0\xf0\xbc\xc4\x0b\x9e\xc7\xfb\x22\x09\xd2\xd0\xf3\x42\x77\x49\x31\x0b\x32\xcf\x5f\xd0\x24\xce\xca\xbc\x4c\x97\x61\xea\x4f\xab\x88\x26\x71\x3a\xcf\xa3\xd3\xe3\xb8\x39\x7e\x7f\x03\x78\x6c\x0a\x22\x29\x29\x51\x0a\x33\xf2\xa7\x70\xef\xb9\xff\x30\x40\x9a\x1b\xe9\x11\x3b\xdf\x01\xa4\xf7\xf7\x02\x32\x5f\x86\x51\x4c\xbd\x38\x48\x63\x12\x87\x55\x19\x2e\xc3\x22\xce\x48\xe5\x05\x24\x8d\xe7\x95\x3b\x8d\x62\x3f\x27\xae\xfb\x43\x40\xc6\x41\x32\x4d\x67\xc1\xdc\xcf\xf3\x60\x46\x7d\x37\x9e\x67\x61\xe4\x65\x45\x14\xa6\x99\xef\xa6\x19\xcd\x16\x71\x92\x65\xee\xf3\x80\x9c\x46\x18\xfa\x41\x39\xa3\x49\xe8\x16\xd3\x59\xea\x56\x99\x1b\x7b\x41\x80\x5e\x14\xbb\x5e\x95\xa5\x6e\x96\xa5\x41\x14\x9f\x8d\xcc\xfd\x9f\x68\x02\x37\x5a\x48\xb2\xc6\x91\xea\x7f\xfb\x5b\xfd\x8a\xe8\x8d\xcd\x4c\x6d\xee\x8e\xf3\x29\x54\xac\xc6\x11\x40\x4b\xf4\x66\x02\xe7\xba\x69\xcf\xbf\xbc\x2e\xfc\x7f\x49\x34\x19\xdb\x99\x65\x61\xec\xce\x04\xaf\xd8\xba\x93\xc4\x9e\x8e\x8e\x0b\x50\x2b\xbd\xf9\xf3\xcb\xf4\x06\x9e\xac\x96\x53\x2a\x3a\xae\x15\x3c\xe0\x01\x86\x28\x46\x64\x10\x9a\x75\x1e\xf0\x60\xc4\x38\x58\x3c\x0e\x19\xdd\x57\x8f\xc7\xc1\x9d\x01\x90\x05\x42\xbe\x7a\x05\x84\x97\xb0\xf2\x57\x70\xd3\x9f\xe5\xcc\x06\x82\xdc\xec\x10\x23\xc3\xfd\x2f\x85\xd2\x9c\x34\x38\x81\xc7\x17\x81\xd1\x0b\x58\x09\xa9\x07\x33\xc6\xc4\xf7\x55\xcd\xa4\x09\xa4\x6e\xea\x9b\xe5\x4d\x8f\x3a\x5a\xd8\x03\x31\xd0\xd3\xac\xa9\x51\xeb\xb7\x7d\x92\x6e\x5a\xa4\xac\x3a\xc0\x62\xaf\xed\xb9\x0b\x5e\xad\x4e\xbc\xb5\x07\x45\x4a\x38\x14\x08\x12\xcd\x59\xb8\x04\xa2\x81\x55\x50\xe0\x86\xf1\x12\x2e\xf3\x5b\x63\x06\x07\xed\x57\xab\x09\xec\xc6\xfb\xf1\x61\xfc\xb9\x2f\x81\xf1\xba\x53\x58\x3e\x36\x82\x89\xbb\x26\x07\x94\xa6\x10\xd6\x5d\xdb\xc6\x76\xf6\x2d\x6b\x50\x74\x36\x4c\x0e\xa2\x45\x3e\x3c\xfa\x0c\x27\x61\xbb\x87\xda\xd3\xfd\x08\x8e\xe2\x41\x65\x02\x67\x81\xab\x2c\xec\xae\x3a\xec\xf0\x9b\x70\xed\xea\x44\x1d\x38\xdd\x48\xc1\x45\xa7\xcc\xb6\x4c\x51\x29\xc6\xd7\xa3\x4f\x46\xa1\x4f\x46\xff\x64\xa5\xfa\xd0\xbb\xa6\x40\x69\xa8\xd1\xf4\x3c\x4a\x75\x4e\x05\x57\x86\x33\x87\x4d\x7e\xc7\xea\xda\xe4\x85\xd4\xe6\x68\xaf\xfb\xcc\x28\x4d\xa4\xee\xda\x11\x18\xfd\xfb\x5e\xd1\xb0\xa7\x6b\xed\x2f\x25\xa2\x82\xae\x85\xd9\xea\x0e\xe8\x81\xd6\xa8\xfa\x60\xfb\x25\xcc\x3d\x6e\x47\x98\x7d\xeb\x32\x1e\xe3\x16\x0d\x92\x60\x18\xbe\x27\xcc\xc6\xfb\xf6\xa6\xe7\x1f\x4b\xe1\x83\x8f\x12\xb5\x64\x68\xef\x23\x62\x37\xa4\x9b\x80\x26\xca\x50\xb8\xf9\xb9\xee\x27\x58\x26\x1f\x9d\x90\x9e\xb2\xf5\x67\xf4\xeb\x8c\x8d\x8e\x94\x37\x80\x04\x6b\x34\x6c\xb6\xdb\x30\xba\x79\xa4\x43\x18\xb0\x6e\xca\x62\xee\xa3\xc3\x2e\x22\x4c\x06\x87\x63\x4f\x09\xac\xbf\x78\xd0\x4e\x69\xd1\x0c\x8b\x1c\x1b\x71\x78\x16\x1c\x5a\xec\xd2\x62\xfe\xcc\x10\xef\xd9\xe3\xe3\x9f\xed\xf1\xc1\xf0\xe3\xba\xb4\x36\x57\xc5\x1e\x9c\x3f\xed\xd0\xde\x94\x99\x44\xd8\x29\x10\x12\x58\x4b\x87\x17\x41\x52\xd4\x68\xfe\xa5\xf6\xc0\xd5\x67\xd3\x1c\xac\x8c\xe2\xdd\xf5\x9b\x09\x6c\xb4\x6e\x27\xe7\xe7\xf6\x6a\x66\xee\x73\x93\x2c\x0a\x23\xbb\x76\x43\xf6\xac\x31\x21\x0e\xf9\x5c\x13\x13\x13\xa3\xd6\x5e\x4b\x0e\xc7\x04\x4b\xc2\xd5\x70\x61\x64\x1c\x76\xc8\xac\xb6\xef\xc2\xc5\x0e\x19\x70\xb1\x1b\x81\xb1\x75\x41\xd4\xca\x68\x4f\xc0\x77\x1f\xff\xec\xd4\x0b\xa2\xa0\x66\x0d\x1b\xf6\x8a\x92\x55\x15\x4a\x13\xdd\x63\x85\x44\x8b\xc7\xae\x05\xe3\xc7\x1b\x3b\xfb\xf8\x98\x39\xb3\x07\x48\x0b\xb1\xc1\xa6\x91\xe6\x65\xf9\x1a\x0f\x13\x08\x4e\x85\xd7\xb8\x15\x0f\x68\xe5\x51\x74\x14\xf7\x3b\xc5\x4c\x34\x0d\x33\xd4\xf1\x8d\x7c\x25\xf1\x38\xe4\x7d\x31\xc5\x2b\xfd\x96\x71\x3d\x81\xec\x4b\x1c\xc7\xde\xd5\xc2\x42\xb8\xcf\x0f\xff\x52\xb3\xd3\x4c\x0d\xd5\x29\xcb\xfe\xed\x96\x80\xbd\xf4\x5b\x5a\xec\x8b\x04\x5a\xb2\xf5\x1a\x25\x96\x7d\xa7\x6b\xdc\xeb\x23\xfa\xfb\x6e\x8f\x5d\xd3\xee\xcf\x2d\x2c\x91\x94\x20\x78\x7d\x38\x49\xde\xe3\x03\xf6\xd1\xa5\x2f\xa6\xaf\x91\x94\x5f\x9b\xf7\xa2\xc1\xfa\xa5\xc1\xd8\xa9\xef\xad\x10\xb5\xa9\xe8\x63\xc7\x69\x01\x0a\x79\xf9\x0d\x18\xc4\xd6\x32\x5c\x43\xf6\x8f\x8d\xe7\x0f\x99\xfa\xbe\x49\xfb\x74\xb0\x25\xb5\xb5\x7b\xe8\x59\x81\x18\x07\x69\x27\x2d\x1e\x4e\x35\x36\x44\x41\x81\xc8\xa1\x44\x8d\x54\xdb\x34\x1d\x0d\x98\xf5\xcc\x7e\xef\x0f\x11\xcc\x99\xb2\x7d\x60\x2d\x2a\xd1\x3c\xe9\x23\x05\xa5\x38\x7d\x61\x02\xbd\xb7\x1e\x91\xd6\x80\x59\xef\x57\x42\xd4\x39\x35\x6c\xb9\xe0\xc6\x52\x39\x01\x2d\x3b\x34\x2c\x42\xf8\x01\x4a\x2c\xba\xf5\x7a\x60\x6a\xd3\xdc\x96\x17\xd7\x02\xcc\x22\x23\x3b\xda\x93\x48\xdb\x4a\x51\xf5\x20\x3f\xaa\x98\x3d\xc0\x48\x27\x50\x91\x5a\xe1\x68\xd4\xa3\x6e\x78\xab\x6f\x25\xd2\x01\x7c\x76\xc1\x7f\x05\x00\x00\xff\xff\x4c\x50\x3e\x8b\xa0\x18\x00\x00")

func goCentrifugeBuildConfigsDefault_configYamlBytes() ([]byte, error) {
	return bindataRead(
		_goCentrifugeBuildConfigsDefault_configYaml,
		"go-centrifuge/build/configs/default_config.yaml",
	)
}

func goCentrifugeBuildConfigsDefault_configYaml() (*asset, error) {
	bytes, err := goCentrifugeBuildConfigsDefault_configYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "go-centrifuge/build/configs/default_config.yaml", size: 6304, mode: os.FileMode(420), modTime: time.Unix(1554810360, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _goCentrifugeBuildConfigsTesting_configYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x53\xb9\x8e\xe4\x36\x10\xcd\xf5\x15\x04\x1d\x4c\xd2\x07\x4f\xf1\xc8\x1c\x1a\x0b\x3b\xb1\x81\x8d\x8b\x64\xb1\x87\xe8\xd1\x61\x92\x9a\xd9\xc6\x62\xff\xdd\x50\x6f\x8f\xbd\x99\x67\xb3\xaa\xd2\x3b\xaa\xa4\xa7\x88\x73\xaf\x25\x6f\x17\xfc\x03\xfb\xdb\x52\xaf\x9e\x74\x6c\xbd\xcc\x97\x01\xfb\x33\x56\xdc\x26\x3f\x10\x02\x31\x2e\xdb\xdc\xdb\x5e\x13\x32\x41\x99\x3d\xb9\x97\x84\x5c\xf1\xe6\xc9\xd3\x57\x0a\x29\x55\x6c\x8d\x7a\x6a\x5d\x60\x60\x47\x6d\x65\x54\x4a\x29\x88\x39\x19\x1e\xd4\x28\x91\x25\x19\xb5\x06\xe4\x8a\x0b\xd0\xf4\x40\x63\xbd\xad\x7d\xa1\xfe\x2b\x8d\x65\x7d\xc6\x4a\x3d\x05\x6c\x47\x2e\xec\x31\xf6\xba\x03\xee\xe3\x8e\x5f\x3a\xf5\x34\x1a\xe3\xb2\x95\xc6\x25\x63\x58\x72\x22\xe6\xc8\x53\x4a\x0a\x6c\x96\x3c\x69\x60\x90\xa2\xcd\x02\x58\x10\xc0\x15\xe3\xd2\xb0\x24\x47\xc9\xb2\xb4\x91\x45\x0b\xff\xea\xad\x50\x61\x6a\xbb\x6d\x79\xa5\x9e\xca\x31\xf2\xd1\xa2\x91\x21\x3b\xcb\x32\x1a\x1d\x98\x11\x26\x5b\xc7\xc0\x70\x48\xf4\xdb\x81\x5e\x53\xa6\x9e\xb6\xfb\xc2\xf4\xde\xfe\x27\x92\xae\x2f\x38\x53\x2f\xc5\x81\xce\xd4\x8b\x51\x70\xa5\x0e\x74\xa5\x9e\x1f\x68\xa5\xde\x1e\x68\x83\x97\xfd\x80\x84\x3c\x20\x1f\x51\x46\x67\xb9\x53\x2a\x71\x8c\x20\x82\x0d\xc2\xa0\xc2\x11\x59\xd0\x21\x07\x25\x03\x32\x69\x46\xd0\xc9\x5a\xeb\x32\x8c\xc6\x81\xb0\x5c\x88\x7d\x91\x09\xe2\xfe\x2a\x22\x17\x36\x58\xae\xb5\xd6\x01\x38\x42\x32\x11\xd0\xb1\x91\xa1\xb5\x4a\x40\x8e\x60\xa5\x1e\x13\x1b\x95\xd6\x21\x39\xd0\x46\x8b\x00\x63\x8e\x91\x39\x81\x79\x57\x2a\x89\x7a\xaa\x34\xb2\x91\xc1\x78\x4c\x02\xf0\xa8\x64\xb0\x47\x27\x44\x3e\x2a\x65\x85\x53\xce\x25\x69\x12\x3d\xd0\x57\xac\xad\x2c\xfb\x91\xdf\x9e\x1e\x1f\x7e\x85\xd6\xde\x96\x9a\x3c\x79\x7a\x1f\x3d\x32\xe0\xc9\x47\x23\x30\x0c\x25\xe1\xdc\x4b\xbf\xfd\x96\x3c\xa1\xec\xcb\x87\xb3\x33\x0c\xbf\x90\x5f\x1f\xa9\xdc\x33\x48\x5a\x5f\x2a\x5c\x70\xf8\x31\xaa\x57\xbc\xed\x63\xf4\xe4\xdc\xa7\xf5\xfc\xfe\x68\x18\xfe\xde\x70\xc3\x1d\x31\x6f\xd3\xe7\xa5\x5e\xb1\x36\x4f\xc4\x40\xc8\xdb\xbd\xf9\x0c\xa5\xff\x55\x26\xfc\xfd\x4f\x4f\xf8\x30\xec\x32\x3b\x78\x15\xeb\xf7\x1f\x60\xdd\xc2\x4b\x89\x9f\xf6\xe4\x9f\x4e\xe7\xd3\xe9\x1c\xb6\xf2\x92\xce\x15\xdb\xb2\xd5\x88\xed\xbc\x8a\xf5\x13\xde\x4e\xeb\x16\x4e\x2b\x4e\xdf\x39\xb5\xbc\x42\xc7\xff\x27\x5d\x77\xe2\x9d\xd4\xca\x65\x2e\xf3\xe5\x83\x9e\x0f\xf4\xcf\xfb\xfe\x40\x7c\xf7\x1e\x60\x8e\xcf\x4b\x7d\x98\xaf\x15\xe3\x32\x4d\xa5\x7b\xd2\xeb\x86\xff\x04\x00\x00\xff\xff\x1f\xaf\xbe\x5d\x34\x04\x00\x00")

func goCentrifugeBuildConfigsTesting_configYamlBytes() ([]byte, error) {
	return bindataRead(
		_goCentrifugeBuildConfigsTesting_configYaml,
		"go-centrifuge/build/configs/testing_config.yaml",
	)
}

func goCentrifugeBuildConfigsTesting_configYaml() (*asset, error) {
	bytes, err := goCentrifugeBuildConfigsTesting_configYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "go-centrifuge/build/configs/testing_config.yaml", size: 1076, mode: os.FileMode(420), modTime: time.Unix(1552945185, 0)}
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
	"go-centrifuge/build/configs/default_config.yaml": goCentrifugeBuildConfigsDefault_configYaml,
	"go-centrifuge/build/configs/testing_config.yaml": goCentrifugeBuildConfigsTesting_configYaml,
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
	"go-centrifuge": &bintree{nil, map[string]*bintree{
		"build": &bintree{nil, map[string]*bintree{
			"configs": &bintree{nil, map[string]*bintree{
				"default_config.yaml": &bintree{goCentrifugeBuildConfigsDefault_configYaml, map[string]*bintree{}},
				"testing_config.yaml": &bintree{goCentrifugeBuildConfigsTesting_configYaml, map[string]*bintree{}},
			}},
		}},
	}},
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

