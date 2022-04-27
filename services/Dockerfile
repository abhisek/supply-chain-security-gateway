FROM golang:1.18-buster AS build

WORKDIR /build

COPY go.mod .
COPY go.sum .
COPY Makefile .

RUN go mod download
RUN make oapi-codegen-install

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN make

FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=build /build/out /app/server

ENV PATH "${PATH}:/app/server"
EXPOSE 9000 9001

USER nonroot:nonroot
