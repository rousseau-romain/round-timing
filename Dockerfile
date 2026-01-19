FROM golang:1.25.5

ARG COMMIT_ID

ENV COMMIT_ID=${COMMIT_ID}

RUN echo "This image was built from commit: ${COMMIT_ID}"

WORKDIR /app

COPY . .

RUN apt update && apt-get install -y nodejs \
&& apt-get install -y npm \
&& apt-get install -y jq \
&& apt-get clean \
&& go install github.com/a-h/templ/cmd/templ@v0.2.793 \
&& npm install \
&& go mod download \
&& make build/commit-id ${COMMIT_ID} \
&& make build/tailwind \
&& make build/templ \
&& go build -o /main .

EXPOSE 2468

CMD ["/main"]