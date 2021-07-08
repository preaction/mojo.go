package mojo_test

import (
	"bytes"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/preaction/mojo.go"
)

func TestNewAsset(t *testing.T) {
	asset := mojo.NewAsset([]byte("foobar"))
	if _, ok := asset.(*mojo.MemoryAsset); !ok {
		t.Errorf("NewAsset([]byte) did not return MemoryAsset. Got: %t", asset)
	}

	asset = mojo.NewAsset("foobar")
	if _, ok := asset.(*mojo.MemoryAsset); !ok {
		t.Errorf("NewAsset(string) did not return MemoryAsset. Got: %t", asset)
	}

	file, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatalf("Could not create temp file: %v", err)
	}
	defer os.RemoveAll(file.Name())
	asset = mojo.NewAsset(file)
	if _, ok := asset.(*mojo.FileAsset); !ok {
		t.Errorf("NewAsset(os.File) did not return FileAsset. Got: %t", asset)
	}

	path := mojo.NewFile(file.Name())
	asset = mojo.NewAsset(path)
	if _, ok := asset.(*mojo.FileAsset); !ok {
		t.Errorf("NewAsset(mojo.File) did not return FileAsset. Got: %t", asset)
	}
}

func TestMemoryAssetLength(t *testing.T) {
	buf := []byte(`{"foo":"bar"}`)
	asset := mojo.NewAsset(buf)

	if asset.Length() != int64(len(buf)) {
		t.Errorf(`MemoryAsset.Length() incorrect. Got: %d, Expect: %d`, asset.Length(), len(buf))
	}
}

func TestMemoryAssetServe(t *testing.T) {
	buf := []byte(`{"foo":"bar"}`)
	asset := mojo.NewAsset(buf)
	w := httptest.NewRecorder()
	err := asset.Serve(w)
	if err != nil {
		t.Errorf(`MemoryAsset.Serve() returned error: %v`, err)
	}
	resBody := w.Body.Bytes()
	if !bytes.Equal(resBody, buf) {
		t.Errorf(`MemoryAsset.Serve() body incorrect. Got: %s, Expect %s`, resBody, buf)
	}
}

func TestMemoryAssetAddChunk(t *testing.T) {
	buf := []byte(`{"foo":"bar"}`)
	asset := mojo.NewAsset(buf)
	asset.AddChunk([]byte("\n"))
	buf = append(buf, byte('\n'))
	if asset.String() != string(buf) {
		t.Errorf("AddChunk() did not append. Expect: %s; Got: %s", buf, asset.String())
	}
}

func TestMemoryAssetRange(t *testing.T) {
	buf := []byte(`{"foo":"bar"}`)
	asset := mojo.NewAsset(buf)
	w := httptest.NewRecorder()
	asset.Range(1, 5)
	err := asset.Serve(w)
	if err != nil {
		t.Errorf(`MemoryAsset.Serve() with range returned error: %v`, err)
	}
	resBody := w.Body.Bytes()
	if !bytes.Equal(resBody, []byte(`"foo"`)) {
		t.Errorf(`MemoryAsset.Serve() with range body incorrect. Got: %s, Expect %s`, resBody, []byte(`"foo"`))
	}
}

func TestFileAssetLength(t *testing.T) {
	buf := []byte(`{"foo":"bar"}`)
	temp := mojo.TempFile()
	temp.Spurt(buf)
	asset := mojo.NewAsset(temp)

	if asset.Length() != int64(len(buf)) {
		t.Errorf(`FileAsset.Length() incorrect. Got: %d, Expect: %d`, asset.Length(), len(buf))
	}
}

func TestFileAssetServe(t *testing.T) {
	buf := []byte(`{"foo":"bar"}`)
	temp := mojo.TempFile()
	temp.Spurt(buf)
	asset := mojo.NewAsset(temp)

	w := httptest.NewRecorder()
	err := asset.Serve(w)
	if err != nil {
		t.Errorf(`FileAsset.Serve() returned error: %v`, err)
	}
	resBody := w.Body.Bytes()
	if !bytes.Equal(resBody, buf) {
		t.Errorf(`FileAsset.Serve() body incorrect. Got: %s, Expect %s`, resBody, buf)
	}
}

func TestFileAssetAddChunk(t *testing.T) {
	buf := []byte(`{"foo":"bar"}`)
	temp := mojo.TempFile()
	temp.Spurt(buf)
	asset := mojo.NewAsset(temp)

	asset.AddChunk([]byte("\n"))
	buf = append(buf, byte('\n'))
	if asset.String() != string(buf) {
		t.Errorf("AddChunk() did not append. Expect: %s; Got: %s", buf, asset.String())
	}
}

func TestFileAssetRange(t *testing.T) {
	buf := []byte(`{"foo":"bar"}`)
	temp := mojo.TempFile()
	temp.Spurt(buf)
	asset := mojo.NewAsset(temp)

	w := httptest.NewRecorder()
	asset.Range(1, 5)
	err := asset.Serve(w)
	if err != nil {
		t.Errorf(`FileAsset.Serve() with range returned error: %v`, err)
	}
	resBody := w.Body.Bytes()
	if !bytes.Equal(resBody, []byte(`"foo"`)) {
		t.Errorf(`FileAsset.Serve() with range body incorrect. Got: %s, Expect %s`, resBody, []byte(`"foo"`))
	}
}
