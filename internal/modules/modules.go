package modules

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	xpextv1 "github.com/crossplane/crossplane/apis/apiextensions/v1"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	crv1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pkg/errors"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/scale/scheme"
)

const (
	moduleApiGroup     = "modules.krateo.io"
	errParseValidation = "cannot parse validation schema"
)

type PullOpts struct {
	ModuleName string
	Username   string
	Password   string
}

func PullXRD(imageURL string, opts PullOpts) (*xpextv1.CompositeResourceDefinition, error) {
	auth := authn.Anonymous
	if len(opts.Username) > 0 {
		auth = &authn.Basic{Username: opts.Username, Password: opts.Password}
	} else if len(opts.Password) > 0 {
		auth = &authn.Bearer{Token: opts.Password}
	}

	img, err := pullImage(imageURL, auth)
	if err != nil {
		return nil, err
	}

	dat, err := extractMultiYAML(img)
	if err != nil {
		return nil, err
	}

	return findXRD(dat, func(xrd *xpextv1.CompositeResourceDefinition) bool {
		ok := xrd.Spec.Group == moduleApiGroup
		ok = ok && xrd.Spec.Names.Singular == opts.ModuleName
		return ok
	}), nil
}

type acceptFunc func(xrd *xpextv1.CompositeResourceDefinition) bool

func findXRD(dat []byte, accept acceptFunc) (xrd *xpextv1.CompositeResourceDefinition) {
	sch := runtime.NewScheme()
	_ = scheme.AddToScheme(sch)
	_ = xpextv1.AddToScheme(sch)

	manifests := string(dat)
	decode := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode
	for _, spec := range strings.Split(manifests, "---") {
		if len(spec) == 0 {
			continue
		}

		obj, gvk, err := decode([]byte(spec), nil, nil)
		if err != nil {
			continue
		}

		if gvk.Kind != xpextv1.CompositeResourceDefinitionKind {
			continue
		}

		xrd = obj.(*xpextv1.CompositeResourceDefinition)
		if !accept(xrd) {
			continue
		}
	}

	return xrd
}

func extractMultiYAML(img crv1.Image) ([]byte, error) {
	all, err := img.Layers()
	if err != nil {
		return nil, err
	}

	if len(all) > 1 {
		return nil, fmt.Errorf("module image should contain only one layer")
	}

	src, err := all[0].Compressed()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	dst := bufio.NewWriter(&buf)
	err = untar(src, dst)

	return buf.Bytes(), err
}

func XRDSpecs(xrd *xpextv1.CompositeResourceDefinition) (map[string]extv1.JSONSchemaProps, []string, error) {
	vr := xrd.Spec.Versions[0]
	return getProps("spec", vr.Schema)
}

func getProps(field string, v *xpextv1.CompositeResourceValidation) (map[string]extv1.JSONSchemaProps, []string, error) {
	if v == nil {
		return nil, nil, nil
	}

	s := &extv1.JSONSchemaProps{}
	if err := json.Unmarshal(v.OpenAPIV3Schema.Raw, s); err != nil {
		return nil, nil, errors.Wrap(err, errParseValidation)
	}

	spec, ok := s.Properties[field]
	if !ok {
		return nil, nil, nil
	}

	return spec.Properties, spec.Required, nil
}

func pullImage(src string, auth authn.Authenticator) (crv1.Image, error) {
	ref, err := name.ParseReference(src)
	if err != nil {
		return nil, fmt.Errorf("parsing reference %q: %w", src, err)
	}

	if auth == nil {
		auth = authn.Anonymous
	}

	return remote.Image(ref, remote.WithAuth(auth))
}

func untar(in io.Reader, out io.Writer) (err error) {
	gzf, err := gzip.NewReader(in)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzf)

	header, err := tarReader.Next()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	if header.Typeflag != tar.TypeReg {
		return nil
	}

	_, err = io.Copy(out, tarReader)
	return err
}
