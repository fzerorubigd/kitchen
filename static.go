package kitchen

import (
	"net/http"
	"path"
	"strings"
)

// StaticMiddleware create a static middleware function for serve the file from a folder
// This part is base on negroni static handler
func StaticMiddleware(dir http.FileSystem, prefix, index string) MiddlewareFunc {
	fn := func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" && r.Method != "HEAD" {
				next.ServeHTTP(rw, r)
				return
			}
			file := r.URL.Path
			// if we have a prefix, filter requests by stripping the prefix
			if prefix != "" {
				if !strings.HasPrefix(file, prefix) {
					next.ServeHTTP(rw, r)
					return
				}
				file = file[len(prefix):]
				if file != "" && file[0] != '/' {
					next.ServeHTTP(rw, r)
					return
				}
			}
			f, err := dir.Open(file)
			if err != nil {
				// discard the error?
				next.ServeHTTP(rw, r)
				return
			}
			defer f.Close()
			fi, err := f.Stat()
			if err != nil {
				next.ServeHTTP(rw, r)
				return
			}
			// try to serve index file
			if fi.IsDir() {
				// redirect if missing trailing slash
				if !strings.HasSuffix(r.URL.Path, "/") {
					http.Redirect(rw, r, r.URL.Path+"/", http.StatusFound)
					return
				}
				file = path.Join(file, index)
				f, err = dir.Open(file)
				if err != nil {
					next.ServeHTTP(rw, r)
					return
				}
				defer f.Close()
				fi, err = f.Stat()
				if err != nil || fi.IsDir() {
					next.ServeHTTP(rw, r)
					return
				}
			}
			http.ServeContent(rw, r, file, fi.ModTime(), f)
		}
		return http.HandlerFunc(fn)
	}
	return fn
}
