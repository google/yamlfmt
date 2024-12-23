FROM golang:alpine AS build
RUN apk add --no-cache git make
WORKDIR /build
COPY . .
ENV CGO_ENABLED=0
RUN make build

FROM alpine:latest
COPY --from=build /build/dist/yamlfmt /bin/yamlfmt
ENTRYPOINT ["/bin/yamlfmt"]
