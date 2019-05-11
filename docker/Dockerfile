# execute from project root

# build inside an own docker
FROM golang:alpine AS build

# without musl-dev gcc it wont work in arm64v8, even if it does on amd64
RUN apk update \
    && apk --no-cache add make git gcc musl-dev

COPY ./ /build/go-dslp

RUN cd /build/go-dslp \
    && go build -x -o go-dslp

###################################
# Set the base image
FROM alpine:latest

# informations
LABEL author="Iceflower S" \
      maintainer="Iceflower S" \
      email="iceflower@gmx.de" \
      version="$VERSION" \
      description="GO DSLP Protocoll implementation"

# add server user and set permissions
RUN adduser -S -D iceflower \
    && mkdir /go-dslp \
    && chown iceflower: -R /go-dslp

WORKDIR /go-dslp
COPY --from=build /build/go-dslp/build/go-dslp .

USER iceflower
ENTRYPOINT ["./go-dslp"]
