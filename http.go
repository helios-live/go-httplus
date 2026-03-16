package httplus

import (
	"encoding/base64"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/rotisserie/eris"
)

// GetHostPort returns host and port from a proxy HTTP request
func GetHostPort(r *http.Request) (host string, port int) {
	host = r.Host
	host, sport, err := net.SplitHostPort(r.Host)
	if err == nil {
		port, _ = strconv.Atoi(sport)
	} else {
		if r.Host != "" {
			host = r.Host
		}
		if r.URL.Scheme == "" || r.URL.Scheme == "http" {
			port = 80
		}
		if r.URL.Scheme == "https" {
			port = 443
		}
	}
	return host, port
}

// CopyHeader copies all the headers from src to dst
func CopyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// GetAuth decodes the Proxy-Authorization header
// @param *http.Request
// @returns []string
func GetAuth(r *http.Request) (ret []string, err error) {
	s := r.Header.Get("Proxy-Authorization")
	if s == "" {
		return ret, eris.New("getauth empty header")
	}
	ss := strings.Split(s, " ")
	if ss[0] != "Basic" {
		return ret, eris.New("getauth not basic")
	}
	b, err := base64.StdEncoding.DecodeString(ss[1])
	if err != nil {
		return ret, eris.Wrap(err, "getauth base64decode failed")
	}
	return strings.Split(string(b), ":"), nil
}
