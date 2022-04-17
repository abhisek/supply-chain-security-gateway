FROM golang:1.17-buster AS build

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN make

FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=build /build/out /app/server

ENV PATH "${PATH}:/app/server"
EXPOSE 9000

USER nonroot:nonroot