package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func echo(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, "OK")
}

func main() {
	fmt.Printf("START")
	http.HandleFunc("/v2/animus/nft/webhook", echo)

	http.ListenAndServe(":3000", nil)
	fmt.Printf("STOP")
}
