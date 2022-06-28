package xrds

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	//"github.com/google/go-cmp/cmp"
)

func TestParserAll(t *testing.T) {
	def, err := ioutil.ReadFile("../../testdata/definition.yaml")
	if err != nil {
		t.Fatal(err)
	}

	fields, err := ParseBytes(def)
	if err != nil {
		t.Fatal(err)
	}

	if got, exp := len(fields), 31; got != exp {
		t.Fatalf("there should be %d fields, got: %d", exp, got)
	}

	//for _, it := range fields {
	//	fmt.Println(it)
	//}
}

func TestParserRequiredOnly(t *testing.T) {
	def, err := ioutil.ReadFile("../../testdata/definition.yaml")
	if err != nil {
		t.Fatal(err)
	}

	fields, err := ParseBytes(def)
	if err != nil {
		t.Fatal(err)
	}

	if got, exp := len(fields), 31; got != exp {
		t.Fatalf("there should be %d fields, got: %d", exp, got)
	}

	for _, it := range fields {
		fmt.Println(it)
	}
}

func TestParseDefault(t *testing.T) {
	stringJSON := &JSON{
		RawMessage: []byte(`"foo"`),
	}

	var got interface{}
	err := json.Unmarshal(stringJSON.RawMessage, &got)
	if err != nil {
		t.Fatal(err)
	}

	if got != "foo" {
		t.Fatalf("raw message should be: foo - got: %s", got)
	}
}
