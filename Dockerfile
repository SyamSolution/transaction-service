FROM golang:latest
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main ./cmd/http/api.go
EXPOSE 4000
CMD [ "./main" ]