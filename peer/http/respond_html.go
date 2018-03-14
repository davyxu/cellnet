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

func compile(pset cellnet.PropertySet) *template.Template {

	var (
		templateDir   string
		delimsLeft    string
		delimsRight   string
		templateExts  []string
		templateFuncs []template.FuncMap
	)

	pset.GetProperty("TemplateDir", &templateDir)
	pset.GetProperty("TemplateExts", &templateExts)
	pset.GetProperty("TemplateFuncs", &templateFuncs)
	pset.GetProperty("DelimsLeft", &delimsLeft)
	pset.GetProperty("DelimsRight", &delimsRight)

	if templateDir == "" {
		templateDir = "."
	}

	if len(templateExts) == 0 {
		templateExts = []string{".tpl", ".html"}
	}

	t := template.New(templateDir)

	t.Delims(delimsLeft, delimsRight)
	// parse an initial template in case we don't have any
	//template.Must(t.Parse("Martini"))

	filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		r, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		ext := getExt(r)

		for _, extension := range templateExts {
			if ext == extension {

				buf, err := ioutil.ReadFile(path)
				if err != nil {
					panic(err)
				}

				name := r[0 : len(r)-len(ext)]
				tmpl := t.New(filepath.ToSlash(name))

				// add our funcmaps
				for _, funcs := range templateFuncs {
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

	log.Debugf("#http.recv(%s) '%s' %s | [%d] HTML %s",
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
