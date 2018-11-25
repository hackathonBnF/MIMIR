FROM golang:1.10 AS builder

# Download and install the latest release of dep
# ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
# RUN chmod +x /usr/bin/dep

# Copy the code from the host and compile it
COPY . ./

RUN mkdir res
RUN mkdir -p /go/src/gitlab.com/GT-RIMi/RIMo

COPY settings.json res
ADD static res/static
ADD templates res/templates

RUN cp -r goproviders media query settings main.go /go/src/gitlab.com/GT-RIMi/RIMo
RUN cd /go/src/gitlab.com/GT-RIMi/RIMo

RUN go get -v -d ./...
RUN cd /go

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

# We use Alpine for it's ca-certificates needed by http lib
FROM alpine:latest
RUN apk add --no-cache ca-certificates apache2-utils
COPY --from=builder /app ./
COPY --from=builder /go/res ./


ENTRYPOINT ["./app"]
