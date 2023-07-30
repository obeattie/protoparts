package main

import (
	"io/ioutil"
	"log"
	"os"

	protopartstest "github.com/obeattie/protoparts/test"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("error reading stdin: %v", err)
	}
	n := protopartstest.Fuzz(b)
	log.Printf("fuzz = %d", n)
}
