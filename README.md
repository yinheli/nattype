# NAT Type Detection

[![Go Reference](https://pkg.go.dev/badge/github.com/yinheli/nattype.svg)](https://pkg.go.dev/github.com/yinheli/nattype)
[![Go Report Card](https://goreportcard.com/badge/github.com/yinheli/nattype)](https://goreportcard.com/report/github.com/yinheli/nattype)
[![CI](https://github.com/yinheli/nattype/actions/workflows/ci.yml/badge.svg)](https://github.com/yinheli/nattype/actions/workflows/ci.yml)

A Go package for detecting NAT type using STUN protocol. This package helps determine the type of NAT (Network Address Translation) that exists between your application and the Internet.

## NAT Types

- `UdpBlocked`: UDP is blocked or network/STUN issue
- `OpenInternet`: Open Internet, no NAT
- `FullCone`: Full Cone NAT
- `RestrictedCone`: Restricted Cone NAT
- `PortRestrictedCone`: Port Restricted Cone NAT
- `Symmetric`: Symmetric NAT

## Installation

```bash
go get github.com/yinheli/nattype
```


## Usage

```go
import "github.com/yinheli/nattype"

natType, ip, err := nattype.DetectNATType(stunServer)
if err != nil {
	panic(err)
}
```

## Contributing

1. Fork the repository
2. Create a new branch
3. Make your changes
4. Run `go test` to ensure your changes pass all tests
5. Run `go fmt` to format your code
6. Run `go vet` to check your code for potential issues
7. Submit a pull request

## Contributors

![Contributors](https://contrib.rocks/image?repo=yinheli/nattype)

## License

MIT
