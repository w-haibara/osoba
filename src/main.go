package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"osoba/config"
	"osoba/deploy"
	"osoba/handler"

	"github.com/k0kubun/pp"
)

var configFile = flag.String("f", "config.json", "path to the configuration file")

func main() {
	flag.Parse()

	c := configure()
	chanDeployInfo := make(chan deploy.Info)

	http.Handle("/", handler.Logging(handler.CheckMethods(handler.Auth(*c.Auth, http.HandlerFunc(handler.Main)), http.MethodGet)))
	http.Handle("/deploy", handler.Logging(handler.CheckMethods(handler.Webhook(*c.DB, chanDeployInfo), http.MethodPost)))

	go deploy.AwaitDeploy(chanDeployInfo)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func configure() config.Config {
	log.Printf("config file: %v\n", *configFile)

	json, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Panic(err)
	}

	c, err := config.Configure(json)
	if err != nil {
		log.Panic(err)
	}

	pp.Println(c)

	c.DB.InitDB()

	return c
}
