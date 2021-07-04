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

func TestMemoryAsset(t *testing.T) {
	buf := []byte(`{"foo":"bar"}`)
	asset := mojo.NewAsset(buf)

	if asset.Length() != int64(len(buf)) {
		t.Errorf(`MemoryAsset.Length() incorrect. Got: %d, Expect: %d`, asset.Length(), len(buf))
	}

	w := httptest.NewRecorder()
	err := asset.Serve(w)
	if err != nil {
		t.Errorf(`MemoryAsset.Serve() returned error: %v`, err)
	}
	resBody := w.Body.Bytes()
	if !bytes.Equal(resBody, buf) {
		t.Errorf(`MemoryAsset.Serve() body incorrect. Got: %s, Expect %s`, resBody, buf)
	}

	w = httptest.NewRecorder()
	err = asset.ServeRange(w, 1, 5)
	if err != nil {
		t.Errorf(`MemoryAsset.ServeRange() returned error: %v`, err)
	}
	resBody = w.Body.Bytes()
	if !bytes.Equal(resBody, []byte(`"foo"`)) {
		t.Errorf(`MemoryAsset.ServeRange() body incorrect. Got: %s, Expect %s`, resBody, []byte(`"foo"`))
	}

	asset.AddChunk([]byte("\n"))
	buf = append(buf, byte('\n'))
	if asset.String() != string(buf) {
		t.Errorf("AddChunk() did not append. Expect: %s; Got: %s", buf, asset.String())
	}
}

func TestFileAsset(t *testing.T) {
	buf := []byte(`{"foo":"bar"}`)
	temp := mojo.TempFile()
	temp.Spurt(buf)
	asset := mojo.NewAsset(temp)

	if asset.Length() != int64(len(buf)) {
		t.Errorf(`FileAsset.Length() incorrect. Got: %d, Expect: %d`, asset.Length(), len(buf))
	}

	w := httptest.NewRecorder()
	err := asset.Serve(w)
	if err != nil {
		t.Errorf(`FileAsset.Serve() returned error: %v`, err)
	}
	resBody := w.Body.Bytes()
	if !bytes.Equal(resBody, buf) {
		t.Errorf(`FileAsset.Serve() body incorrect. Got: %s, Expect %s`, resBody, buf)
	}

	w = httptest.NewRecorder()
	err = asset.ServeRange(w, 1, 5)
	if err != nil {
		t.Errorf(`FileAsset.ServeRange() returned error: %v`, err)
	}
	resBody = w.Body.Bytes()
	if !bytes.Equal(resBody, []byte(`"foo"`)) {
		t.Errorf(`FileAsset.ServeRange() body incorrect. Got: %s, Expect %s`, resBody, []byte(`"foo"`))
	}

	asset.AddChunk([]byte("\n"))
	buf = append(buf, byte('\n'))
	if asset.String() != string(buf) {
		t.Errorf("AddChunk() did not append. Expect: %s; Got: %s", buf, asset.String())
	}
}
