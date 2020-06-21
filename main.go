package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/prometheus/client_golang/prometheus"

	_ "github.com/gadumitrachioaiei/slotserver/docs"
	"github.com/gadumitrachioaiei/slotserver/httpserver"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "", "port where to run the app")
	flag.Parse()
	if port == "" {
		log.Fatal("no port given")
	}
	register := prometheus.NewRegistry()
	register.MustRegister(prometheus.NewGoCollector())
	register.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{Namespace: "slotserver"}))
	// we register our service with the local consul agent
	//if err := consulRegister(port, fmt.Sprintf("http://localhost:%s/ping", port)); err != nil {
	//	log.Fatal("cannot register with consul", err)
	//	os.Exit(1)
	//}
	s := httpserver.New(":"+port, register)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}

// handleNotify handles notifications for talking to systemd manager
func handleNotify() {
	if ok, err := daemon.SdNotify(false, daemon.SdNotifyReady); err != nil || !ok {
		log.Fatal(ok, err)
	}
	c := make(chan os.Signal, 100)
	signal.Notify(c)
	go func() {
		for {
			s := <-c
			switch s {
			case syscall.SIGTERM:
				if ok, err := daemon.SdNotify(false, daemon.SdNotifyStopping); err != nil || !ok {
					log.Fatal("cannot send shutdown notification")
				}
				os.Exit(0)
			default:
				log.Println("we received signal", s)
			}
		}
	}()
}

// consulRegister registers a service with a health check with the local consul agent
func consulRegister(registerPort string, healthURL string) error {
	u := "http://localhost:8500/v1/agent/service/register?replace-existing-checks=true"
	body := fmt.Sprintf(`
	{
	  "id": "slotserver%s",
	  "name": "slotserver",
	  "tags": [
		"games",
		"metrics",
        "urlprefix-/ weight=0.2"
	  ],
	  "port": %s,
	  "check": {
		"http": "%s",
		"interval": "10s",
		"timeout": "5s"
	  }
	}`, registerPort, registerPort, healthURL)
	c := http.Client{
		Timeout: time.Second,
	}
	req, err := http.NewRequest(http.MethodPut, u, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("cannot make up the request object: %v", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("cannot make the request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("cannot read response: %v", err)
		}
		return fmt.Errorf("bad status response: %d %s", resp.StatusCode, r)
	}
	return nil
}
