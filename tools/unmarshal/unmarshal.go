package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	input, err := ioutil.ReadFile("./tools/unmarshal/input.txt")
	if err != nil {
		fmt.Println(err)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(input))
	if err != nil {
		log.Fatalf("failed to decode config: %s", err)
	}
	// write the config to a file
	if err = os.WriteFile("output.yaml", decoded, 0o600); err != nil {
		log.Fatalf("failed to write config: %s", err)
	}
}
