FROM golang:1.25.6

WORKDIR /app

COPY . .

RUN apt update && apt-get install -y nodejs \
&& apt-get install -y npm \
&& apt-get install -y jq \
&& apt-get clean \
&& go install github.com/a-h/templ/cmd/templ@v0.3.977 \
&& npm install \
&& go mod download \
&& make build/tailwind \
&& make build/templ \
&& go build -buildvcs -o /main .

EXPOSE 2468

CMD ["/main"]