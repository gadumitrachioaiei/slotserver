package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/pprof"

	"github.com/gadumitrachioaiei/slotserver/slot"
)

// Server is our http server
type Server struct {
	addr string
}

// New returns a new server for the address
func New(addr string) *Server {
	return &Server{addr}
}

// Start starts the atkins-diet slot machine server in the same goroutine
func (s *Server) Start() error {
	mux := http.NewServeMux()
	profile(mux)
	mux.Handle("/api/machines/atkins-diet/spins", slotMachineHandler{slot.NewMachine()})
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

// slotMachineHandler is the handler for atkins-diet slot machine spins
type slotMachineHandler struct {
	m *slot.Machine
}

func (h slotMachineHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(rw, "Bad http method", http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "Cannot read request body", http.StatusBadRequest)
		return
	}
	var body Request
	if err := json.Unmarshal(data, &body); err != nil {
		http.Error(rw, "Request body is not correct", http.StatusBadRequest)
		return
	}
	result, err := h.m.Bet(body.Chips, body.Bet)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	r := Response{Result: result, JWT: Request{UID: body.UID, Chips: result.Chips, Bet: result.Bet}}
	rData, err := json.Marshal(r)
	if err != nil {
		http.Error(rw, "Internal Server error", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Write(rData)
}

func profile(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
