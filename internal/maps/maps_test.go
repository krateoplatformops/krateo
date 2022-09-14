package maps

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMap(t *testing.T) {

	m := NewMap("name", "Mat", "age", 29, "bool", true)

	assert.Equal(t, "Mat", m["name"])
	assert.Equal(t, 29, m["age"])
	assert.Equal(t, true, m["bool"])

	assert.Panics(t, func() {
		NewMap(1, "Mat", "age", 29, "bool", true)
	}, "Non string key should panic")
	assert.Panics(t, func() {
		NewMap("name", "Mat", "age", 29, "bool")
	}, "Wrong number of arguments should panic")

}

func TestM(t *testing.T) {

	m := M("name", "Mat", "age", 29, "bool", true)

	assert.Equal(t, "Mat", m["name"])
	assert.Equal(t, 29, m["age"])
	assert.Equal(t, true, m["bool"])

	assert.Panics(t, func() {
		M(1, "Mat", "age", 29, "bool", true)
	}, "Non string key should panic")
	assert.Panics(t, func() {
		M("name", "Mat", "age", 29, "bool")
	}, "Wrong number of arguments should panic")

}

func TestCopy(t *testing.T) {

	d1 := make(Map)
	d1["name"] = "Tyler"
	d1["location"] = "UT"

	d2 := d1.Copy()
	d2["name"] = "Mat"

	assert.Equal(t, d1["name"], "Tyler")
	assert.Equal(t, d2["name"], "Mat")

}

func TestMerge(t *testing.T) {

	d := make(Map)
	d["name"] = "Mat"

	d1 := make(Map)
	d1["name"] = "Tyler"
	d1["location"] = "UT"

	merged := d.Merge(d1)

	assert.Equal(t, merged["name"], d1["name"])
	assert.Equal(t, merged["location"], d1["location"])
	assert.Nil(t, d["location"])

}

func TestMSI(t *testing.T) {

	m := NewMap("name", "Mat", "age", 29, "bool", true)

	var msi map[string]interface{} = m.MSI()

	assert.Equal(t, "Mat", msi["name"])
	assert.Equal(t, 29, msi["age"])
	assert.Equal(t, true, msi["bool"])

}

func TestMergeHere(t *testing.T) {

	d := make(Map)
	d["name"] = "Mat"

	d1 := make(Map)
	d1["name"] = "Tyler"
	d1["location"] = "UT"

	merged := d.MergeHere(d1)

	assert.Equal(t, d, merged, "With MergeHere, it should return the first modified map")
	assert.Equal(t, merged["name"], d1["name"])
	assert.Equal(t, merged["location"], d1["location"])
	assert.Equal(t, merged["location"], d["location"])

}

func TestTransform(t *testing.T) {

	d1 := make(Map)
	d1["name"] = "Tyler"
	d1["location"] = "UT"
	d1["language"] = "English"

	d2 := d1.Transform(func(k string, v interface{}) (string, interface{}) {
		return strings.ToLower(k), strings.ToLower(v.(string))
	})

	assert.NotEqual(t, d1, d2, "Should be a new map")

	assert.Equal(t, "tyler", d2["name"])
	assert.Equal(t, "ut", d2["location"])
	assert.Equal(t, "english", d2["language"])

}

func TestTransformKeys(t *testing.T) {

	d1 := make(Map)
	d1["name"] = "Tyler"
	d1["location"] = "UT"
	d1["language"] = "English"
	d1["untouched"] = true

	mapping := map[string]string{
		"name":     "Username",
		"location": "Place",
		"language": "Speaks",
	}

	d2 := d1.TransformKeys(mapping)

	assert.Equal(t, "Tyler", d2["Username"])
	assert.Equal(t, "UT", d2["Place"])
	assert.Equal(t, "English", d2["Speaks"])
	assert.Equal(t, true, d2["untouched"])

}

func TestExclude(t *testing.T) {

	d := make(Map)
	d["name"] = "Mat"
	d["age"] = 29
	d["secret"] = "ABC"

	excluded := d.Exclude([]string{"secret"})

	assert.Equal(t, d["name"], excluded["name"])
	assert.Equal(t, d["age"], excluded["age"])
	assert.False(t, excluded.Has("secret"), "secret should be excluded")

}

func TestHas(t *testing.T) {

	d := make(Map)
	d["name"] = "Mat"

	assert.True(t, d.Has("name"))
	assert.False(t, d.Has("nope"))

}

func TestHas_WithDeepNesting(t *testing.T) {

	var l Map = Map{"request": Map{"url": "http://www.stretchr.com/"}}

	assert.True(t, l.Has("request.url"))
	assert.False(t, l.Has("request.method"))
	assert.False(t, l.Has("nothing"))

}

func TestGet(t *testing.T) {

	var l Map = Map{"request": Map{"url": "http://www.stretchr.com/"}}

	assert.Equal(t, "http://www.stretchr.com/", l.Get("request.url"))

	// test some fail cases
	assert.Nil(t, l.Get("something.that.doesnt.exist"))
	assert.Nil(t, l.Get("request.url.somethingelse"))
	assert.Nil(t, l.Get("request.somethingelse"))

}

func TestGetOrDefault(t *testing.T) {

	var defaultValue string = "Default"
	var l Map = Map{"request": Map{"url": "http://www.stretchr.com/"}}

	assert.Equal(t, defaultValue, l.GetOrDefault("request.nope", defaultValue))
	assert.Equal(t, "http://www.stretchr.com/", l.GetOrDefault("request.url", defaultValue))

}

func TestGetString(t *testing.T) {

	var l Map = Map{"request": Map{"url": "http://www.stretchr.com/"}}

	assert.Equal(t, l.GetString("request.url"), "http://www.stretchr.com/")

}

func TestGetStringOrDefault(t *testing.T) {

	var l Map = Map{"request": Map{"url": "http://www.stretchr.com/"}}

	assert.Equal(t, l.GetStringOrDefault("request.url", "default"), "http://www.stretchr.com/")
	assert.Equal(t, l.GetStringOrDefault("request.nope", "default"), "default")

}

func TestGetStringOrEmpty(t *testing.T) {

	var l Map = Map{"request": Map{"url": "http://www.stretchr.com/"}}

	assert.Equal(t, l.GetStringOrEmpty("request.url"), "http://www.stretchr.com/")
	assert.Equal(t, l.GetStringOrEmpty("request.nope"), "")

}

func TestGet_WithNativeMap(t *testing.T) {

	var l Map = Map{"request": map[string]interface{}{"url": "http://www.stretchr.com/"}}

	assert.Equal(t, "http://www.stretchr.com/", l.Get("request.url"))

}

func TestSet_Simple(t *testing.T) {
	// https://github.com/stretchr/stew/issues/2

	var m Map = make(Map)
	assert.Equal(t, m, m.Set("name", "Tyler"))

	assert.Equal(t, "Tyler", m["name"])

}

func TestSet_Deep(t *testing.T) {
	// https://github.com/stretchr/stew/issues/2

	var m Map = make(Map)
	assert.Equal(t, m, m.Set("personal.info.name.first", "Tyler"))

	assert.Equal(t, "Tyler", m.Get("personal.info.name.first"))

	nameObj := m.Get("personal.info.name")
	if assert.NotNil(t, nameObj) {
		assert.Equal(t, "Tyler", nameObj.(Map)["first"])
	}

}

func Test_GetMap(t *testing.T) {

	var parent Map = make(Map)
	var child Map = make(Map)
	child.Set("name", "child")

	parent.Set("child", child)

	var gottenChild Map = parent.GetMap("child")
	assert.Equal(t, "child", gottenChild.Get("name"))

}

func TestMapJSON(t *testing.T) {

	m := make(Map)

	m.Set("name", "tyler")

	json, err := m.JSON()

	if assert.NoError(t, err) {
		assert.Equal(t, json, "{\"name\":\"tyler\"}")
	}

}

func TestMapNewMapFromJSON(t *testing.T) {

	m, err := NewMapFromJSON("{\"name\":\"tyler\"}")

	if assert.NotNil(t, m) && assert.NoError(t, err) {
		assert.Equal(t, m.Get("name").(string), "tyler")
	}

}

func TestMapBase64(t *testing.T) {

	m := make(Map)

	m.Set("name", "tyler")

	b64, err := m.Base64()

	if assert.NoError(t, err) {
		assert.Equal(t, b64, "eyJuYW1lIjoidHlsZXIifQ==")
	}

}

func TestNewMapFromBase64String(t *testing.T) {

	m, err := NewMapFromBase64String("eyJuYW1lIjoidHlsZXIifQ==")

	if assert.NotNil(t, m) && assert.NoError(t, err) {
		assert.Equal(t, m.Get("name").(string), "tyler")
	}

}

func TestMapHash(t *testing.T) {

	m := make(Map)

	m.Set("name", "tyler")

	hash, err := m.Hash()

	if assert.NoError(t, err) {
		assert.Equal(t, hash, "4100f62944bafb39f3cd36a08fe7094482b69207")
	}

}

func TestNewMapFromURLQuery(t *testing.T) {

	m, err := NewMapFromURLQuery("name=tyler&state=UT")
	if assert.NoError(t, err) && assert.NotNil(t, m) {
		assert.Equal(t, "tyler", m.Get("name"))
		assert.Equal(t, "UT", m.Get("state"))
	}

}

func TestURLQuery(t *testing.T) {

	m := make(Map).
		Set("name", "tyler").
		Set("state", "UT")

	encoded, err := m.URLQuery()
	if assert.NoError(t, err) {
		assert.Equal(t, encoded, "name=tyler&state=UT")
	}

}

func TestURLValues(t *testing.T) {

	m := make(Map).
		Set("name", "tyler").
		Set("state", "UT")

	values := m.URLValues()
	if assert.NotNil(t, values) {
		encoded := values.Encode()
		assert.Equal(t, encoded, "name=tyler&state=UT")
	}

}
