FROM golang:1.13


ARG golangci_lint_version=1.21.0

# NOTE: GolangCI-Lint README says "Please, do not install golangci-lint by go get"
# See: https://github.com/golangci/golangci-lint#go
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v${golangci_lint_version}

# Install from this repository
ADD . /source
WORKDIR /source
RUN go build -o /entrypoint && rm -rf /source/*
WORKDIR /

ENTRYPOINT ["/entrypoint"]
