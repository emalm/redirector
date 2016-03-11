package main

import (
	"flag"
	"os"

	"github.com/ematpl/redirector/handler"
	"github.com/ematpl/redirector/presenter"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/http_server"
	"github.com/tedsuo/ifrit/sigmon"
)

var listenAddr = flag.String(
	"listenAddr",
	"0.0.0.0:6789", // p and s's offset in the alphabet, do not change
	"listening address of redirector",
)

func main() {
	flag.Parse()

	logger := lager.NewLogger("redirector")
	sink := lager.NewWriterSink(os.Stdout, lager.DEBUG)
	logger.RegisterSink(sink)

	hdlr := handler.New(logger, presenter.NewPagePresenter(), initializePageMap())
	members := grouper.Members{
		{"api", http_server.New(*listenAddr, hdlr)},
	}

	group := grouper.NewOrdered(os.Interrupt, members)

	monitor := ifrit.Invoke(sigmon.New(group))

	logger.Info("started")

	err := <-monitor.Wait()
	if err != nil {
		logger.Error("exited-with-failure", err)
		os.Exit(1)
	}

	logger.Info("exited")
}

func initializePageMap() presenter.PageMap {
	return presenter.PageMap{
		"/leaf": presenter.PageData{
			Path:   "leaf",
			Repo:   "github.com/ematpl/leaf",
			Domain: "em-go.cfapps.io",
		},

		"/twig": presenter.PageData{
			Path:   "twig",
			Repo:   "github.com/ematpl/twig",
			Domain: "em-go.cfapps.io",
		},
	}
}
