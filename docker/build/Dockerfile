FROM golang:1.12-alpine
RUN apk add --no-cache bash git outils-sha256 zip
WORKDIR /code
CMD ["/code/scripts/build-all-platforms.sh"]
