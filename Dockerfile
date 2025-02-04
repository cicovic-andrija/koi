FROM golang:1.22 AS build-stage
WORKDIR /src
COPY go.mod ./
# COPY go.sum ./
RUN go mod download
COPY main.go ./
ADD server ./server
ADD set ./set
# Make sure to statically link the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /bin/koipond ./main.go

FROM scratch
WORKDIR /srv
COPY --from=build-stage /bin/koipond ./koipond
ADD data ./data
CMD ["/srv/koipond"]
