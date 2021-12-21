FROM golang:alpine3.14
RUN mkdir -p /app
COPY . /app
WORKDIR /app
RUN go build -o main .
EXPOSE 4000
CMD ["./main"]