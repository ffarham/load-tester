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
Total Requests:         1000
Successful requests:    170
Failed Requests:        0
Throttled Requests:     830
Refused Connections:    0

[Successful Requests]
Total Time:             1.641219s
Min Response Time:      0.000938s
Average Response Time:  0.009654s
Max Response Time:      0.015749s
Requests Per Second:    103.581532
```