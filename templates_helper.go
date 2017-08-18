package authproxy

import (
	"bytes"
	"text/template"
)

// Don't edit the following line. This allows go generate ./... to work.
//go:generate go-bindata -pkg authproxy -o templates.go public/

// RenderTemplate is a generic text/template render wrapper.
func RenderTemplate(tpl string, data interface{}) (b []byte, err error) {
	// Templates don't actually need names.
	t, err := template.New("").Parse(tpl)
	if err != nil {
		return b, err
	}

	buf := &bytes.Buffer{}

	err = t.Execute(buf, data)
	if err != nil {
		return b, nil
	}

	return buf.Bytes(), nil
}
