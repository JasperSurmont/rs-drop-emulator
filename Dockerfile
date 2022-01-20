#build stage
FROM golang:1.17-alpine

WORKDIR /app


ENV RS_DROP_EMULATOR_ENV=PROD

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /rs-drop-emulator 

CMD ["/rs-drop-emulator"]