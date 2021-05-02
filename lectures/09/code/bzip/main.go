package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/dbeliakov/mipt-golang-course/lectures/09/code/cgo"
)

func main() {
	writer := cgo.NewWriter(os.Stdout)
	defer func() {
		err := writer.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	_, err = writer.Write(in)
	if err != nil {
		log.Fatal(err)
	}
}
