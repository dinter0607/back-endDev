FROM golang:1.14.2

RUN adduser --disabled-password --gecos '' api
USER api

WORKDIR /go/src/Project/app
COPY . .

RUN go install -v ./...

CMD [ "app" ]