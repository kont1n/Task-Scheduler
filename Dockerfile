FROM golang:1.23

ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE ${TODO_PORT}

RUN go build -o /my_app

CMD ["/my_app"]