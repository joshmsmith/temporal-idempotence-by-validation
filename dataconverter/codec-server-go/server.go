package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gogo/protobuf/jsonpb"
	commonpb "go.temporal.io/api/common/v1"
	converter "go.temporal.io/sdk/converter"

	dataconverter "idempotence-by-validation/dataconverter"
)

// Main
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("*", HomeHandler)
	//router.HandleFunc("/decode", DecodeHandler).Methods("GET","HEAD","PUT","PATCH","POST","DELETE")
	router.HandleFunc("/decode", DecodeHandler)

	log.Print("Serve Http on 8888")
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8888",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

// HomeHandler - just for testing
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/: called")
	params := r.URL.Query()
	for k, v := range params {
		log.Println("URL Params:", k, " => ", v)
	}
	log.Println("method:", r.Method) // request method
	fmt.Fprint(w, "HomeHandler")
	log.Println("/: done.")
}

// DecodeHandler -
func DecodeHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("/decode: called, method:", r.Method)

	params := r.URL.Query()
	for k, v := range params {
		log.Println("URL Params:", k, " => ", v)
	}

	// set response header options
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,PUT,PATCH,POST,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "content-type,x-namespace")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Keep-Alive", "timeout=5")

	if r.Method == "OPTIONS" {
		return
	}

	// Expecting a POST
	if r.Method == "POST" {
		r.ParseForm()
		w.Header().Set("Content-Type", "application/json")

		// from: converter.codec.ServeHTTP
		path := r.URL.Path

		if !strings.HasSuffix(path, remotePayloadCodecEncodePath) && !strings.HasSuffix(path, remotePayloadCodecDecodePath) {
			http.NotFound(w, r)
			return
		}

		var payloadspb commonpb.Payloads
		var err error
		if r.Body == nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = jsonpb.Unmarshal(r.Body, &payloadspb)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		payloads := payloadspb.Payloads

		e := codecHTTPHandler{}
		switch {
		case strings.HasSuffix(path, remotePayloadCodecDecodePath):
			if payloads, err = e.decode(payloads); err != nil {
				fmt.Println("codecHTTPHandler.decode returned error:", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		default:
			http.NotFound(w, r)
			return
		}

		// Encrypted to stdout for info:
		for i := range payloads {
			//fmt.Println(ColorBlue, "Encrypted: payload[", i, "]:", payloads[i], ColorReset)
			fmt.Println(ColorBlue, "Encrypted: payload[", i, "].Metadata:", payloads[i].Metadata, ColorReset)
			fmt.Println(ColorBlue, "Encrypted: payload[", i, "] metadata encoding string:", string(payloads[i].Metadata["encoding"]), ColorReset)
			fmt.Println(ColorBlue, "Encrypted: payload[", i, "] metadata encryption-key-id string:", string(payloads[i].Metadata["encryption-key-id"]), ColorReset)
			fmt.Println(ColorBlue, "Encrypted: payload[", i, "].Data:", payloads[i].Data, ColorReset)
			fmt.Println(ColorBlue, "Encrypted: payload[", i, "] data string:", string(payloads[i].Data), ColorReset)
		}
		fmt.Println("-----")

		// Decode/Decrypt:
		// Instance of local codec used in workflow client/workers
		codec := dataconverter.Codec{
			KeyID: "test-key-test-key-test-key-test!",
		}
		dpayloads, _ := codec.Decode(payloads)

		// Output to stdout for info
		for i := range dpayloads {
			//fmt.Println(ColorGreen, "Decrypted: payload[", i, "]:", dpayloads[i], ColorReset)
			fmt.Println(ColorGreen, "Decrypted: payload[", i, "].Metadata:", dpayloads[i].Metadata, ColorReset)
			fmt.Println(ColorGreen, "Decrypted: payload[", i, "] metadata encoding string:", string(dpayloads[i].Metadata["encoding"]), ColorReset)
			fmt.Println(ColorGreen, "Decrypted: payload[", i, "].Data:", dpayloads[i].Data, ColorReset)
			fmt.Println(ColorGreen, "Decrypted: payload[", i, "] data string:", string(dpayloads[i].Data), ColorReset)
		}

		// Response
		err = json.NewEncoder(w).Encode(commonpb.Payloads{Payloads: dpayloads})
		if err != nil {
			fmt.Println("json.NewEncoder(w).Encode(commonpb.Payloads{Payloads: dpayloads}) returned error:", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	} else {
		fmt.Fprint(w, "POST method required")
	}

	log.Println("/decode: done.")
	return
}

// codecHTTPHandler -
const remotePayloadCodecEncodePath = "/encode"
const remotePayloadCodecDecodePath = "/decode"

type codecHTTPHandler struct {
	codecs []converter.PayloadCodec
}

func (e *codecHTTPHandler) encode(payloads []*commonpb.Payload) ([]*commonpb.Payload, error) {
	var err error
	for i := len(e.codecs) - 1; i >= 0; i-- {
		if payloads, err = e.codecs[i].Encode(payloads); err != nil {
			return payloads, err
		}
	}
	return payloads, nil
}

func (e *codecHTTPHandler) decode(payloads []*commonpb.Payload) ([]*commonpb.Payload, error) {
	var err error
	for _, codec := range e.codecs {
		if payloads, err = codec.Decode(payloads); err != nil {
			return payloads, err
		}
	}
	return payloads, nil
}

// console colours
var ColorReset = "\033[0m"
var ColorRed = "\033[31m"
var ColorGreen = "\033[32m"
var ColorYellow = "\033[33m"
var ColorBlue = "\033[94m"
var ColorMagenta = "\033[35m"
var ColorCyan = "\033[36m"
var ColorWhite = "\033[37m"
