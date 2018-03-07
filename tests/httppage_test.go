package tests

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	httppeer "github.com/davyxu/cellnet/peer/http"
	"github.com/davyxu/cellnet/proc"
	"net/http"
	"testing"
)

func TestPrintPage(t *testing.T) {
	p := peer.NewPeer("http.Acceptor")
	pset := p.(cellnet.PropertySet)
	pset.SetProperty("Name", "httpserver")
	pset.SetProperty("Address", "127.0.0.1:8081")
	proc.BindProcessor(p, "http", func(raw cellnet.Event) {

		switch {
		case raw.Session().(httppeer.RequestMatcher).Match("GET", "/"):

			raw.Session().Send(httppeer.HTML{
				Code:          http.StatusOK,
				PageTemplate:  "index",
				TemplateModel: "world",
			})
		}

	})

	p.Start()

	validPage(t, "http://127.0.0.1:8081", "<h1>Hello world</h1>")

}

func init() {
	//cellnet.RegisterHttpMeta(&cellnet.HttpMeta{
	//	URL:    "/",
	//	Method: "GET",
	//})
}
