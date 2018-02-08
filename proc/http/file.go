package http

import (
	"github.com/davyxu/cellnet"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type StaticFileProc struct {
	dir http.Dir
}

func (self *StaticFileProc) InitDir(set cellnet.PropertySet) {

	var (
		httpDir  string
		httpRoot string
	)
	set.GetProperty("HttpDir", &httpDir)
	set.GetProperty("HttpRoot", &httpRoot)

	if filepath.IsAbs(httpDir) {
		self.dir = http.Dir(httpDir)
	} else {
		self.dir = http.Dir(filepath.Join(httpRoot, httpDir))
	}

	workDir, _ := os.Getwd()
	log.Debugf("Http serve file: %s (%s)", self.dir, workDir)

}

func (self *StaticFileProc) ServeFile(res http.ResponseWriter, req *http.Request) error {
	if req.Method != "GET" && req.Method != "HEAD" {
		return nil
	}

	file := req.URL.Path

	f, err := self.dir.Open(file)
	if err != nil {

		if err != nil {
			return errNotFound
		}
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return errNotFound
	}

	// try to serve index file
	if fi.IsDir() {
		// redirect if missing trailing slash
		if !strings.HasSuffix(req.URL.Path, "/") {
			dest := url.URL{
				Path:     req.URL.Path + "/",
				RawQuery: req.URL.RawQuery,
				Fragment: req.URL.Fragment,
			}
			http.Redirect(res, req, dest.String(), http.StatusFound)
			return nil
		}

		file = path.Join(file, "index.html")
		f, err = self.dir.Open(file)
		if err != nil {
			return errNotFound
		}
		defer f.Close()

		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			return errNotFound
		}
	}

	log.Debugln("#file ", file)

	http.ServeContent(res, req, file, fi.ModTime(), f)

	return nil
}

func (self *StaticFileProc) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	httpContext := ses.(HttpContext)
	req := httpContext.Request()
	resp := httpContext.Response()

	if self.dir == "" {
		self.InitDir(ses.Peer().(cellnet.PropertySet))
	}

	return nil, self.ServeFile(resp, req)
}

func (StaticFileProc) OnSendMessage(ses cellnet.Session, raw interface{}) error {

	return nil
}
