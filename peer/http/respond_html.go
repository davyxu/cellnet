package http

import (
	"bytes"
	"github.com/davyxu/cellnet"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func getExt(s string) string {
	if strings.Index(s, ".") == -1 {
		return ""
	}
	return "." + strings.Join(strings.Split(s, ".")[1:], ".")
}

func (self *httpAcceptor) SetTemplateDir(dir string) {

	self.templateDir = dir
}

func (self *httpAcceptor) SetTemplateDelims(delimsLeft, delimsRight string) {
	self.delimsLeft = delimsLeft
	self.delimsRight = delimsRight
}

func (self *httpAcceptor) SetTemplateExtensions(exts []string) {
	self.templateExts = exts
}

func (self *httpAcceptor) SetTemplateFunc(f []template.FuncMap) {
	self.templateFuncs = f
}

func (self *httpAcceptor) Compile() *template.Template {

	if self.templateDir == "" {
		self.templateDir = "."
	}

	if len(self.templateExts) == 0 {
		self.templateExts = []string{".tpl", ".html"}
	}

	t := template.New(self.templateDir)

	t.Delims(self.delimsLeft, self.delimsRight)
	// parse an initial template in case we don't have any
	//template.Must(t.Parse("Martini"))

	filepath.Walk(self.templateDir, func(path string, info os.FileInfo, err error) error {
		r, err := filepath.Rel(self.templateDir, path)
		if err != nil {
			return err
		}

		ext := getExt(r)

		for _, extension := range self.templateExts {
			if ext == extension {

				buf, err := ioutil.ReadFile(path)
				if err != nil {
					panic(err)
				}

				name := r[0 : len(r)-len(ext)]
				tmpl := t.New(filepath.ToSlash(name))

				// add our funcmaps
				for _, funcs := range self.templateFuncs {
					tmpl.Funcs(funcs)
				}

				// Bomb out if parse fails. We don't want any silent server starts.
				template.Must(tmpl.Parse(string(buf)))
				break
			}
		}

		return nil
	})

	return t
}

type HTMLRespond struct {
	StatusCode int

	PageTemplate string

	TemplateModel interface{}
}

func (self *HTMLRespond) WriteRespond(ses *httpSession) error {

	peerInfo := ses.Peer().(cellnet.PeerProperty)

	log.Debugf("#http.send(%s) '%s' %s | [%d] HTML %s",
		peerInfo.Name(),
		ses.req.Method,
		ses.req.URL.Path,
		self.StatusCode,
		self.PageTemplate)

	buf := make([]byte, 64)

	bb := bytes.NewBuffer(buf)
	bb.Reset()

	err := ses.t.ExecuteTemplate(bb, self.PageTemplate, self.TemplateModel)

	if err != nil {
		return err
	}

	// template rendered fine, write out the result
	ses.resp.Header().Set("Content-Type", "text/html")
	ses.resp.WriteHeader(self.StatusCode)
	io.Copy(ses.resp, bb)

	return nil
}

type TextRespond struct {
	StatusCode int
	Text       string
}

func (self *TextRespond) WriteRespond(ses *httpSession) error {

	if log.IsDebugEnabled() {
		peerInfo := ses.Peer().(cellnet.PeerProperty)

		log.Debugf("#http.send(%s) '%s' %s | [%d] HTML '%s'",
			peerInfo.Name(),
			ses.req.Method,
			ses.req.URL.Path,
			self.StatusCode,
			self.Text)
	}

	ses.resp.Header().Set("Content-Type", "text/html;charset=utf-8")
	ses.resp.WriteHeader(self.StatusCode)
	ses.resp.Write([]byte(self.Text))

	return nil
}
