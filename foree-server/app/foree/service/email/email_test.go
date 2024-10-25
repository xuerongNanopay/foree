package foree_email_service

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"
)

func TestHowHTMLTemplateWorks(t *testing.T) {
	t.Run("Demo 1", func(t *testing.T) {
		buf := new(bytes.Buffer)
		template, _ := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
		_ = template.ExecuteTemplate(buf, "T", "<script>alert('you have been pwned')</script>")
		fmt.Println(buf.String())
	})

	t.Run("Test Root Layout", func(t *testing.T) {
		buf := new(bytes.Buffer)
		template, _ := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
		_ = template.ExecuteTemplate(buf, "T", "<script>alert('you have been pwned')</script>")
		fmt.Println(buf.String())
	})
}
