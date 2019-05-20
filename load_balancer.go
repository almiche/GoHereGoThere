package main

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/almiche/GoHereGoThere/balancerAlgos"
	"github.com/gorilla/mux"
)

type LoadBalancer struct {
	Balancer   balancerAlgos.Balancer
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
	nextNode := b.Balancer.Balance()
	scheme := r.URL.Scheme
	if r.URL.Scheme == "" {
		scheme = "https"
	}

	r.URL = &url.URL{
		Host:   nextNode,
		Scheme: scheme,
	}
	log.Printf("Incoming request dispatching to:%v", nextNode)

	r.RequestURI = ""
	resp, err := b.HTTPClient.Do(r)
	if err != nil {
		log.Fatal("An error has occured")
	}

	for key, value := range resp.Header {
		w.Header().Set(key, strings.Join(value, ","))
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
	userConfig := GetConfig()
	balancer := balancerAlgos.MapOfAlgos()[userConfig.BalancerAlgo]
	balancer.SetNodes(userConfig.Nodes)

	return &LoadBalancer{
		Balancer:   balancer,
		HTTPClient: &http.Client{},
	}
}

func GetConfig() *BalancerConfig {
	configuration := BalancerConfig{}

	path,_ := filepath.Abs(os.Args[1])
	yamlFiles, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("An error has occured reading the file you have provided")
	}

	err = yaml.Unmarshal([]byte(yamlFiles), &configuration)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return &configuration
}
