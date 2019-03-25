FROM golang:1.12.1 as builder

WORKDIR /app/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch
COPY --from=builder /app/go-blockchain .
EXPOSE 5000
CMD ["./go-blockchain"]
