FROM golang:1.15.3-alpine3.12 as build
WORKDIR /build
ENV CGO_ENABLED=0

# download all imports and prebuild in cache
COPY go.mod go.sum ./
COPY ./internal/imports ./internal/imports
RUN go build ./internal/imports

# no cache
COPY . .
RUN go build .

FROM scratch
WORKDIR /app
COPY --from=build /build/k8gb-discovery /app/k8gb-discovery
ENTRYPOINT ["./k8gb-discovery"]

