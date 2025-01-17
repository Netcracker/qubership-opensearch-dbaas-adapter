FROM golang:1.22.5-alpine3.20 as builder

ENV GOPRIVATE=github.com/Netcracker \
    GOSUMDB=off \
    CGO_ENABLED=0 \
    GO111MODULE=on

ARG GH_ACCESS_TOKEN

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

RUN apk add --no-cache git

RUN git config --global url."https://${GH_ACCESS_TOKEN}@github.com/".insteadOf "https://github.com/"

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY api api
COPY backup backup
COPY basic basic
COPY client client
COPY cluster cluster
COPY cmd cmd
COPY common common
COPY health health
COPY physical physical
COPY server server

RUN go mod tidy

RUN go test -v ./...

#RUN go test -v ./adapter...
RUN GOOS=linux GOARCH=amd64 go build -a -o main cmd/main.go

FROM scratch
COPY --from=builder /workspace/main /bin/sh

USER 1000:0

CMD ["/bin/sh", "-mode", "service"]
