package stub

import (
	"embed"
	"io"
	"text/template"

	"github.com/pkg/errors"
)

//go:embed files/*
var stub embed.FS

func parseTemplate(filename string, wr io.Writer, data any) error {
	t, err := template.New(filename).ParseFS(stub, "files/*.stub")
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(t.Execute(wr, data))
}
