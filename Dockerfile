FROM ubuntu:19.04 AS client
WORKDIR /app
RUN apt-get update && \
        apt-get -y install sudo
RUN apt-get install -y nodejs npm

COPY ./pororift-client .

RUN npm install && \
        npm run test && \
        npm run build

FROM golang:1.11
WORKDIR /go-source

RUN mkdir /go-source/bin
ENV GOBIN=/go-source/bin

COPY ./src ./src
COPY ./go.mod .
COPY ./go.sum .

RUN go install src/pororift.go

COPY --from=client /app/dist dist



ENTRYPOINT ["./bin/pororift"]
