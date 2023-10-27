FROM golang:alpine as build
RUN apk add --no-cache --update git
ADD . /go/src/app
WORKDIR /go/src/app
RUN go get ./...
RUN go build \
    -a -tags timetzdata \
    -o nv \
    -ldflags="-s -w -X 'github.com/boggydigital/novus/cli.GitTag=`git describe --tags --abbrev=0`'" \
    main.go

FROM alpine:latest
COPY --from=build /go/src/app/nv /usr/bin/nv
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# https://en.wikipedia.org/wiki/Acta_Diurna
EXPOSE 59222
#app logs
VOLUME /var/log/novus
#app artifacts: content, redux
VOLUME /var/lib/novus

ENTRYPOINT ["/usr/bin/nv"]
CMD ["serve","-port", "59222", "-stderr"]
