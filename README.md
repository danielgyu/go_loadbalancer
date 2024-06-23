PoC load balancer project in Go

# Goals
- load balancer should balance load for multiple nodes
- each node can have different rate limits
- rate limits are measured in two ways
  - BPM (HTTP Body Bytes Per Minute)
  - RPM (Requests Per Minute)
- Rate limits can be hit across any of the options above


# Testing
- Open four shell processes
- For each of the processes, run the following
  - 1) `go run cmd/loadbalancer/main.go`
  - 2) `go run cmd/backend/main.go -port 8001`
  - 3) `go run cmd/backend/main.go -port 8002 -request 1 -sleep 10`
  - 4) `go run cmd/backend/main.go -port 8003`
- Open two shell processes, and run the following in order
  - `curl localhost:8080/`

This would set up 1 load balancer process and 3 backend processes connected to the load balancer. Through first 4 curls, the backend processes would receive the requests in order. However, at the 5th curl, if the first request to the 8002 port server hasn't done sleeping, the request would be relayed to the 8003 port server due to rate limiting of 1.

# Known limitations
- if there're no nodes available, the load balancer goes into a infinite loop searching
- no health checking logic of downstream nodes to remove them from list
- downstream node has no way of un-register from the load balancer
- the implemented load balancing strategy is round-robin, can be enhanced to more sophiscated algorithm
