package main

import (
	"GoHereGoThere/balancerAlgos"
	"bytes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type LoadBalancer struct {
	BALANCER   balancerAlgos.Balancer
	HTTPClient *http.Client
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
	nextNode := b.BALANCER.Balance()
	scheme := r.URL.Scheme
	if r.URL.Scheme == "" {
		scheme = "https"
	}

	r.URL = &url.URL{
		Host: nextNode,
		Scheme: scheme,
	}
	log.Printf("Incoming request dispatching to:%v", nextNode)

	r.RequestURI = ""
	resp, err := b.HTTPClient.Do(r)
	if err != nil {
		log.Fatal("An error has occured")
	}

	for key,value := range resp.Header {
		w.Header().Set(key,strings.Join(value,","))
	}

	w.WriteHeader(resp.StatusCode)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.Fatal("An error has occured writing back to the user")
	}
}

func CreateLoadBalancer() *LoadBalancer {
	user_config := GetConfig()
	balancer := balancerAlgos.MapOfAlgos()[user_config.BalancerAlgo]
	balancer.SetNodes(user_config.Nodes)

	return &LoadBalancer{
		BALANCER: balancer,
		HTTPClient: &http.Client{},
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
