FROM golang:1.21-bullseye AS dev
COPY . /src/web
WORKDIR /src/web

RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y git mariadb-client

RUN go install github.com/cosmtrek/air@v1.40.4
RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/go-delve/delve/cmd/dlv@v1.20.0
RUN go install gotest.tools/gotestsum@v1.8.2
RUN go install github.com/google/wire/cmd/wire@v0.5.0
RUN go install github.com/golang/mock/mockgen@v1.6.0

RUN apt-get install bash-completion
COPY bashrc /root/.bashrc

FROM golang:1.21-bullseye AS build
COPY . /go/src/web
WORKDIR /go/src/web
RUN apt update && apt upgrade -y && apt install -y git
RUN go install github.com/google/wire/cmd/wire@v0.5.0
RUN make build_prod

FROM gcr.io/distroless/base
WORKDIR /go/src/web
COPY --from=build /go/src/web/app .

EXPOSE 3000
ENV SERVER_PORT=3000
EXPOSE 3000
CMD ["./app"]
