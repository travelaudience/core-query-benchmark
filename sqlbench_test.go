package sqlbench

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	_, e := New("file_does_not_exists.json")
	if e == nil {
		t.Error("Expected an error when file does not exists", e)
	}

	_, e = New(createJSONFile("{}"))
	if e != nil {
		t.Error("Expected a new object with no error", e)
	}

	_, e = New(createJSONFile("{"))
	if e == nil {
		t.Error("Expected a error if json is not correct", e)
	}

}

func createJSONFile(s string) string {
	content := []byte(s)
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	return tmpfile.Name()
}
