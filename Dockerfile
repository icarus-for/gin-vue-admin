FROM golang:1.22.2-alpine

WORKDIR /app

COPY server/ /app/

RUN go mod tidy && go build -o main main.go

EXPOSE 8888

CMD ["./main"]
