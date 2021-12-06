package tmpl

// FuncMap returns a copy of the basic function map as a map[string]interface{}.
func FuncMap() map[string]interface{} {
	gfm := make(map[string]interface{}, len(genericMap))
	for k, v := range genericMap {
		gfm[k] = v
	}
	return gfm
}

var genericMap = map[string]interface{}{
	// Defaults
	"default": dfault,
	"empty":   empty,

	// Encoding:
	"b64enc": base64encode,
	"b64dec": base64decode,

	// Strings:

}
