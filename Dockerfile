FROM golang:1.22

WORKDIR /app

COPY . .

RUN go mod tidy 

RUN go build -o /main . 

EXPOSE 2468

CMD ["/main"]