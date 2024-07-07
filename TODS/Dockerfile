FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./

COPY *.go ./

RUN go build -o /tods-app

EXPOSE 8080

CMD [ "/tods-app" ]
