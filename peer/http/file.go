package http

import (
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

func (self *httpAcceptor) SetFileServe(dir string, root string) {

	self.httpDir = dir
	self.httpRoot = root
}

func (self *httpAcceptor) GetDir() http.Dir {

	if filepath.IsAbs(self.httpDir) {
		return http.Dir(self.httpDir)
	} else {
		return http.Dir(filepath.Join(self.httpRoot, self.httpDir))
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

	dir := self.GetDir()

	if dir == "" {
		return nil, errNotFound, false
	}

	err, handled = self.ServeFile(res, req, dir)

	return
}
