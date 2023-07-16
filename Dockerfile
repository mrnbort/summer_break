FROM golang:1.20-alpine as build

ADD . /build
WORKDIR /build

RUN go build -o /build/summer_break -ldflags "-s -w"


FROM alpine:3.17

COPY --from=build /build/summer_break /srv/summer_break
COPY ./testdata /srv/testdata
RUN chmod +x /srv/summer_break

WORKDIR /srv
EXPOSE 8080
CMD ["/srv/summer_break"]