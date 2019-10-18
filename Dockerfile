#FROM golang:1.13-alpine as build
#
#WORKDIR /go/src/app
#COPY . .
#ENV GO111MODULE=on
#RUN apk add --no-cache make git gcc
#RUN go mod download
#RUN CGO_ENABLED=0 go test -v ./...
#RUN go get -v ./...
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app

FROM app-base-build as build
COPY . .
RUN make build

FROM alpine

WORKDIR /app
COPY --from=build /go/src/app .

CMD ["./app"]