package main

import (
	"fmt"
    "log"
    //"net/http"
    "go-aviatrix/goaviatrix"
)

func main() {
	client, err := aviatrix.NewClient("rakesh", "av1@Tr1x", "https://13.126.166.7/v1/api")
	if err != nil {
		fmt.Println("Error")
  		log.Fatal(err)
	}
	if err==nil {
		fmt.Println(client.Username)
	}
}

