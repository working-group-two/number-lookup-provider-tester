FROM golang:1.21.4 AS builder
COPY . /go/src/
WORKDIR /go/src/
RUN set -Eeux && go mod download && go mod verify

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" ./...

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /go/src/number-lookup-provider-tester ./number-lookup-provider-tester

ENTRYPOINT ["./number-lookup-provider-tester"]
