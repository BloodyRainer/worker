FROM golang:1.11.0-alpine3.8 AS build-env
WORKDIR /fastworker/

COPY . /fastworker/

# go module needs alpine-sdk and git
RUN apk update \
    && apk add alpine-sdk \
    && apk add git

RUN go build -o /fw

FROM alpine
RUN apk update \
    && apk add ca-certificates \
    && apk add curl \
    && rm -rf /var/cache/apk/*

COPY --from=build-env /fw /

EXPOSE 8080

CMD /fw