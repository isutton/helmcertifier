FROM golang:1.15 AS build

WORKDIR /tmp/src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN ./hack/build.sh

FROM fedora:31

COPY --from=build /tmp/src/out/helmcertifier /app/helmcertifier

ENTRYPOINT ["/app/helmcertifier"]
