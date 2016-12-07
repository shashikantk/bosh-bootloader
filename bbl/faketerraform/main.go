package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	backendURL string
)

func main() {
	if os.Args[1] == "fast-fail" {
		log.Fatal("failed to terraform")
	}

	if os.Args[1] == "output" {
		resp, err := http.Get(fmt.Sprintf("%s/output/%s", backendURL, os.Args[2]))
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			panic(err)
		}
		fmt.Print(string(body))
	}

	if os.Args[1] == "apply" {
		err := ioutil.WriteFile("terraform.tfstate", []byte(`{"key":"value"}`), os.ModePerm)
		if err != nil {
			panic(err)
		}

		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		fmt.Printf("working directory: %s\n", dir)
		fmt.Printf("terraform %s/n", removeBrackets(fmt.Sprintf("%+v", os.Args)))
	}
}

func removeBrackets(contents string) string {
	contents = strings.Replace(contents, "[", "", -1)
	contents = strings.Replace(contents, "]", "", -1)
	return contents
}