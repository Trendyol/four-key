package main

import (
	"four-key/cmd"
	"io/ioutil"
	"log"
)

func main() {
	log.SetOutput(ioutil.Discard)
	cmd.Execute()
}
