FROM golang:1.18-alpine

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql and make
RUN apk update
RUN apk add postgresql-client make cmake

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# build app 
RUN go mod download
RUN go build -o ./main ./cmd/main.go
# install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

CMD ["./main"]