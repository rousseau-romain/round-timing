FROM golang:1.25.6

ARG SOURCE_COMMIT
ARG SOURCE_VERSION

WORKDIR /app

COPY . .

RUN apt update && apt-get install -y nodejs npm jq git \
&& apt-get clean \
&& go install github.com/a-h/templ/cmd/templ@v0.3.977 \
&& npm install \
&& go mod download \
&& make build/tailwind \
&& make build/templ \
&& COMMIT=${SOURCE_COMMIT:-$(git rev-parse HEAD 2>/dev/null || echo "unknown")} \
&& VERSION=${SOURCE_VERSION:-$(git describe --tags --abbrev=0 2>/dev/null || echo "")} \
&& go build -ldflags "\
  -X github.com/rousseau-romain/round-timing/config.commit=${COMMIT} \
  -X github.com/rousseau-romain/round-timing/config.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -X github.com/rousseau-romain/round-timing/config.version=${VERSION}" \
  -o /main .

EXPOSE 2468

CMD ["/main"]