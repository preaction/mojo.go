package mojo

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Asset is a container for a Message's contents
type Asset interface {
	Length() int64
	Serve(http.ResponseWriter) error
	ServeRange(w http.ResponseWriter, start int64, end int64) error
	String() string
	AddChunk([]byte)
}

// NewAsset builds an Asset from the given content, which can be a File
// object, a string, or an array of bytes.
func NewAsset(content interface{}) Asset {
	switch v := content.(type) {
	case File:
		return &FileAsset{path: v}
	case *os.File:
		return &FileAsset{path: NewFile(v.Name())}
	case string:
		return &MemoryAsset{buffer: []byte(v)}
	case []byte:
		return &MemoryAsset{buffer: v}
	}
	panic(fmt.Sprintf("mojo.NewAsset: Unknown type %t", content))
}

// FileAsset is an Asset backed by a file on a filesystem
type FileAsset struct {
	path File
}

// Length returns the length of the file
func (asset *FileAsset) Length() int64 {
	return asset.path.Stat().Size()
}

// Serve writes the entire file contents to the given http.ResponseWriter
func (asset *FileAsset) Serve(w http.ResponseWriter) error {
	// XXX: Go makes us copy/paste because receiver can't be interface?
	// XXX: Does calculating the length slow this down?
	return asset.ServeRange(w, 0, asset.Length()-1)
}

// ServeRange writes the given range to the given http.ResponseWriter.
// start and end are byte positions of the first and last byte to send,
// respectively.
func (asset *FileAsset) ServeRange(w http.ResponseWriter, start int64, end int64) error {
	file := asset.path.Open()
	if _, err := file.Seek(start, io.SeekStart); err != nil {
		panic(fmt.Sprintf("Could not Seek to %d: %v", start, err))
	}
	// XXX: Check that all requested was written?
	_, _ = io.CopyN(w, file, 1+end-start)
	return nil
}

// String returns the asset's contents as a string
func (asset *FileAsset) String() string {
	content, err := os.ReadFile(asset.path.String())
	if err != nil {
		panic(fmt.Sprintf("Could not read file %s: %v", asset.path, err))
	}
	return string(content)
}

// AddChunk adds the given data to the end of the file
func (asset *FileAsset) AddChunk(data []byte) {
	file := asset.path.Open()
	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		panic(fmt.Sprintf("Could not Seek to end: %v", err))
	}
	// XXX: Check that all requested was written?
	bytes, err := file.Write(data)
	if bytes != len(data) || err != nil {
		panic(fmt.Sprintf("Error writing %d bytes: %v", len(data), err))
	}
}

// MemoryAsset is an Asset stored in memory
type MemoryAsset struct {
	buffer []byte
}

// Length returns the length of the buffer
func (asset *MemoryAsset) Length() int64 {
	return int64(len(asset.buffer))
}

// Serve writes the entire buffer contents to the given http.ResponseWriter
func (asset *MemoryAsset) Serve(w http.ResponseWriter) error {
	// XXX: Go makes us copy/paste because receiver can't be interface?
	// XXX: Does calculating the length slow this down?
	return asset.ServeRange(w, 0, asset.Length()-1)
}

// ServeRange writes the given range to the given http.ResponseWriter.
// start and end are byte positions of the first and last byte to send,
// respectively.
func (asset *MemoryAsset) ServeRange(w http.ResponseWriter, start int64, end int64) error {
	_, err := w.Write(asset.buffer[start : end+1])
	// XXX: Check that all requested was written?
	return err
}

// String returns the asset's contents as a string
func (asset *MemoryAsset) String() string {
	return string(asset.buffer)
}

// AddChunk adds the given data to the end of the buffer
func (asset *MemoryAsset) AddChunk(data []byte) {
	asset.buffer = append(asset.buffer, data...)
}
