package kitchen

import (
	"net/http"
	"path"
	"strings"
)

type static struct {
	next   http.Handler
	dir    http.FileSystem
	prefix string
	index  string
}

func (s *static) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		s.next.ServeHTTP(rw, r)
		return
	}
	file := r.URL.Path
	// if we have a prefix, filter requests by stripping the prefix
	if s.prefix != "" {
		if !strings.HasPrefix(file, s.prefix) {
			s.next.ServeHTTP(rw, r)
			return
		}
		file = file[len(s.prefix):]
		if file != "" && file[0] != '/' {
			s.next.ServeHTTP(rw, r)
			return
		}
	}
	f, err := s.dir.Open(file)
	if err != nil {
		// discard the error?
		s.next.ServeHTTP(rw, r)
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		s.next.ServeHTTP(rw, r)
		return
	}
	// try to serve index file
	if fi.IsDir() {
		// redirect if missing trailing slash
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(rw, r, r.URL.Path+"/", http.StatusFound)
			return
		}
		file = path.Join(file, s.index)
		f, err = s.dir.Open(file)
		if err != nil {
			s.next.ServeHTTP(rw, r)
			return
		}
		defer f.Close()
		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			s.next.ServeHTTP(rw, r)
			return
		}
	}
	http.ServeContent(rw, r, file, fi.ModTime(), f)
}

// StaticMiddlewareGenerator create a static middleware function for serve the file from a folder
// This part is base on negroni static handler
func StaticMiddlewareGenerator(dir http.FileSystem, prefix, index string) Middleware {
	fn := func(next http.Handler) http.Handler {
		return &static{next, dir, prefix, index}
	}
	return fn
}
