FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . ./

EXPOSE 8080

CMD ["go", "run", "cmd/main.go"]

#CMD CompileDaemon --build="go build -o main ." --command="./main"
