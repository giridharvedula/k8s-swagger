FROM golang:1.22-alpine AS build
WORKDIR /src
COPY main.go .
RUN go build -o /proxy main.go

FROM alpine
COPY --from=build /proxy /proxy
ENTRYPOINT ["/proxy"]
