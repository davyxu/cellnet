package http

import (
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

type StaticFile struct {
	dir http.Dir
}

func (self *StaticFile) ServeFile(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" && req.Method != "HEAD" {
		return
	}

	file := req.URL.Path

	f, err := self.dir.Open(file)
	if err != nil {

		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return
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
			return
		}

		file = path.Join(file, "index.html")
		f, err = self.dir.Open(file)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		defer f.Close()

		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			return
		}
	}

	log.Debugln("[Static] Serving ", file)

	http.ServeContent(res, req, file, fi.ModTime(), f)
}

func newStaticFile(dirStr string, root string) *StaticFile {

	if !filepath.IsAbs(dirStr) {
		dirStr = filepath.Join(root, dirStr)
	}

	return &StaticFile{
		dir: http.Dir(dirStr),
	}
}
