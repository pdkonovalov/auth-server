FROM golang:1.23

WORKDIR /usr/src/auth-server

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/auth-server ./cmd/auth-server/...

CMD go test -v ./cmd/auth-server/ && auth-server