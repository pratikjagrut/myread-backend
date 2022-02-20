FROM golang:1.16-alpine

RUN mkdir /app
ADD . /app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod vendor            

RUN GO111MODULE=on GOFLAGS=-mod=vendor go build -v -o main cmd/main.go

EXPOSE 8080

CMD [ "./main" ]