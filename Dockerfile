FROM golang:1.10-alpine3.8 as builder
RUN apk update && apk add --update git
WORKDIR /go/src/github.com/endiangroup/transferwiser/

RUN go get github.com/golang/dep/cmd/dep
COPY Gopkg.toml .
COPY Gopkg.lock .
RUN dep ensure -vendor-only

COPY . .
RUN go install -v ./cmd

FROM alpine:3.8
RUN apk update && apk add ca-certificates
COPY --from=builder /go/bin/cmd transferwiser

EXPOSE 443
EXPOSE 80

ENV TRANSFERWISER_PORT=443
ENV TRANSFERWISER_LETSENCRYPTPORT=80

ENTRYPOINT ["./transferwiser"]
