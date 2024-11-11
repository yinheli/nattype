package nattype

import (
	"fmt"
	"net"

	"github.com/pion/stun/v3"
)

// NATType represents the type of NAT (Network Address Translation).
type NATType int

const (
	// UdpBlocked indicates UDP is blocked or there are network/STUN issues
	UdpBlocked NATType = iota
	// OpenInternet indicates direct internet access without NAT
	OpenInternet
	// FullCone indicates a full cone NAT where any external host can send packets
	// to the internal host once the internal host has sent a packet to any external host
	FullCone
	// RestrictedCone indicates a restricted cone NAT where an external host can send packets
	// to the internal host only if the internal host has previously sent a packet to that external host
	RestrictedCone
	// PortRestrictedCone indicates a port restricted cone NAT where an external host can send packets
	// to the internal host only if the internal host has previously sent a packet to that external host's IP and port
	PortRestrictedCone
	// Symmetric indicates a symmetric NAT where each request from the same internal IP:port
	// to a specific destination IP:port is mapped to a unique external source IP:port
	Symmetric
)

// String returns the string representation of the NAT type.
func (n NATType) String() string {
	return []string{"UdpBlocked", "OpenInternet", "FullCone", "RestrictedCone", "PortRestrictedCone", "Symmetric"}[n]
}

// DetectNATType detects the NAT type by performing multiple STUN tests.
//
// The function performs three STUN tests to determine the NAT type:
//  1. First test establishes basic connectivity and gets the initial mapped address
//  2. Second test checks if the mapped address changes with a new connection
//  3. Third test verifies the behavior with the original connection
//
// Parameters:
//   - stunServer: The STUN server address in the format "host:port"
//
// Returns:
//   - NATType: The detected NAT type
//   - net.IP: The external IP address as seen by the STUN server
//   - error: Any error that occurred during detection
//
// Example:
//
//	natType, ip, err := nattype.DetectNATType("stun.miwifi.com:3478")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("NAT Type: %s, External IP: %s\n", natType, ip)
func DetectNATType(stunServer string) (NATType, net.IP, error) {
	xorAddr1, err := doStunTest(stunServer)
	if err != nil {
		return UdpBlocked, nil, err
	}

	mappedIP := xorAddr1.IP

	xorAddr2, err := doStunTest(stunServer)
	if err != nil {
		return UdpBlocked, nil, fmt.Errorf("second STUN test error: %v", err)
	}

	if !xorAddr1.IP.Equal(xorAddr2.IP) || xorAddr1.Port != xorAddr2.Port {
		return Symmetric, mappedIP, nil
	}

	xorAddr3, err := doStunTest(stunServer)
	if err != nil {
		return RestrictedCone, mappedIP, nil
	}

	if xorAddr1.IP.Equal(xorAddr3.IP) && xorAddr1.Port == xorAddr3.Port {
		return FullCone, mappedIP, nil
	}

	return PortRestrictedCone, mappedIP, nil
}

// doStunTest performs a single STUN test and returns the mapped address.
func doStunTest(stunServer string) (*stun.XORMappedAddress, error) {
	conn, err := stun.Dial("udp4", stunServer)
	if err != nil {
		return nil, fmt.Errorf("STUN dial error: %v", err)
	}
	defer conn.Close() //nolint:errcheck

	var xorAddr stun.XORMappedAddress
	message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)

	err = conn.Do(message, func(res stun.Event) {
		if res.Error != nil {
			err = fmt.Errorf("STUN test error: %v", res.Error)
			return
		}
		if getErr := xorAddr.GetFrom(res.Message); getErr != nil {
			err = fmt.Errorf("failed to get XOR-MAPPED-ADDRESS: %v", getErr)
			return
		}
	})

	if err != nil {
		return nil, err
	}

	return &xorAddr, nil
}
