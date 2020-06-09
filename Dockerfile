FROM golang

RUN mkdir -p /home/go/app
WORKDIR /home/go/app
COPY . /home/go/app
RUN go build main.go
ENTRYPOINT ./main

EXPOSE 9765

