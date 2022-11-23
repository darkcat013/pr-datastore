FROM golang:alpine

ARG config

RUN mkdir /app

COPY . /app

COPY ./${config}/ /app/config/

WORKDIR /app

RUN go mod download

RUN go build -o main .

CMD ["/app/main"]