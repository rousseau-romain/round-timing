FROM golang:1.25.6

WORKDIR /app

COPY . .

RUN apt update && apt-get install -y nodejs npm jq git \
&& apt-get clean \
&& go install github.com/a-h/templ/cmd/templ@v0.3.977 \
&& npm install \
&& go mod download \
&& make build/tailwind \
&& make build/templ \
&& go build -ldflags "\
  -X github.com/rousseau-romain/round-timing/config.commit=$(git rev-parse HEAD) \
  -X github.com/rousseau-romain/round-timing/config.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -X github.com/rousseau-romain/round-timing/config.vcsModified=false" \
  -o /main .

EXPOSE 2468

CMD ["/main"]