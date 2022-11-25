package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

func main() {
	input, err := ioutil.ReadFile("./tools/marshal/input.yaml")
	if err != nil {
		fmt.Println(err)
	}
	encoded := base64.StdEncoding.EncodeToString(input)
	fmt.Println(encoded)
}
