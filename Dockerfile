FROM golang:1.16.4-buster AS builder
COPY . /gene_list_svc
WORKDIR /gene_list_svc
RUN go build .

FROM debian:buster-20210511-slim
RUN apt-get update && \
    apt-get install --no-install-recommends -y ca-certificates=20200601~deb10u2 && \
    rm -rf /var/lib/apt/lists/*
COPY --from=builder /gene_list_svc/gene_list_svc /usr/local/bin/
