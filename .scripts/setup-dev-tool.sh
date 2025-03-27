#!/bin/sh
cd "$(go env GOBIN)"

if command -v golangci-lint &> /dev/null
then
    printf "golangci-lint command is already installed\n"
else
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    printf "golangci-lint command install done\n"
fi

if command -v govulncheck &> /dev/null
then
    printf "govulncheck command is already installed\n"
else
    go install golang.org/x/vuln/cmd/govulncheck@latest
    printf "govulncheck command install done\n"
fi

if command -v swag &> /dev/null
then
    printf "swag command is already installed\n"
else
    # go install github.com/swaggo/swag/cmd/swag@latest
    go install github.com/go-swagger/go-swagger/cmd/swagger@latest

    printf "swagger command install done\n"
fi

if command -v godoc &> /dev/null
then
    printf "godoc command is already installed\n"
else
    go install golang.org/x/tools/cmd/godoc@latest
    printf "godoc command install done\n"
fi

if command -v godoc &> /dev/null
then
    printf "pkgsite command is already installed\n"
else
    go install golang.org/x/pkgsite/cmd/pkgsite@latest
    printf "pkgsite command install done\n"
fi
