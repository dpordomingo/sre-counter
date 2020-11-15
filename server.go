package main

import (
    "fmt"
    "net/http"
    "os"
)

var REDIS_HOST_PORT = os.Getenv("REDIS_HOST_PORT")

func hello(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "hello\n")
}

func counter(w http.ResponseWriter, req *http.Request) {
    val, err := GetFromRedis(REDIS_HOST_PORT, "hi")
    if err != nil {
	fmt.Fprintf(w, fmt.Errorf("redis returned an error. %s", err).Error())
    }
    fmt.Fprintf(w, val)
}

func put(w http.ResponseWriter, req *http.Request) {
    PutInRedis(REDIS_HOST_PORT, "hi", "aaass")
}

func headers(w http.ResponseWriter, req *http.Request) {
    for name, headers := range req.Header {
        for _, h := range headers {
            fmt.Fprintf(w, "%v: %v\n", name, h)
        }
    }
}

const port = 8090
func main() {
    fmt.Println(REDIS_HOST_PORT)
    fmt.Println("Starting server...")
    http.HandleFunc("/hello", hello)
    http.HandleFunc("/headers", headers)
    http.HandleFunc("/counter", counter)
    http.HandleFunc("/put", put)
    fmt.Printf("Listening on :%d\n", port)
    err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
    if err != nil {
        fmt.Println(fmt.Errorf("Server rised a panic. %s", err))
    }
}
