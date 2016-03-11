package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/pivotal-golang/lager"
)

type redirectHandler struct {
	logger   lager.Logger
	template *template.Template
	pageMap  PageMap
}

func New(logger lager.Logger, pageMap PageMap) http.Handler {
	template := template.Must(template.New("page").Parse(pageTemplate))
	return redirectHandler{
		logger:   logger,
		template: template,
		pageMap:  pageMap,
	}
}

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

func (h redirectHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.logger.Debug("handling-request", lager.Data{"url": fmt.Sprintf("%#v", req.URL)})

	data, found := h.pageMap.Match(req.URL.Path)
	if found {
		if req.URL.Query().Get("go-get") == "1" {
			err := writeTemplate(h.logger, h.template, w, data)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.Header().Set("Location", "https://"+data.Repo)
			w.WriteHeader(http.StatusFound)
		}
	}
}

func writeTemplate(logger lager.Logger, t *template.Template, w io.Writer, data PageData) error {
	buf := &bytes.Buffer{}

	err := t.Execute(buf, data)
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
