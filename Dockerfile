FROM golang:1.23 AS builder

WORKDIR /app

COPY . ./
RUN go mod download

RUN echo "BEFORE BUILD ---->" && ls -a

RUN go build -o mailservice .

# FINAL CONTAINER

FROM ubuntu:latest

WORKDIR /app

RUN ls -l

COPY --from=builder /app/mailservice /app/mailservice
COPY --from=builder /app/.env /app/.env

RUN ls -l
# RUN chmod +x mailservice

RUN pwd

EXPOSE 3000

CMD ["./mailservice"]

