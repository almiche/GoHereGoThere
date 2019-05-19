package main

import (
	"go-fuckery/balancer_algos"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type LoadBalancer struct {
	balancer balancer_algos.Balancer
}

func main() {
	app := &LoadBalancer{
		balancer: balancer_algos.NewRoundRobin(Nodes()),
	}

	r := mux.NewRouter()
	r.HandleFunc("/", app.BalanceRequest)
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Print(srv.ListenAndServe())
}

func (b LoadBalancer)BalanceRequest(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	next_node := b.balancer.Balance()
	log.Printf("Incoming request dispatching to:%v",next_node)
}

func Nodes() []string{
	return []string {"192.250.78.1","134.45.65.76","192.45.65.02","156.46.45.21"}
}