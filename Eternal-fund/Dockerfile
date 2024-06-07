FROM golang:alpine

RUN apk update && apk add --no-cache postgresql-client

WORKDIR /app

COPY . .

COPY .env /app

RUN go mod tidy

RUN go build -o eternal-fund

ENTRYPOINT [ "/app/eternal-fund" ]
