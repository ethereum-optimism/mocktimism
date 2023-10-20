FROM golang:1.21.1-alpine3.18 as builder

ARG VERSION=v0.0.0

WORKDIR /go/src/mocktimism

RUN apk --update --no-cache add ca-certificates gcc libtool make musl-dev git

COPY Makefile go.mod go.sum ./
RUN make init && go mod download 
COPY . .
RUN make build

FROM alpine:3.18

COPY --from=builder /go/src/mocktimism/bin/mocktimism /usr/local/bin
ENTRYPOINT ["mocktimism"]
CMD ["config"]
