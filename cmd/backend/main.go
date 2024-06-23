package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"
)

type UrlBody struct {
  RequestUrl string `json:"requestUrl"`
  MaximumByte int `json:"maximumByte"`
  MaximumRequest int `json:"maximumRequest"`
}

func main() {
  log.Println("registering to load balacner")

  port := flag.String("port", "0", "port number")
  maxBytes := flag.Int("bytes", 10000, "port number")
  maxRequests := flag.Int("requests", 10, "port number")
  sleep := flag.Int("sleep", 0, "port number")
  flag.Parse()

	body, _ := json.Marshal(UrlBody{
    "localhost:"+*port,
    *maxBytes,
    *maxRequests,
	})

  loadBalancerUrl := "http://127.0.0.1:8080/register"
  r, err := http.NewRequest("POST", loadBalancerUrl, bytes.NewBuffer(body))
  if err != nil {
    log.Fatal("erorr registering to load balancer", err)
    return
  }
  r.Header.Add("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(r)
  if err != nil {
    log.Fatal("erorr registering to load balancer", err)
    return
  }

  if resp.StatusCode != http.StatusOK {
    log.Fatal("erorr registering to load balancer")
    return
  }

  log.Println("Backend | registered to load balancer")
  log.Println("Backend | starting backend at port", *port)

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    log.Println("sleeping")
    time.Sleep(time.Duration(*sleep) * time.Second)
    log.Println("hello")
  })
  log.Fatal(http.ListenAndServe(":"+*port, nil))
}
