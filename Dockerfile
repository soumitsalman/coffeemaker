FROM golang:1.22.3-alpine
WORKDIR /app
COPY . .

ENV GIN_MODE release
RUN go get
RUN go build -o bin .

ENTRYPOINT [ "/app/bin" ]
