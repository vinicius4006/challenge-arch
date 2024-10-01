FROM golang:1.23.1-alpine3.20 AS base
RUN apk add --update --no-cache mysql
RUN apk add --update --no-cache wget

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# suggested libraries by vs code for go, used to debug and test features in dev
RUN go install golang.org/x/tools/gopls@v0.16.2
RUN go install github.com/cweill/gotests/gotests@v1.6.0
RUN go install github.com/fatih/gomodifytags@v1.17.0
RUN go install github.com/josharian/impl@v1.4.0
RUN go install github.com/haya14busa/goplay/cmd/goplay@v1.0.0
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install honnef.co/go/tools/cmd/staticcheck@latest


# Install MySQL client for testing
RUN apk add mysql-client
RUN apk add moreutils

FROM base AS dev
CMD ["go", "run", "cmd/ordersystem/main.go"]