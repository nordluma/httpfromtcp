# Http From TCP

An implementation of the most common features of HTTP/1.1 with golang.

# Running

```bash
go run cmd/httpserver/main.go
```

Available endpoints:

- `/yourproblem`: allways returns a 400 error
- `/myproblem`: returns a 500 error.
- GET `/httpbin/stream/{number_of_responses}`: returns defined number for
  responses from `https://httpbin.org`.

Running tests:

```bash
go test ./...
```
