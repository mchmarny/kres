# BUILD STAGE
FROM golang:latest as build

# copy
WORKDIR /go/src/github.com/mchmarny/kres/
COPY . /src/

# dependancies
WORKDIR /src/
ENV GO111MODULE=on
RUN go mod download

# build
WORKDIR /src/cmd/service/
RUN CGO_ENABLED=0 go build -o /kres


# RUN STAGE
FROM alpine as release
RUN apk add --no-cache ca-certificates

# app executable
COPY --from=build /kres /app/

# run
WORKDIR /app/
ENTRYPOINT ["./kres"]