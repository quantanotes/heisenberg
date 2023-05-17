FROM golang:1.20-alpine

RUN apk add --no-cache --update go gcc g++

ENV CGO_ENABLED=1

ENV GOFLAGS=-mod=mod

WORKDIR /app

COPY . .

RUN go build -o heisenberg . 

EXPOSE 8080

CMD ["./heisenberg"]
