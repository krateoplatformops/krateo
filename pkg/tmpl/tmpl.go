package tmpl

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"path"

	"github.com/ghodss/yaml"
	"github.com/krateoplatformops/krateo/pkg/clients/helmclient"
)

//go:embed assets/*
var assetsFS embed.FS

func Execute(name string, data interface{}) ([]byte, error) {
	fn := fmt.Sprintf("assets/%s", name)
	t, err := template.New(path.Base(fn)).Funcs(FuncMap()).ParseFS(assetsFS, fn)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Chart(name string, data interface{}) (*helmclient.ReleaseParameters, error) {
	src, err := Execute(name, data)
	if err != nil {
		return nil, err
	}

	var res helmclient.ReleaseParameters
	if err := yaml.Unmarshal(src, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
