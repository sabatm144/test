FROM alpine:3.6
MAINTAINER mahendrank@benseron.com

RUN apk update && apk upgrade \
  && rm -rf /var/cache/apk/*

COPY public ./public/
COPY linga-syncserver ./linga-syncserver
RUN chmod +x ./linga-syncserver

EXPOSE 9096
CMD ["./linga-syncserver"]
