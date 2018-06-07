package tests

import (
	"encoding/json"
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	_ "github.com/davyxu/cellnet/codec/httpform"
	_ "github.com/davyxu/cellnet/codec/httpjson"
	"github.com/davyxu/cellnet/peer"
	httppeer "github.com/davyxu/cellnet/peer/http"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/http"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"
)

const httpTestAddr = "127.0.0.1:8081"

func TestHttp(t *testing.T) {

	p := peer.NewGenericPeer("http.Acceptor", "httpserver", httpTestAddr, nil)

	proc.BindProcessorHandler(p, "http", func(raw cellnet.Event) {

		switch raw.Message().(type) {
		case *HttpEchoREQ:

			raw.Session().Send(&httppeer.MessageRespond{
				StatusCode: http.StatusOK,
				Msg: &HttpEchoACK{
					Status: 0,
					Token:  "ok",
				},
				CodecName: "httpjson",
			})

		}

	})

	p.Start()

	requestThenValid(t, "GET", "/hello_form", &HttpEchoREQ{
		UserName: "kitty_form",
	}, &HttpEchoACK{
		Token: "ok",
	})

	requestThenValid(t, "POST", "/hello_json", &HttpEchoREQ{
		UserName: "kitty_json",
	}, &HttpEchoACK{
		Token: "ok",
	})

	p.Stop()

	//validPage(t, "http://127.0.0.1:8081", "")
}

func requestThenValid(t *testing.T, method, path string, req, expectACK interface{}) {

	p := peer.NewGenericPeer("http.Connector", "httpclient", httpTestAddr, nil).(cellnet.HTTPConnector)

	ack, err := p.Request(method, path, req)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(ack, expectACK) {
		t.Log("unexpect token result", err)
		t.FailNow()
	}

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

func postForm(t *testing.T) {
	resp, err := http.PostForm("http://127.0.0.1:8081/hello",
		url.Values{"UserName": {"kitty"}})

	if err != nil {
		t.Log("http req failed", err)
		t.FailNow()
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log("http response failed", err)
		t.FailNow()
	}

	var ack HttpEchoACK
	if err := json.Unmarshal(body, &ack); err != nil {
		t.Log("json unmarshal failed", err)
		t.FailNow()
	}

	if ack.Token != "ok" {
		t.Log("unexpect token result", err)
		t.FailNow()
	}

}

type HttpEchoREQ struct {
	UserName string
}

type HttpEchoACK struct {
	Token  string
	Status int32
}

func (self *HttpEchoREQ) String() string { return fmt.Sprintf("%+v", *self) }
func (self *HttpEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellnet.RegisterHttpMeta(&cellnet.HttpMeta{
		Path:         "/hello_form",
		Method:       "GET",
		RequestCodec: codec.MustGetCodec("httpform"),
		RequestType:  reflect.TypeOf((*HttpEchoREQ)(nil)).Elem(),

		// 请求方约束
		ResponseCodec: codec.MustGetCodec("httpjson"),
		ResponseType:  reflect.TypeOf((*HttpEchoACK)(nil)).Elem(),
	})

	cellnet.RegisterHttpMeta(&cellnet.HttpMeta{
		Path:         "/hello_json",
		Method:       "POST",
		RequestCodec: codec.MustGetCodec("httpjson"),
		RequestType:  reflect.TypeOf((*HttpEchoREQ)(nil)).Elem(),

		ResponseCodec: codec.MustGetCodec("httpjson"),
		ResponseType:  reflect.TypeOf((*HttpEchoACK)(nil)).Elem(),
	})

}
