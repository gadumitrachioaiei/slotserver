package server

import (
	"net/http"

	"github.com/gadumitrachioaiei/slotserver/slot"
)

// Server is our http server
type Server struct {
	addr        string
	userService slot.UserService
}

// New returns a new server for the address
func New(addr string, userService slot.UserService) *Server {
	return &Server{addr: addr, userService: userService}
}

// Start starts the atkins-diet slot machine server in the same goroutine
func (s *Server) Start() error {
	mux := http.NewServeMux()
	slotMachineHandler := slotMachineHandler{slot.NewMachine(s.userService)}
	mux.Handle("/api/machines/atkins-diet/spins", slotMachineHandler)
	server := http.Server{Addr: s.addr, Handler: mux}
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
	return nil
}
