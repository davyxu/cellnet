package http

import (
	"github.com/davyxu/cellnet"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

func (self *httpAcceptor) GetDir(set cellnet.PropertySet) http.Dir {

	var (
		httpDir  string
		httpRoot string
	)
	set.GetProperty("HttpDir", &httpDir)
	set.GetProperty("HttpRoot", &httpRoot)

	if filepath.IsAbs(httpDir) {
		return http.Dir(httpDir)
	} else {
		return http.Dir(filepath.Join(httpRoot, httpDir))
	}

	//workDir, _ := os.Getwd()
	//log.Debugf("Http serve file: %s (%s)", self.dir, workDir)
}

func (self *httpAcceptor) ServeFile(res http.ResponseWriter, req *http.Request, dir http.Dir) (error, bool) {
	if req.Method != "GET" && req.Method != "HEAD" {
		return nil, false
	}

	file := req.URL.Path

	f, err := dir.Open(file)
	if err != nil {

		if err != nil {
			return errNotFound, false
		}
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return errNotFound, false
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
			return nil, false
		}

		file = path.Join(file, "index.html")
		f, err = dir.Open(file)
		if err != nil {
			return errNotFound, false
		}
		defer f.Close()

		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			return errNotFound, false
		}
	}

	http.ServeContent(res, req, file, fi.ModTime(), f)

	return nil, true
}

func (self *httpAcceptor) ServeFileWithDir(res http.ResponseWriter, req *http.Request) (msg interface{}, err error, handled bool) {

	dir := self.GetDir(&self.CorePropertySet)

	if dir == "" {
		log.Warnln("peer's 'HttpDir' 'HttpRoot' property not set")
		return nil, errNotFound, false
	}

	err, handled = self.ServeFile(res, req, dir)

	return
}
