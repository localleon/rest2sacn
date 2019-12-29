package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Hundemeier/go-sacn/sacn"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

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

	// rest API setup
	router := mux.NewRouter()
	router.HandleFunc("/sacn/reset/{universe}", resetHandler)
	router.HandleFunc("/sacn/send/{universe}/{channel}/{value}", sacnHandler)

	srv := &http.Server{
		Handler: router,
		Addr:    c.Ip,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

// UniverseOutput stores all configured universes and their current data
var UniverseOutput map[int]Universe

// Universe represents a DMX Universe sent out over sACN
type Universe struct {
	number uint16
	data   [512]byte
	output chan<- [512]byte
}

// Send outputs the currently stored data in the universe
func (u *Universe) Send() {
	u.output <- u.data
}

// setupSACN activates all universes and makes them ready for transmitting
func setupSACN(c Config) {
	UniverseOutput = make(map[int]Universe) // global storage to access from rest Handlers
	// setup main transmitter.
	trans, err := sacn.NewTransmitter("", [16]byte{1, 2, 3}, "rest2sacn")
	if err != nil {
		log.Fatal(err)
	}
	// Configuring all universes
	for _, univ := range c.Universe {
		// Activate and store channel in global var
		ch, err := trans.Activate(univ)
		if err != nil {
			log.Fatal(err)
		}
		u := Universe{
			number: univ,
			output: ch,
		}
		UniverseOutput[int(univ)] = u
		// Send universe to be sent out unicast to destination
		trans.SetDestinations(u.number, []string{c.Destination})
		// finished
		log.Println("Setup completed for universe", u.number)
	}
}

// closeSACN goes through all open sACN channels and closes them
func closeSACN() {
	for _, univ := range UniverseOutput {
		close(univ.output)
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
