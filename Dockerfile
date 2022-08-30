FROM golang:1.18-alpine

# install psql and make
RUN apk update
RUN apk add postgresql-client make cmake

WORKDIR /go/src/github.com/Viquad/crud-app
COPY cmd cmd
COPY configs configs
COPY docs docs
COPY internal internal
COPY pkg pkg
COPY schema schema
COPY go.mod go.sum ./

# build app
RUN go mod download
RUN go build -v -o /main ./cmd/main.go

# install migrate
# RUN go get github.com/golang-migrate/migrate/v4
# RUN go install github.com/golang-migrate/migrate/v4

CMD ["/main"]