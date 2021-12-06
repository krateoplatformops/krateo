package tmpl

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTemplate(t *testing.T) {
	data := map[string]string{
		"User": "projectkerberus",
		"Pass": "XXXX",
	}

	auth := fmt.Sprintf("%s:%s", data["User"], data["Pass"])
	data["Auth"] = base64.StdEncoding.EncodeToString([]byte(auth))

	exptected := "cHJvamVjdGtlcmJlcnVzOlhYWFg="
	if data["Auth"] != exptected {
		t.Fatal(cmp.Diff(exptected, data["Auth"]))
	}
	fmt.Printf("%v\n", data)

	_, err := Execute("dockerconfig.json", data)
	if err != nil {
		t.Fatal(err)
	}
}
