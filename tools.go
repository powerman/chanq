//go:build generate

//go:generate mkdir -p .gobincache
//go:generate -command GOINSTALL env "GOBIN=$PWD/.gobincache" go install

package tools

//go:generate GOINSTALL github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2
//go:generate GOINSTALL github.com/mattn/goveralls@v0.0.11
