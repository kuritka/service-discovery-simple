FROM golang:1.16.1-alpine3.13 as build
WORKDIR /build
ENV CGO_ENABLED=0
ENV USER=disco
ENV UID=12345
ENV GID=23456

# download all imports and prebuild in cache
COPY go.mod go.sum ./
COPY ./internal/imports ./internal/imports
RUN go build ./internal/imports
RUN addgroup -g ${GID} ${USER} && \
    adduser -D -u ${UID} -G ${USER} ${USER}

# no cache
COPY . .
RUN go build .

FROM scratch
WORKDIR /app
COPY --from=build /build/k8gb-discovery /app/k8gb-discovery
USER ${USER}

ENTRYPOINT ["./k8gb-discovery","listen"]
