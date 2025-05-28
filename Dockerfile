FROM golang:1.24

RUN apt-get update && apt-get install -y postgresql-client
RUN mkdir -p /app/data
COPY ./src/app/database/data.sql /app/data/data.sql
WORKDIR /app

COPY /go.* ./
RUN go mod download

COPY ./ .

RUN go build -o main ./src/cmd/main.go

EXPOSE 8080

CMD ["./main"]