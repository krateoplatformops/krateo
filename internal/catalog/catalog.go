package catalog

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Machiel/slugify"
)

const (
	catalogURL = "https://raw.githubusercontent.com/krateoplatformops/catalog/main/index.json"
)

type Catalog struct {
	Items []PackageInfo `json:"packages"`
}

type PackageInfo struct {
	Image       string `json:"image"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Cli         bool   `json:"cli"`
	Manifest    string `json:"package"`
}

func Fetch() (*Catalog, error) {
	client := &http.Client{Timeout: 40 * time.Second}
	r, err := client.Get(catalogURL)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	target := &Catalog{}
	err = json.NewDecoder(r.Body).Decode(target)
	if err != nil {
		return nil, err
	}
	return target, nil
}

func FetchManifest(info *PackageInfo) ([]byte, error) {
	client := &http.Client{Timeout: 40 * time.Second}
	r, err := client.Get(info.Manifest)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	res := strings.ReplaceAll(string(data), "VERSION", info.Version)
	return []byte(res), nil
}

type FilterFunc func(PackageInfo) bool

func ForCLI() FilterFunc {
	return func(info PackageInfo) bool {
		return info.Cli == true
	}
}

func FilterBy(criteria FilterFunc) (*Catalog, error) {
	all, err := Fetch()
	if err != nil {
		return nil, err
	}

	if criteria == nil {
		return all, nil
	}

	res := &Catalog{Items: []PackageInfo{}}
	for _, el := range all.Items {
		if criteria(el) {
			el.Name = slugify.Slugify(el.Name)
			res.Items = append(res.Items, el)
		}
	}

	return res, nil
}
