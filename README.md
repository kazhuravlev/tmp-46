# Sample HTTP request "proxy"

This project demonstrate how to accept new requests, parse and validate it and work with bufferpool, taskpool and http
requests. Also you can find a simple graceful shutdown implementation.

## How to run

```shell
git clone git@github.com:kazhuravlev/tmp-46.git tmp-46
cd tmp-46

# run the server
go run ./cmd/server

# send requests
curl \
  -X POST \
  http://localhost:8888/task \
  -d '{"method":"GET", "url":"http://example.com"}'
```
