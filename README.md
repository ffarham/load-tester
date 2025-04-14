### Load Tester

Tool to simulate high load on an API endpoint.

### Example Usage

```sh
go run cmd/main.go -m=POST -url=http://localhost:8080/ -w=10 -r=100 -t=10 -f=input.json
```

- `m`: **required** flag for specifying the HTTP method (GET, POST, e.t.c.).
- `url`: **required** flag for specifying the URL to load test.
- `w`: flag for specifying the number of concurrent workers. Default is 1.
- `r`: flag for specifying the total number of requests to send. Default is 1.
- `t`: flag for specifying the HTTP request timeout in seconds. Default is 3.
- `f`: flag for specifying the filename to load the json payload from. Payload is empty for post reqests when this flag is not provided, as well as for other requests. Default is "".


### Example Output
```
        ---Load Test Results---
Total Requests:         100
Successful requests:    17
Throttled Requests:     83
Other Failed Requests:  0
Request Per Second:     93.577747

Total Time:             0.181667s
Min Response Time:      0.000359s
Average Response Time:  0.001817s
Max Response Time:      0.011258s
```