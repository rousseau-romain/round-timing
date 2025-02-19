FROM golang:1.22.9

WORKDIR /app

COPY . .

RUN apt update

RUN apt-get install -y nodejs

RUN apt-get install -y npm

RUN go install github.com/a-h/templ/cmd/templ@v0.2.793
	
RUN npm install

RUN go mod download

RUN make build/commit-id

RUN make build/tailwind

RUN make build/templ

RUN go build -o /main . 

EXPOSE 2468

CMD ["/main"]