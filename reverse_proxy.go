package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// NewPathDirectorReverseProxy returns a new ReverseProxy that directs
// network traffic based on the provided mapping.
func NewPathDirectorReverseProxy(mapping map[string]*url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		log.Println(req.URL.Path)
		log.Println(req.URL.Scheme)
		log.Println(req.URL.Host)
		log.Println("yes I printed ********")
		path := strings.TrimSuffix(req.URL.Path, "/")
		if target, ok := mapping[path]; ok {
			log.Println("Path matched!!!")
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
			targetQuery := target.RawQuery
			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}
			if _, ok := req.Header["User-Agent"]; !ok {
				// explicitly disable User-Agent so it's not set to default value
				req.Header.Set("User-Agent", "")
			}
		} else {
			log.Println("path didn't match")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

func main() {
	const (
		defaultPort      = "9090"
		defaultPortUsage = "default server port, ':9090'"
	)
	port := flag.String("port", defaultPort, defaultPortUsage)
	flag.Parse()
	mapping := make(map[string]*url.URL)
	mapping["/aaa"] = &url.URL{
		Scheme: "http",
		Host:   "localhost:9091",
	}
	proxy := NewPathDirectorReverseProxy(mapping)
	log.Fatal(http.ListenAndServe(":"+*port, proxy))
}
