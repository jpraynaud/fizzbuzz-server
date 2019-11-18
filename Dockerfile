###############################
# STEP 1: build executable
###############################
FROM golang:1.13-alpine AS gobuilder

# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Install build base
RUN apk add build-base

# Create appuser.
RUN adduser -D -g '' appuser

# Wordkir
WORKDIR /app
COPY . .

# Show Go version
RUN go version

# Fetch dependencies.
RUN go get -d -t -v ./...

# Run tests
RUN go test -v -cover ./...

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o "/go/bin/fizzbuzz-server" cmd/server/main.go

###############################
# STEP 2: build a small image
###############################
FROM scratch

# Import the user and group files from the builder.
COPY --from=gobuilder /etc/passwd /etc/passwd

# Copy our static executable.
COPY --from=gobuilder "/go/bin/fizzbuzz-server" "/go/bin/fizzbuzz-server"

#Workdir
WORKDIR /go/bin/

# Use an unprivileged user.
USER appuser

# Run the binary.
ENTRYPOINT ["/go/bin/fizzbuzz-server"]