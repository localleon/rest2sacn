package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func sacnHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Mux should do our error checking
	data, err := parameterParse(vars)
	if err != nil {
		log.Println("Error on receiving sACN Data Request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// check and send
	if val, ok := UniverseOutput[data[0]]; ok {
		log.Println("Executing sACN Send with:", vars)
		val.data[data[1]-1] = byte(data[2]) // update current values with new values from parameters
		val.output <- val.data              // send out data to channel

		w.WriteHeader(http.StatusAccepted) // request processed, data was sent out
		UniverseOutput[data[0]] = val      // update global var to new data
	} else {
		// If the universe is not found in the config file, it cant be controlled via API
		log.Println("/sacn/send: Universe not found.")
		w.WriteHeader(http.StatusNotImplemented)
	}

}

// parameterParse given arguemnts from the mux endpoint /sacn/send and outputs them as 0: universe, 1: channelId, 2:channelVal
func parameterParse(v map[string]string) ([3]int, error) {
	univID, univErr := strconv.ParseInt(v["universe"], 10, 64)
	channelID, chanErr := strconv.ParseInt(v["channel"], 10, 64)
	channelVal, valErr := strconv.ParseInt(v["value"], 10, 64)
	if univErr != nil || chanErr != nil || valErr != nil {
		return [3]int{}, errors.New("parameterParse() couldn't convert all parameters to int64")
	}

	// Check if the parameters match the number range of typical dmx values
	if univID <= 63999 && univID >= 1 && channelID <= 512 && channelID >= 1 && channelVal <= 255 && channelVal >= 0 {
		return [3]int{int(univID), int(channelID), int(channelVal)}, nil
	}
	return [3]int{}, errors.New("parameterParse() found that on of the parameters is not in the range of the DMX Standard")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {

}
