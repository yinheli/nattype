package nattype_test

import (
	"fmt"
	"log"

	"github.com/yinheli/nattype"
)

func Example() {
	// Use a STUN server to detect NAT type
	stunServer := "stun.miwifi.com:3478"

	natType, ip, err := nattype.DetectNATType(stunServer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("NAT Type: %s\n", natType)
	fmt.Printf("External IP: %v\n", ip)
}

func ExampleDetectNATType() {
	stunServer := "stun.miwifi.com:3478"
	natType, ip, err := nattype.DetectNATType(stunServer)
	if err != nil {
		log.Fatal(err)
	}

	switch natType {
	case nattype.UdpBlocked:
		fmt.Println("UDP is blocked")
	case nattype.OpenInternet:
		fmt.Println("Direct internet access")
	case nattype.FullCone:
		fmt.Println("Full cone NAT")
	case nattype.RestrictedCone:
		fmt.Println("Restricted cone NAT")
	case nattype.PortRestrictedCone:
		fmt.Println("Port restricted cone NAT")
	case nattype.Symmetric:
		fmt.Println("Symmetric NAT")
	}

	fmt.Printf("External IP: %v\n", ip)
}
