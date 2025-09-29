FROM golang:1.24-alpine

WORKDIR /app

COPY go.work ./
COPY publisher/go.mod publisher/
COPY consumer/go.mod consumer/
COPY nf-server/go.mod nf-server/


RUN go work sync
RUN go mod download

COPY publisher ./publisher
COPY nf-server ./nf-server
COPY consumer ./consumer

RUN go build -o /consumer ./consumer

CMD ["/consumer"]