# httperror

httperror is simple struct for RFC7807 error response.

## Usage

### standalone

```go
package main

import (
	"fmt"
	"github.com/snowmerak/httperror"
	"net/http"
)

func main() {
	response := httperror.New("invalid username", http.StatusBadRequest, "http://localhost:3000/error/400").WithInstance("/api/auth/login").WithDetail("the username is not valid format").WithExtensionMember("username", "한글 미지원인데요?")

	data, err := response.ToJSON(nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	received, err := httperror.FromJSON(data, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", received)
	fmt.Printf("%v\n", received.ExtensionMembers["username"])
}
```

```shell
{"type":"http://localhost:3000/error/400","title":"invalid username","status":400,"detail":"the username is not valid format","instance":"/api/auth/login","username":"한글 미지원인데요?"}
invalid username on /api/auth/login reference http://localhost:3000/error/400
한글 미지원인데요?
```

### with http response writer

```go
package main

import (
	"github.com/snowmerak/httperror"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		response := httperror.New("just error", http.StatusNotImplemented, "http://localhost:8080/error/501").WithInstance(request.URL.Path).WithDetail("this service is not implemented")

		if _, err := response.WriteToHttpResponseWriter(writer, nil); err != nil {
			log.Printf("[ERR] %v", err)
			return
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Printf("[ERR] %v", err)
	}
}
```

```
HTTP/1.1 501 Not Implemented
Content-Type: application/problem+json
Date: Tue, 21 May 2024 09:14:54 GMT
Content-Length: 134

{"type":"http://localhost:8080/error/501","title":"just error","status":501,"detail":"this service is not implemented","instance":"/"}
```
