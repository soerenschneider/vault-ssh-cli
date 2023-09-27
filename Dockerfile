FROM golang:1.21.1 as builder

ENV MODULE=github.com/soerenschneider/vault-ssh-cli
ENV CGO_ENABLED=0

WORKDIR /build/
ADD go.mod go.sum /build/
RUN go mod download

ADD . /build/
RUN make build

FROM gcr.io/distroless/base
COPY --from=builder /build/vault-ssh-cli /vault-ssh-cli
ENTRYPOINT ["/vault-ssh-cli"]
