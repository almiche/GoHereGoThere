# GoHereGoThere

![GoHereGoThere](https://user-images.githubusercontent.com/14715156/57989207-6f5a8680-7a65-11e9-899f-465cd1aa26a5.png)

# What is this? (WIP)

GoHereGoThere is a simple application layer load balancer. It's been designed to be extensible and simple to use.

# How can I use this ?

Simple supply the application with a yaml file such as the one below, `./load_balancer.go examples/test.yaml`:
```
BalancerAlgo:
  RoundRobin
Nodes:
  - 192.123.43.2
  - 192.023.34.2
  - 182.345.43.2
  - 136.54.65.7
```
The application is written in such a way that one can decide on which load balancing algorithm is used on the fly by changing the BalancerAlgo variable in the configuration YAML (In this case round robin will be used).
