package loadbalancer

import (
	"log"
	"net/http"
)

func RunLoadBalancer() {
  var manager Manager = Manager{[100]Node{}, 0, 0, 10000, 100};

  http.HandleFunc("/register", manager.registerNode)
  http.HandleFunc("/", manager.balanceLoad)

  port := ":8080"
  log.Println("starting load balancer at port", port)
  log.Fatal(http.ListenAndServe(":8080", nil))
}
