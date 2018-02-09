
# Example http parameters oneof

## How to use

```bash
$ go run main.go

$ curl -v http://localhost:4000
# hello world

$ curl -v http://localhost:4000/?prefix=start-
# start-hello world

$ curl -v http://localhost:4000/?suffix=-end
# hello world-end

$ curl -v http://localhost:4000/?prefix=start-&suffix=-end
# {"id":"63d7b11a-26f4-406e-4d35-00f0df466c9c","status":422,"message":"Unprocessable Entity","causes":[{"resource":"query","field":"prefix","code":"conflict"},{"resource":"query","field":"suffix","code":"conflict"}]}
```
