FROM golang:1.22.3

# install packages
WORKDIR /app

# Build the go code
COPY . .
RUN go get
RUN go build -o indexer .

# Set entrypoint.
ENTRYPOINT ["/app/indexer"]

