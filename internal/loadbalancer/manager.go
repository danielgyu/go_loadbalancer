package loadbalancer

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type UrlBody struct {
  RequestUrl string `json:"requestUrl"`
  MaximumByte int `json:"maximumByte"`
  MaximumRequest int `json:"maximumRequest"`
}

type Manager struct {
  nodes [100]Node  // assume 100 is the maximum capacity for one load balancer
  active int
  current int
  maxByteLimit int
  maxRequestLimit int
}

func (manager *Manager) registerNode(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    log.Println("Manager | error registering | wrong method")
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  if manager.active >= len(manager.nodes) {
    log.Println("Manager | error registering | maximum actives")
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  var urlBody UrlBody
  err := json.NewDecoder(r.Body).Decode(&urlBody)
  if err != nil {
    log.Println("Manager | error registering | wrong body")
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  remoteUrl, err := url.Parse("http://"+urlBody.RequestUrl)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  byteLoad, requestLoad := 0, 0
  manager.nodes[manager.active] = Node{
    &byteLoad,
    &requestLoad,
    urlBody.MaximumByte,
    urlBody.MaximumRequest,
    httputil.NewSingleHostReverseProxy(remoteUrl),
  };
  manager.active++
  log.Println("Manager | added", remoteUrl);
  log.Println("Manager | max bytes:", urlBody.MaximumByte);
  log.Println("Manager | max requests", urlBody.MaximumRequest);
  log.Println("Manager | active", manager.active);
}

func (manager *Manager) balanceLoad(w http.ResponseWriter, r *http.Request) {
  // check current status of manager.nodes[current]
  // possible inifinite loop if no nodes are available
  target := manager.current
  selectedNode := manager.nodes[target]
  for (
    *selectedNode.byteLoad >= selectedNode.maximumByte ||
    *selectedNode.requestLoad >= selectedNode.maximumRequest) {
    target = (target+1) % manager.active
    selectedNode = manager.nodes[target]
  }

  *selectedNode.requestLoad += 1
  *selectedNode.byteLoad += int(r.ContentLength)
  manager.current = (target + 1) % manager.active

  // forward request to selected node
  selectedNode.reverseProxy.ServeHTTP(w, r)

  *selectedNode.requestLoad -= 1
  *selectedNode.byteLoad -= int(r.ContentLength)
}

