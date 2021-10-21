FROM golang:1.17

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY *.go ./

RUN go mod download
RUN go build -o dbmigrate
CMD [ "./dbmigrate", "-truncate" ]