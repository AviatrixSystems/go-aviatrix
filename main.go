package main

import (
	"fmt"
    "log"
    "crypto/tls"
    "net/http"
    "github.com/go-aviatrix/goaviatrix"
)

func main() {

    tr := &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    //client,_ := goaviatrix.NewClient(c.Username, c.Password, c.ControllerIP, &http.Client{Transport: tr})
	client, err := goaviatrix.NewClient("rakesh", "av1@Tr1x", "13.126.166.7", &http.Client{Transport: tr})
	if err != nil {
		fmt.Println("Error")
  		log.Fatal(err)
	}
	if err==nil {
		fmt.Println(client.CID)
	}

	err = client.CreateGateway(&goaviatrix.Gateway{
		Action: "connect_container",
		CloudType: 1,
		AccountName: "devops",
		GwName: "avtxgw3",
		VpcID: "vpc-0d7b3664",
		VpcRegion: "ap-south-1",
		VpcSize: "t2.micro",
		VpcNet: "avtxgw3_sub1~~10.3.0.0/24~~ap-south-1a",
		})
	if err!=nil {
		fmt.Println(err)
	}

	// err1 := client.DeleteGateway(&goaviatrix.Gateway{
	// 	CloudType: 1,
	// 	GwName: "avtxgw1",
	// 	})
	// if err1!=nil {
	// 	fmt.Println(err1)
	// }

	// err = client.CreateTunnel(&goaviatrix.Tunnel{
	//  	VpcName1: "avtxgw1",
	//  	VpcName2: "avtxgw2",
	// })

	// if err!=nil {
	// 	fmt.Println(err)
	// }

	// tun, err := client.GetTunnel(&goaviatrix.Tunnel{
	//  	VpcName1: "avtxgw1",
	//  	VpcName2: "avtxgw2",
	// })

	// if err!=nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(tun.VpcName1, tun.VpcName2)

	// err = client.DeleteTunnel(&goaviatrix.Tunnel{
	//  	VpcName1: "avtxgw1",
	//  	VpcName2: "avtxgw2",
	// })

	// if err!=nil {
	// 	fmt.Println(err)
	// }
}

