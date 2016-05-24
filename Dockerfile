# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM iron/go:dev

# Copy the local package files to the container's workspace.
COPY . /go/src/github.com/deciphernow/synonyms

# Build the synonyms command inside the container.
RUN go install github.com/deciphernow/synonyms

# Run the synonyms command by default when the container starts.
ENTRYPOINT ["/go/bin/synonyms"]

# Document that the service listens on port 8080.
EXPOSE 8080
