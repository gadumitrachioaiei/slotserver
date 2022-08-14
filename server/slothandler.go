package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gadumitrachioaiei/slotserver/slot"
)

// slotMachineHandler is the handler for atkins-diet slot machine spins
type slotMachineHandler struct {
	m *slot.Machine
}

func (h slotMachineHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(rw, "Bad http method", http.StatusBadRequest)
		return
	}
	data, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "Cannot read request body", http.StatusBadRequest)
		return
	}
	var body Request
	if err := json.Unmarshal(data, &body); err != nil {
		http.Error(rw, "Request body is not correct", http.StatusBadRequest)
		return
	}
	result, err := h.m.Bet(req.Context(), body.UID, body.Bet)
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
