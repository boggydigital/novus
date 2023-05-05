FROM golang:alpine as build
RUN apk add --no-cache --update git
ADD . /go/src/app
WORKDIR /go/src/app
RUN go get ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -tags timetzdata -o nv main.go

FROM scratch
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
