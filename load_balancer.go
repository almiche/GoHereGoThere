package main

import (
	"GoHereGoThere/balancer_algos"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type LoadBalancer struct {
	balancer balancer_algos.Balancer
}

type BalancerConfig struct {
	BalancerAlgo string   `yaml:"BalancerAlgo"`
	Nodes        []string `yaml:"Nodes"`
}

func main() {
	load_balancer := CreateLoadBalancer()

	r := mux.NewRouter()
	r.HandleFunc("/", load_balancer.BalanceRequest)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Print(srv.ListenAndServe())
}

func (b LoadBalancer) BalanceRequest(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	next_node := b.balancer.Balance()
	log.Printf("Incoming request dispatching to:%v", next_node)
}

func CreateLoadBalancer() *LoadBalancer {
	user_config := GetConfig()
	balancer := balancer_algos.MapOfAlgos()[user_config.BalancerAlgo]
	balancer.SetNodes(user_config.Nodes)

	return &LoadBalancer{
		balancer: balancer,
	}
}

func GetConfig() *BalancerConfig {
	configuration := BalancerConfig{}

	yaml_file, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("An error has occured reading the file you have provided")
	}

	err = yaml.Unmarshal([]byte(yaml_file), &configuration)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return &configuration
}
