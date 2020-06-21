package httpserver

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gadumitrachioaiei/slotserver/slot"
)

// Server is our http server
type Server struct {
	addr     string
	register *prometheus.Registry
}

// New returns a new server for the address
func New(addr string, register *prometheus.Registry) *Server {
	return &Server{addr: addr, register: register}
}

// Start starts the atkins-diet slot machine server in the same goroutine
func (s *Server) Start() error {
	mux := http.NewServeMux()
	profile(mux)
	health(mux)
	mux.Handle("/metrics", promhttp.HandlerFor(s.register, promhttp.HandlerOpts{}))
	mux.Handle("/api/machines/atkins-diet/spins", newSlotMachineHandler(slot.NewMachine(), s.register))
	server := http.Server{Addr: s.addr, Handler: mux}
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
	return nil
}

// Request describes a request
type Request struct {
	UID   string // user id
	Chips int    // chips balance
	Bet   int    // bet size
}

// Response describes a response
type Response struct {
	*slot.Result
	JWT Request
}

// ErrorResponse describes an error response
type ErrorResponse struct {
	Code    int
	Message string
}

// slotMachineHandler is the handler for atkins-diet slot machine spins
type slotMachineHandler struct {
	m        *slot.Machine
	register prometheus.Registerer
	// metrics for our http server
	latencies  prometheus.Histogram
	latenciesS prometheus.Summary
	badReq     prometheus.Counter
}

func newSlotMachineHandler(m *slot.Machine, register prometheus.Registerer) *slotMachineHandler {
	latencies := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "slotserver",
		Name:      "http_request_duration_seconds",
		Help:      "Latency for a request",
		Buckets:   []float64{0.0001, 0.0002, 0.0003, 0.0005, 0.001, 0.05, 0.1, 0.5, 0.9, 1.0},
	})
	latenciesS := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace:  "slotserver",
		Name:       "http_request_duration_summary_seconds",
		Help:       "Latency for a request",
		Objectives: map[float64]float64{0.2: 0.05, 0.4: 0.05, 0.5: 0.05, 0.75: 0.03, 0.9: 0.01, 0.99: 0.001},
	})
	badReq := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "slotserver",
		Name:      "http_request_errors_total",
		Help:      "Count of http error requests",
	})
	if err := register.Register(latencies); err != nil {
		log.Printf("registering latencies: %v", err)
	}
	if err := register.Register(latenciesS); err != nil {
		log.Printf("registering latenciesS: %v", err)
	}
	if err := register.Register(badReq); err != nil {
		log.Printf("registering badReq: %v", err)
	}
	h := slotMachineHandler{m: m, register: register, latencies: latencies, latenciesS: latenciesS, badReq: badReq}
	return &h
}

func (h slotMachineHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	start := time.Now()
	defer func() {
		h.latencies.Observe(float64(time.Since(start)) / float64(time.Second))
		h.latenciesS.Observe(float64(time.Since(start)) / float64(time.Second))
	}()
	if req.Method != http.MethodPost {
		h.badReq.Inc()
		writeError(rw, "Bad http method", http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.badReq.Inc()
		writeError(rw, "Cannot read request body "+err.Error(), http.StatusBadRequest)
		return
	}
	var body Request
	if err := json.Unmarshal(data, &body); err != nil {
		h.badReq.Inc()
		writeError(rw, "Request body is not correct "+err.Error(), http.StatusBadRequest)
		return
	}
	result, err := h.m.Bet(body.Chips, body.Bet)
	if err != nil {
		h.badReq.Inc()
		writeError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	r := Response{Result: result, JWT: Request{UID: body.UID, Chips: result.Chips, Bet: result.Bet}}
	rData, err := json.Marshal(r)
	if err != nil {
		h.badReq.Inc()
		writeError(rw, "Internal Server error", http.StatusInternalServerError)
		return
	}
	rw.Write(rData)
}

func writeError(rw http.ResponseWriter, message string, status int) {
	rw.WriteHeader(status)
	data, err := json.Marshal(ErrorResponse{Code: status, Message: message})
	if err == nil {
		rw.Write(data)
	}
}

func profile(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func health(mux *http.ServeMux) {
	mux.HandleFunc("/ping", func(_ http.ResponseWriter, _ *http.Request) {})
}
