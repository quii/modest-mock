package main

import (
	"flag"
	"fmt"
	"github.com/quii/modest-mock"
	"log"
	"os"
)

// /home/qui/go/src/github.com/mergermarket/run-amqp/connection/amqpchannel.go
// AMQPChannel

func main() {
	var filePath, interfaceName string

	flag.StringVar(&filePath, "path", "", "path to where interface declaration file")
	flag.StringVar(&interfaceName, "name", "", "name of interface you wish to mock")

	flag.Parse()

	f, err := os.Open(filePath)

	if err != nil {
		log.Printf("problem opening file %s %v", filePath, err)
		os.Exit(1)
	}

	mock, err := modestmock.New(f, interfaceName)

	if err != nil {
		log.Printf("problem creating mock %s defined in %s %v", interfaceName, filePath, err)
		os.Exit(1)
	}

	generatedMock, err := modestmock.GenerateMockCode(mock)

	if err != nil {
		log.Printf("problem generating mock code from %v %v", mock, err)
		os.Exit(1)
	}

	fmt.Print(generatedMock)
}
