FROM golang:1.4-onbuild
RUN mkdir /var/log/gtfs-api
EXPOSE 4000