FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

RUN go mod vendor 

RUN go mod download

COPY . .

RUN go build -o /main

EXPOSE 2468

CMD ["/main"]
