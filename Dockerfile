FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum .env ./
RUN go mod download
COPY *.go ./

RUN go build -o mailservice

RUN ls -a

EXPOSE 3000


CMD ["./mailservice"]
