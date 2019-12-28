package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Hundemeier/go-sacn/sacn"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

// Storing universe channels
var universes map[int]chan<- [512]byte

func main() {
	// Flags for the Application
	configpath := flag.String("config", "sample-config.yaml", "Path to config file")
	flag.Parse()
	log.Println("rest2sacn started")

	// parsing config file
	c := readConfigFile(*configpath)
	err := c.Check()
	if err != nil {
		log.Fatal(err)
	}

	// set all required sACN Options/Parameters
	setupSACN(c)
	defer closeSACN()
	fmt.Println(universes)

	for {
		universes[69] <- [512]byte{255, 255}
		time.Sleep(1 * time.Second)
		universes[69] <- [512]byte{255, 0}
		log.Println("changed")
	}

	// rest API setup
	router := mux.NewRouter()
	router.HandleFunc("/index", indexHandler)
	router.HandleFunc("/sacn/reset", resetHandler)
	router.HandleFunc("/sacn/{universe:[1-63999]}/{channel:[1-512]}/{value:[0-255]}", sacnHandler)

	srv := &http.Server{
		Handler: router,
		Addr:    c.Ip,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

// closeSACN goes through all open sACN channels and closes them
func closeSACN() {
	for _, channel := range universes {
		close(channel)
	}
}

func setupSACN(c Config) {
	universes = make(map[int]chan<- [512]byte)

	trans, err := sacn.NewTransmitter("", [16]byte{1, 2, 3}, "rest2sacn")
	if err != nil {
		log.Fatal(err)
	}

	// Configuring all universes
	for _, univ := range c.Universe {
		var err error = nil
		universes[int(univ)], err = trans.Activate(univ)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Setup completed for universe", univ)

		// Send universe to be sent out unicast to destination
		trans.SetDestinations(univ, []string{c.Destination})
	}

}

// Config represents the YAML Config file of the application
type Config struct {
	Universe    []uint16
	Ip          string
	Destination string
}

// Check tests if the config is valid
func (c Config) Check() error {
	// TODO: Error Check if every required option is set
	return nil
}

func readConfigFile(path string) Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Couldn't read config file. Exiting...")
	}
	var c Config
	yaml.Unmarshal(data, &c)
	return c
}
