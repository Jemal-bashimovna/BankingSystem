FROM golang:1.22-alpine

COPY . /appGo

WORKDIR /appGo

EXPOSE 8080

RUN go build main.go

CMD [ "./main" ]