# build stage
FROM golang:1.18-alpine as builder
WORKDIR /app
COPY cmd cmd
COPY docs docs
COPY internal internal
COPY pkg pkg
COPY schema schema
COPY go.mod go.sum ./
RUN go mod download
RUN go build -v -o /app/main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY configs configs
COPY --from=builder /app/main .

CMD ["/app/main"]