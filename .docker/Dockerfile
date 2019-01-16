FROM golang:alpine as builder
RUN mkdir /go/src/semtest 
ADD . /go/src/semtest/
WORKDIR /go/src/semtest
RUN apk add --update \
    curl \
    && rm -rf /var/cache/apk/* \
    && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o crawler cmd/crawler/main.go
FROM scratch
COPY --from=builder /go/src/semtest /app/
WORKDIR /app
CMD ["./crawler"]