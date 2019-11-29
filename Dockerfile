FROM golang:1.12.6-alpine3.10 as builder

ENV KADEMLIA /code

WORKDIR $KADEMLIA

COPY . .

RUN dep ensure

RUN go build . -o /main

FROM alpine:latest
WORKDIR /code
COPY --from=builder /main .
RUN chmod +x /main
ENTRYPOINT ["/code/main"]
