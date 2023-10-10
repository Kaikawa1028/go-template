package errors

import (
	"html/template"
	"strings"
)

func BuildTemplate(errorMessage string, params map[string]interface{}) (string, error) {
	writer := new(strings.Builder)
	t, err := template.New("").Parse(errorMessage)
	if err != nil {
		return "", Wrap(err)
	}
	if err = t.Execute(writer, params); err != nil {
		return "", Wrap(err)
	}

	return writer.String(), nil
}
