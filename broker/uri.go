package broker

import (
        "strings"
        "net/url"
        "errors"
)

type URI struct {
        scheme string
        host   string
        port   string
}

var defaultUri = URI{
        scheme: "tcp",
        host:   "localhost",
        port:   "9001",
}

var errURIWhitespace = errors.New("URI contain whitespace")
var errURIScheme = errors.New("URI Scheme should be tcp")
var errURIInvalidPort = errors.New("URI port invalid")

func parseUri(uri string) (URI, error) {
        du := defaultUri 
        if strings.Contains(uri, " ") == true {
                return du, errURIWhitespace
        }
        u, err := url.Parse(uri)
        if err != nil {
                return du, err
        }
        if u.Scheme != "tcp" {
                return du, errURIScheme
        }
        h := u.Hostname()
        p := u.Port()

        if h != "" {
                du.host = h
        }
        if p != "" {
                if p != du.port {
                        return du, errURIInvalidPort
                }
        }
        return du, nil
}
