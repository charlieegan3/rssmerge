FROM golang:1.10 as build

WORKDIR /go/src/github.com/charlieegan3/rssmerge

COPY . .

RUN go get -u github.com/gobuffalo/packr/packr
RUN CGO_ENABLED=0 packr build -o rssmerge


FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/charlieegan3/rssmerge/rssmerge /

CMD ["/rssmerge"]
