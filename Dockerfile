FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go mod tidy 
# try again

RUN go build -o /main . 
#if . is main package -> /app

EXPOSE 2468

CMD ["/main"]