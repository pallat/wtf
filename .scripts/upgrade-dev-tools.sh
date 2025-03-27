#!/bin/sh
cd "$(go env GOBIN)"


go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
printf "golangci-lint command install done\n"

go install golang.org/x/vuln/cmd/govulncheck@latest
printf "govulncheck command install done\n"

go install github.com/go-swagger/go-swagger/cmd/swagger@latest
printf "swagger command install done\n"

go install golang.org/x/tools/cmd/godoc@latest
printf "godoc command install done\n"

go install golang.org/x/pkgsite/cmd/pkgsite@latest
printf "pkgsite command install done\n"
