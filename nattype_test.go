package nattype

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testStunServer = "stun.miwifi.com:3478"

func TestDetectNATType(t *testing.T) {
	natType, ip, err := DetectNATType(testStunServer)
	require.NoError(t, err, "DetectNATType should not return error")

	assert.GreaterOrEqual(t, int(natType), int(UdpBlocked), "NAT type should be >= UdpBlocked")
	assert.LessOrEqual(t, int(natType), int(Symmetric), "NAT type should be <= Symmetric")

	require.NotNil(t, ip, "IP should not be nil")
	assert.True(t, ip.IsGlobalUnicast(), "IP should be a global unicast address")
	assert.NotNil(t, ip.To4(), "IP should be IPv4")

	t.Logf("Detected NAT Type: %s", natType)
	t.Logf("External IP: %s", ip.String())
}

func TestDetectNATTypeWithInvalidServer(t *testing.T) {
	tests := []struct {
		name   string
		server string
	}{
		{
			name:   "invalid domain",
			server: "invalid.server:3478",
		},
		{
			name:   "invalid port",
			server: "stun.invalid.com:invalid",
		},
		{
			name:   "empty server",
			server: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			natType, ip, err := DetectNATType(tt.server)
			assert.Error(t, err, "Should return error")
			assert.Equal(t, UdpBlocked, natType, "Should return UdpBlocked")
			assert.Nil(t, ip, "Should return nil IP")
		})
	}
}

func TestMultipleDetections(t *testing.T) {
	var prevNatType NATType
	var prevIP net.IP

	for i := 0; i < 3; i++ {
		natType, ip, err := DetectNATType(testStunServer)
		require.NoError(t, err, "Detection %d should not fail", i+1)

		if i > 0 {
			assert.Equal(t, prevNatType, natType, "NAT type should be consistent across detections")
			assert.True(t, ip.Equal(prevIP), "IP should be consistent across detections")
		}

		prevNatType = natType
		prevIP = ip

		t.Logf("Detection %d - NAT Type: %s, IP: %s", i+1, natType, ip.String())
	}
}

func BenchmarkDetectNATType(b *testing.B) {
	b.Skip("skipping benchmark")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		natType, ip, err := DetectNATType(testStunServer)
		require.NoError(b, err, "Detection should not fail")
		require.GreaterOrEqual(b, int(natType), int(UdpBlocked))
		require.LessOrEqual(b, int(natType), int(Symmetric))
		require.NotNil(b, ip)
		require.True(b, ip.IsGlobalUnicast())
	}
}
