package loadbalancer

import (
	"net/http/httputil"
)

type Node struct {
  byteLoad *int
  requestLoad *int
  maximumByte int
  maximumRequest int
  reverseProxy *httputil.ReverseProxy
}
