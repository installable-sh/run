FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go get -u github.com/installable-sh/lib@main
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o RUN .

FROM scratch
COPY --from=builder /build/RUN /usr/local/bin/RUN
ENTRYPOINT ["/usr/local/bin/RUN"]
