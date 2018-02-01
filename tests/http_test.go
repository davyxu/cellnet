package tests

import (
	"encoding/json"
	"fmt"
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/httpform"
	_ "github.com/davyxu/cellnet/codec/json"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/http"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestHttp(t *testing.T) {

	p := peer.NewPeer("http.Acceptor")
	//infoSetter := p.(cellnet.PeerInfo)
	//infoSetter.SetAddress("127.0.0.1:8081")
	//eventInitor := p.(cellnet.DuplexEventInitor)
	//eventInitor.SetRaw(func(raw cellnet.EventParam) cellnet.EventResult {
	//
	//	switch ev := raw.(type) {
	//	case *cellnet.RecvMsgEvent:
	//
	//		ev.Send(&HttpEchoACK{
	//			Status: 0,
	//			Token:  "ok",
	//		})
	//	}
	//
	//	return nil
	//
	//}, nil)
	p.Start()

	fmt.Scanln()
}

func sendRequest(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:8081/hello?UserName=kitty")
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
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec:  cellnet.MustGetCodec("httpform"),
		URL:    "/hello",
		Method: "GET",
		Type:   reflect.TypeOf((*HttpEchoREQ)(nil)).Elem(),
	})

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: cellnet.MustGetCodec("json"),
		Type:  reflect.TypeOf((*HttpEchoACK)(nil)).Elem(),
	})
}
