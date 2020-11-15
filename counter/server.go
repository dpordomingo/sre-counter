package counter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type handler func(w http.ResponseWriter, req *http.Request)

type route struct {
	endpoint string
	handler  handler
}

// Server is the HTTP server with the counter functionality
type Server interface {
	Run() error
}

// NewServer returns the Server
func NewServer(cacheClient CacheClient, port, instanceName string) Server {
	serv := server{
		port:         port,
		cacheClient:  cacheClient,
		instanceName: instanceName,
	}

	serv.routes = []route{
		route{"/", serv.index},
		route{"/counter", serv.counter},
		route{"/reset", serv.reset},

		// TODO: The ones below can be removed later.
		route{"/extra", serv.extra},
		route{"/set", serv.set},
	}

	return &serv
}

const redisCounterKey = "counter"
const redisExtraKey = "extra"

type server struct {
	port         string
	instanceName string
	cacheClient  CacheClient
	routes       []route
}

func (s *server) Run() error {
	fmt.Printf("Listening on :%s\n", s.port)

	for _, h := range s.routes {
		fmt.Printf("- %s\n", h.endpoint)
		http.HandleFunc(h.endpoint, instanceHandler(s.instanceName, h.handler))
	}

	http.HandleFunc("/healthcheck", s.healthcheck)

	return http.ListenAndServe(fmt.Sprintf(":%s", s.port), nil)
}

func instanceHandler(name string, h handler) handler {
	return func(w http.ResponseWriter, req *http.Request) {
		h(w, req)
		fmt.Fprintf(w, "\n... from %s\n", name)
	}
}

func (s *server) counter(w http.ResponseWriter, req *http.Request) {
	counterStr, err := s.cacheClient.Get(redisCounterKey)
	if err != nil {
		fmt.Fprintf(w, "%s\n", fmt.Errorf("error. redis returned an error. %s", err).Error())
		return
	}

	counter, err := strconv.Atoi(counterStr)
	if err != nil {
		fmt.Fprintf(w, "%s\n", fmt.Errorf("error. counter corrupted. %s", err).Error())
		return
	}

	err = s.cacheClient.Set(redisCounterKey, strconv.Itoa(counter+1))
	if err != nil {
		fmt.Fprintf(w, "%s\n", fmt.Errorf("error. cound not set counter '%d'. %s", counter, err).Error())
		return
	}

	fmt.Fprintf(w, "%d\n", counter)
}

func (s *server) reset(w http.ResponseWriter, req *http.Request) {
	counter := 0
	err := s.cacheClient.Set(redisCounterKey, strconv.Itoa(counter))
	if err != nil {
		fmt.Fprintf(w, "%s\n", fmt.Errorf("error. cound not reset counter to '%d'. %s", counter, err).Error())
		return
	}

	fmt.Fprintf(w, "counter reset to %d\n", counter)
}

func (s *server) extra(w http.ResponseWriter, req *http.Request) {
	extra, err := s.cacheClient.Get(redisExtraKey)
	if err != nil {
		fmt.Fprintf(w, "%s\n", fmt.Errorf("error. redis returned an error. %s", err).Error())
	}

	fmt.Fprintf(w, "%s\n", extra)
}

func (s *server) set(w http.ResponseWriter, req *http.Request) {
	extras, ok := req.URL.Query()[redisExtraKey]
	if !ok || len(extras[0]) < 1 {
		fmt.Fprintf(w, "%s\n", fmt.Errorf("error. '%s' query param was not sent", redisExtraKey).Error())
		return
	}

	extra := extras[0]
	err := s.cacheClient.Set(redisExtraKey, extra)
	if err != nil {
		fmt.Fprintf(w, "%s\n", fmt.Errorf("error. cound not set '%s'. %s", extra, err).Error())
		return
	}

	fmt.Fprintf(w, "set %s\n", extra)
}

func (s *server) index(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, `<pre>
index:
- <a href="/counter">counter</a>
- <a href="/reset">reset</a>
</pre>
`)
}

type healthcheck struct {
	Healthy bool `json:"healthy"`
}

func (s *server) healthcheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	content, err := json.Marshal(healthcheck{true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(content))
}
