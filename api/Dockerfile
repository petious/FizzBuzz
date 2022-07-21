FROM golang:1.16-alpine


WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY fizzbuzz ./fizzbuzz
COPY sqlUtils ./sqlUtils
COPY docs ./docs
COPY *.go ./


RUN go build -o /fizzbuzz

EXPOSE 8080

CMD [ "/fizzbuzz" ]

ENV PORT=8080
ENV MYSQL_PASSWORD="Jvend2?"
ENV MYSQL_DATABASE="fizzbuzz"
ENV MYSQL_USER="api"
ENV MYSQL_PORT=3306
ENV MYSQL_HOST="db"
ENV MAX_SEC_RETRY=60