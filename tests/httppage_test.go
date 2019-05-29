package tests

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	httppeer "github.com/davyxu/cellnet/peer/http"
	"github.com/davyxu/cellnet/proc"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const pageAddress = "127.0.0.1:10087"

func TestPrintPage(t *testing.T) {

	p := peer.NewGenericPeer("http.Acceptor", "httpserver", pageAddress, nil)

	proc.BindProcessorHandler(p, "http", func(raw cellnet.Event) {

		switch {
		case raw.Session().(httppeer.RequestMatcher).Match("GET", "/"):

			raw.Session().Send(&httppeer.HTMLRespond{
				StatusCode:    http.StatusOK,
				PageTemplate:  "index",
				TemplateModel: "world",
			})
		}

	})

	p.Start()

	validPage(t, fmt.Sprintf("http://%s", pageAddress), "<h1>Hello world</h1>")

	p.Stop()

}

func validPage(t *testing.T, url, expectAck string) {
	c := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := c.Get(url)
	if err != nil {
		t.Log("http req failed", err)
		t.FailNow()
	}

	defer resp.Body.Close()
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log("http response failed", err)
		t.FailNow()
	}

	body := string(bodyData)

	if body != expectAck {
		t.Log("unexpect result", err, body)
		t.FailNow()
	}
}
