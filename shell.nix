{ pkgs ? import (builtins.fetchGit {
         # Descriptive name to make the store path easier to identify                
         name = "dev-go";                                                 
         url = "https://github.com/NixOS/nixpkgs";                       
         ref = "refs/heads/nixpkgs-unstable";                     
         rev = "7cf5ccf1cdb2ba5f08f0ac29fc3d04b0b59a07e4"; 
}) {} }:

with pkgs;

mkShell {
  buildInputs = [
    awscli2
    clang-tools
    gitlint
    gnupg
    go_1_19
    go-tools
    go-mockery
    gogetdoc
    golangci-lint
    goreleaser
    gosec
    gotools
    #gocritic
    gofumpt
    golint
    #goreturns
    jq
    mysql80
    openapi-generator-cli
    postgresql
    pre-commit
    protobuf
    protoc-gen-go
    protoc-gen-go-grpc
    sops
    #structlog
  ];

  shellHook =
    ''
      # Setup the terminal prompt.
      export PS1="(nix-shell) \W $ "
      export AWS_REGION=us-east-1
      export AWS_PROFILE=root

      # Setup the binaries installed via `go install` to be accessible globally.
      export PATH="$(go env GOPATH)/bin:$PATH"

      # Create the temporary folder for AWS CLI.
      mkdir -p $HOME/.aws

      # Install pre-commit hooks.
      pre-commit install

      # Install Go binaries.
      which protoc-gen-grpc-gateway || go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0
      which protoc-gen-openapiv2 || go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.5.0
      which enumer || go install github.com/dmarkham/enumer@v1.5.3
      which gocritic || go install github.com/go-critic/go-critic/cmd/gocritic@latest
      which goreturns || go install github.com/sqs/goreturns@latest
      which swag || go get -u github.com/swaggo/swag
      which mockgen || go install github.com/golang/mock/mockgen@v1.6.0
      
      # Add the repo shared gitconfig
      git config --local include.path ../.gitconfig

      # Clear the terminal screen.
      clear
    '';
}
