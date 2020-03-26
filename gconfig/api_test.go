package gconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
)

type Data struct {
	Int     int
	String  string
	Bool    bool
	Float64 float64
	SubData []*Data
}

var testData *Data

var subDataLen = 10

func init() {
	testData = &Data{
		Int:     gofakeit.Number(1, 1000),
		String:  gofakeit.Letter(),
		Bool:    gofakeit.Bool(),
		Float64: gofakeit.Float64(),
		SubData: make([]*Data, subDataLen, subDataLen),
	}
	for i := 0; i < subDataLen; i++ {
		subData := &Data{
			Int:     gofakeit.Number(1, 1000),
			String:  gofakeit.Letter(),
			Bool:    gofakeit.Bool(),
			Float64: gofakeit.Float64(),
		}
		testData.SubData[i] = subData
	}
}

func createFileForTest(t *testing.T) *os.File {
	data, err := json.Marshal(testData)
	if err != nil {
		t.Error(err)
	}
	path := filepath.Join(os.TempDir(), fmt.Sprintf("file.%d", time.Now().UnixNano()))
	fh, err := os.Create(path)
	if err != nil {
		t.Error(err)
	}
	_, err = fh.Write(data)
	if err != nil {
		t.Error(err)
	}

	return fh
}

func TestScan(t *testing.T) {
	fh := createFileForTest(t)
	path := fh.Name()
	defer func() {
		fh.Close()
		os.Remove(path)
	}()

	result := []*Data{}
	cfg := NewConfig(path)
	if err := cfg.Scan("SubData", &result); err != nil {
		t.Error(err)
	}

	expected := testData.SubData

	if len(result) != len(expected) {
		t.Fatalf("数量对不上, 预期: %d，结果: %d", len(expected), len(result))
	}

	for i, r := range result {
		e := expected[i]
		if r.Int != e.Int || r.String != e.String || r.Float64 != e.Float64 || r.Bool != e.Bool {
			t.Fatalf("字段数据对不上, 预期: %v，结果: %v", e, r)
		}
	}
}
