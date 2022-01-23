FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /jaeger-app


FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /jaeger-app /jaeger-app
EXPOSE 8080
ENTRYPOINT ["/jaeger-app"]
