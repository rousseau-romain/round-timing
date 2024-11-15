FROM golang:1.22.9

WORKDIR /app

COPY . .

RUN apt update

RUN apt-get install -y nodejs

RUN apt-get install -y npm

RUN go mod download

RUN go mod tidy 

RUN make build/install

RUN make build/templ

RUN make build/tailwind

RUN go build -o /main . 

EXPOSE 2468

CMD ["/main"]