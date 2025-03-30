FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy


EXPOSE 8080

#DEV

RUN go install -mod=mod github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -directory=/app -polling=true -build="go build -o main ./cmd" -command="./main"

#PROD
#CMD ["go", "run", "cmd/main.go"]
