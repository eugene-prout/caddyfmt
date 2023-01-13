package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

type Request struct {
	UnformattedFile string
}

type Response struct {
	FormattedFile string
}

func formatRequest(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	if request.Method == "POST" {
		writer.WriteHeader(http.StatusOK)

		decoder := json.NewDecoder(request.Body)
		decoder.DisallowUnknownFields()

		var incoming Request
		if err := decoder.Decode(&incoming); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("Bad Request"))
			return
		}

		var outputFile = string(caddyfile.Format([]byte(incoming.UnformattedFile)))

		outgoing := Response{
			FormattedFile: outputFile,
		}

		json.NewEncoder(writer).Encode(outgoing)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Bad Request"))
	}
	log.Println("Endpoint Hit: /format")
}

func main() {
	fmt.Println("Starting server")
	http.HandleFunc("/api/format", formatRequest)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
