package mojo

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
)

// Asset is a container for a Message's contents
type Asset interface {
	Length() int64
	Range(start int64, end int64)
	Serve(w http.ResponseWriter) error
	String() string
	AddChunk([]byte)
}

// NewAsset builds an Asset from the given content, which can be a File
// object, a string, or an array of bytes.
func NewAsset(content interface{}) Asset {
	switch v := content.(type) {
	case File:
		return &FileAsset{path: v.String(), file: v.Open()}
	case *os.File:
		return &FileAsset{path: v.Name(), file: v}
	// XXX: If this doesn't support Stat() or Seek(), it's probably not
	// a real file and should be slurped into a buffer and made into
	// a MemoryAsset. If `fs.File` ever supports Stat() or Seek(), then
	// we can remove this restriction.
	case fs.File:
		return &FileAsset{file: v}
	case string:
		return &MemoryAsset{buffer: []byte(v)}
	case []byte:
		return &MemoryAsset{buffer: v}
	}
	panic(fmt.Sprintf("mojo.NewAsset: Unknown type %t", content))
}

// FileAsset is an Asset backed by a file on a filesystem
type FileAsset struct {
	path     string
	file     fs.File
	hasRange bool
	start    int64
	end      int64
}

// Length returns the length of the file
func (asset *FileAsset) Length() int64 {
	stat, err := asset.file.Stat()
	if err != nil {
		panic(fmt.Sprintf("Could not Stat(): %v", err))
	}
	return stat.Size()
}

// Range sets a start/end range to serve partial content
func (asset *FileAsset) Range(start int64, end int64) {
	asset.start = start
	asset.end = end
	asset.hasRange = start >= 0 || end >= 0
}

// Serve writes the requested contents to the given http.ResponseWriter
func (asset *FileAsset) Serve(w http.ResponseWriter) error {
	file := asset.open()
	// XXX: Check that all requested was written?
	_, _ = io.Copy(w, file)
	return nil
}

// open opens the file to the correct point
func (asset *FileAsset) open() io.Reader {
	start := int64(0)
	end := int64(-1)
	if asset.hasRange {
		start = asset.start
		end = asset.end
	}

	file := asset.file.(io.ReadSeeker)
	if _, err := file.Seek(start, io.SeekStart); err != nil {
		panic(fmt.Sprintf("Could not Seek to %d: %v", start, err))
	}

	if end != -1 {
		return &io.LimitedReader{file, end - start + 1}
	}
	return file
}

// String returns the asset's contents as a string
func (asset *FileAsset) String() string {
	content, _ := io.ReadAll(asset.open())
	return string(content)
}

// AddChunk adds the given data to the end of the file
func (asset *FileAsset) AddChunk(data []byte) {
	file, ok := asset.file.(io.WriteSeeker)
	if !ok {
		var err error
		file, err = os.OpenFile(asset.path, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Sprintf("Could not open file for appending: %v", err))
		}
	}
	_, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		panic(fmt.Sprintf("Could not Seek to end: %v", err))
	}
	// Check that all requested was written
	bytes, err := file.Write(data)
	if bytes != len(data) || err != nil {
		panic(fmt.Sprintf("Error writing %d bytes: %v", len(data), err))
	}
}

// MemoryAsset is an Asset stored in memory
type MemoryAsset struct {
	buffer   []byte
	hasRange bool
	start    int64
	end      int64
}

// Length returns the length of the buffer
func (asset *MemoryAsset) Length() int64 {
	return int64(len(asset.buffer))
}

// Range sets a start/end range to serve partial content
func (asset *MemoryAsset) Range(start int64, end int64) {
	asset.start = start
	asset.end = end
	asset.hasRange = start >= 0 || end >= 0
}

// Serve writes the requested contents to the given http.ResponseWriter
func (asset *MemoryAsset) Serve(w http.ResponseWriter) error {
	_, err := w.Write([]byte(asset.String()))
	return err
}

// String returns the asset's contents as a string
func (asset *MemoryAsset) String() string {
	buffer := asset.buffer
	if asset.hasRange {
		buffer = buffer[asset.start : asset.end+1]
	}
	return string(buffer)
}

// AddChunk adds the given data to the end of the buffer
func (asset *MemoryAsset) AddChunk(data []byte) {
	asset.buffer = append(asset.buffer, data...)
}
