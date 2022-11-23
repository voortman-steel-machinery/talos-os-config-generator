package main

import (
	"encoding/json"
	"net/http"
)

type GeneratedConfig struct {
	ClusterName              string `json:"ClusterName"`
	ControlEndpoint          string `json:"ControlEndpoint"`
	EncodedControlplanConfig string `json:"EncodedControlplanConfig"`
	EncodedWorkerConfig      string `json:"EncodedWorkerConfig"`
	EncodedTalosConfig       string `json:"EncodedTalosConfig"`
}

func GenerateConfigHandler(w http.ResponseWriter, r *http.Request) {
	generatedConfig := GeneratedConfig{"TestName", "TestEndpoint", "TestEncodedControlplanConfig", "TestEncodedWorkerConfig", "TestEncodedTalosConfig"}

	json.NewEncoder(w).Encode(generatedConfig)
}

func main() {
	http.HandleFunc("/generate", GenerateConfigHandler)
	http.ListenAndServe(":5050", nil)
}
