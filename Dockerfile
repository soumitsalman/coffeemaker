FROM golang:1.22.3-alpine
WORKDIR /app
COPY . .
RUN go get
RUN go build -o bin .
EXPOSE 8080
ENTRYPOINT [ "/app/bin" ]
