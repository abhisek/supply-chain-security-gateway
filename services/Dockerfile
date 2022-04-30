FROM golang:1.18-buster AS build

RUN apt-get update && apt-get install -y protobuf-compiler

WORKDIR /build

COPY go.mod go.sum Makefile ./

RUN go mod download && mkdir gen
RUN make oapi-codegen-install protoc-install

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN make

FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=build /build/out /app/server

ENV PATH "${PATH}:/app/server"
EXPOSE 9000 9001 9002

USER nonroot:nonroot
