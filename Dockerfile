FROM golang:1.12-alpine AS build

RUN mkdir -p /src
WORKDIR /src
COPY . .
RUN apk add make git
RUN make build

FROM alpine

LABEL maintainer="Steve Good (thestarwarsdev@gmail.com)"

RUN mkdir -p /opt/flare
VOLUME /opt/flare

RUN mkdir -p /app
WORKDIR /app
COPY --from=build /src/flare /app/.

# install and update ca-certificates so our app can connect to discord
RUN apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates \
    && update-ca-certificates 2>/dev/null || true

RUN addgroup -g 1000 -S flare \
    && adduser -u 1000 -S flare -G flare \
    && chown -R flare:flare /app \
    && chmod -R 777 /app

RUN chown flare:flare /opt/flare

USER flare

CMD ["/bin/sh", "-c", "./flare"]
