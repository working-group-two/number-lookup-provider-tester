FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git upx

COPY . /go/src/
WORKDIR /go/src/

RUN set -Eeux && go mod download && go mod verify

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" ./...
RUN upx --best --lzma /go/src/number-lookup-provider-tester

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /go/src/number-lookup-provider-tester ./number-lookup-provider-tester

ENTRYPOINT ["./number-lookup-provider-tester"]
