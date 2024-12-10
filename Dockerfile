FROM golang:latest AS builder

ARG VERSION
ARG COMMIT

ADD . $GOPATH/src/github.com/dpc-sdp/bay-section-ip-controller/

WORKDIR $GOPATH/src/github.com/dpc-sdp/bay-section-ip-controller

ENV CGO_ENABLED 0

RUN apt-get install ca-certificates

RUN go build -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" -o build/bay-section-ip-controller

FROM scratch

ARG PORT=80

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/dpc-sdp/bay-section-ip-controller/build/bay-section-ip-controller /usr/local/bin/bay-section-ip-controller

EXPOSE $PORT

CMD [ "bay-section-ip-controller", "-port", $PORT ]