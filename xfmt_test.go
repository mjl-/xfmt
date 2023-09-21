package xfmt

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestXfmt(t *testing.T) {
	files, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatalf("reading testdata/: %s", err)
	}

	config := Config{
		MaxWidth:      80,
		BreakPrefixes: []string{"- ", "* "},
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".in") {
			continue
		}

		func() {
			outName := "testdata/" + strings.TrimSuffix(f.Name(), ".in") + ".out"
			want, err := os.ReadFile(outName)
			if err != nil {
				t.Errorf("read output file: %s", err)
				return
			}

			inName := "testdata/" + f.Name()
			in, err := os.Open(inName)
			if err != nil {
				t.Errorf("read: %s", err)
				return
			}
			defer in.Close()

			w := &bytes.Buffer{}
			err = Format(w, in, config)
			if err != nil {
				t.Errorf("format: %s", err)
			}
			got := w.String()
			if got != string(want) {
				t.Errorf("got:\n%s\nwant:\n%s", got, string(want))
			}
		}()
	}

	badReader := &errReader{}
	err = Format(&bytes.Buffer{}, badReader, config)
	if err == nil {
		t.Errorf("formatting with failing reader did not return error")
	}

	badWriter := &errWriter{}
	err = Format(badWriter, strings.NewReader("test"), config)
	if err == nil {
		t.Errorf("formatting with failing reader did not return error")
	}
}

type errReader struct {
}

func (r *errReader) Read(buf []byte) (int, error) {
	return 0, errors.New("bad reader")
}

type errWriter struct {
}

func (r *errWriter) Write(buf []byte) (int, error) {
	return 0, errors.New("bad writer")
}
