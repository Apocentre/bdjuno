FROM --platform=linux ubuntu:22.04
ARG BUILDARCH

ENV LOCAL=/usr/local

COPY build/bdjuno-linux-${BUILDARCH} ${LOCAL}/bin/bdjuno

ENTRYPOINT [ "bdjuno" ]
