package tmpl

import (
	"testing"
)

func TestDefault(t *testing.T) {
	tpl := `{{"" | default "foo"}}`
	if err := runt(tpl, "foo"); err != nil {
		t.Error(err)
	}
	tpl = `{{default "foo" 234}}`
	if err := runt(tpl, "234"); err != nil {
		t.Error(err)
	}
	tpl = `{{default "foo" 2.34}}`
	if err := runt(tpl, "2.34"); err != nil {
		t.Error(err)
	}

	tpl = `{{ .Nothing | default "123" }}`
	if err := runt(tpl, "123"); err != nil {
		t.Error(err)
	}
	tpl = `{{ default "123" }}`
	if err := runt(tpl, "123"); err != nil {
		t.Error(err)
	}
}

func TestEmpty(t *testing.T) {
	tpl := `{{if empty 1}}1{{else}}0{{end}}`
	if err := runt(tpl, "0"); err != nil {
		t.Error(err)
	}

	tpl = `{{if empty 0}}1{{else}}0{{end}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}
	tpl = `{{if empty ""}}1{{else}}0{{end}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}
	tpl = `{{if empty 0.0}}1{{else}}0{{end}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}
	tpl = `{{if empty false}}1{{else}}0{{end}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}

	dict := map[string]interface{}{"top": map[string]interface{}{}}
	tpl = `{{if empty .top.NoSuchThing}}1{{else}}0{{end}}`
	if err := runtv(tpl, "1", dict); err != nil {
		t.Error(err)
	}
	tpl = `{{if empty .bottom.NoSuchThing}}1{{else}}0{{end}}`
	if err := runtv(tpl, "1", dict); err != nil {
		t.Error(err)
	}
}
