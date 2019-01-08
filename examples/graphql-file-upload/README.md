
# Example graphql upload server

## How to use

```bash
$ go run main.go

$ curl -v -d "query=query{file}" http://localhost:4000/graphql
{"data":{"file":null},"errors":[{"message":"file has not been uploaded yet","locations":[]}]}

$ curl -v -XPOST -F "file=@./image.png" -F "query=mutation{upload}" http://localhost:4000/graphql
{"data":{"upload":"File uploaded: image/png type, 91653 bytes"}}

$ curl -v -d "query=query{file}" http://localhost:4000/graphql
{"data":{"file":"File exists: image/png type, 91653 bytes"}}
```
