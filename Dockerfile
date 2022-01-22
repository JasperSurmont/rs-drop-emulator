#build stage
FROM golang:1.17-alpine

WORKDIR /app

ENV RS_DROP_SIMULATOR_ENV=PROD

COPY . ./

RUN go mod download

RUN go build -o /rs-drop-simulator 

CMD ["/rs-drop-simulator"]