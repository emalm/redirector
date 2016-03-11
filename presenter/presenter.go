package presenter

import (
	"bytes"
	"html/template"
	"io"
	"strings"

	"github.com/pivotal-golang/lager"
)

const pageTemplate = `<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
	<meta name="go-import" content="{{.Domain}}/{{.Path}} git https://{{.Repo}}">
	<meta name="go-source" content="{{.Domain}}/{{.Path}} https://{{.Repo}} https://{{.Repo}}/tree/master{/dir} https://{{.Repo}}/blob/master{/dir}/{file}#L{line}">
	<meta http-equiv="refresh" content="0; url=https://godoc.org/{{.Domain}}/{{.Path}}">
</head>
<body>
Nothing to see here; <a href="https://godoc.org/{{.Domain}}/{{.Path}}">move along</a>.
</body>
</html>
`

type PageData struct {
	Path   string
	Repo   string
	Domain string
}

type PageMap map[string]PageData

func (m PageMap) Match(path string) (PageData, bool) {
	for prefix, data := range m {
		if strings.HasPrefix(path, prefix) {
			return data, true
		}
	}
	return PageData{}, false
}

type PagePresenter struct {
	template *template.Template
}

func NewPagePresenter() *PagePresenter {
	return &PagePresenter{
		template: template.Must(template.New("page").Parse(pageTemplate)),
	}
}

func (p *PagePresenter) WritePage(logger lager.Logger, w io.Writer, data PageData) error {
	buf := &bytes.Buffer{}

	err := p.template.Execute(buf, data)
	if err != nil {
		logger.Error("failed-to-render-template", err)
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		logger.Error("failed-to-write-rendered-template", err)
		return err
	}

	return nil
}
