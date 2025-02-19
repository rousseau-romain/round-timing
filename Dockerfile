FROM golang:1.22.9

WORKDIR /app

COPY . .

RUN apt update && apt-get install -y nodejs \
&& apt-get install -y npm \
&& apt-get install -y jq \
&& apt-get clean \
&& go install github.com/a-h/templ/cmd/templ@v0.2.793 \
&& npm install \
&& go mod download \
&& make build/tailwind \
&& make build/templ \
&& go build -o /main .

EXPOSE 2468

CMD ["/main"]