package handler

import (
	"fmt"
	"net/http"

	"github.com/ematpl/redirector/presenter"
	"github.com/pivotal-golang/lager"
)

type redirectHandler struct {
	logger    lager.Logger
	presenter *presenter.PagePresenter
	pageMap   presenter.PageMap
}

func New(logger lager.Logger, presenter *presenter.PagePresenter, pageMap presenter.PageMap) http.Handler {
	return redirectHandler{
		logger:    logger,
		presenter: presenter,
		pageMap:   pageMap,
	}
}

func (h redirectHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.logger.Debug("handling-request", lager.Data{"url": fmt.Sprintf("%#v", req.URL)})

	data, found := h.pageMap.Match(req.URL.Path)
	if found {
		if req.URL.Query().Get("go-get") == "1" {
			err := h.presenter.WritePage(h.logger, w, data)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.Header().Set("Location", "https://"+data.Repo)
			w.WriteHeader(http.StatusFound)
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
