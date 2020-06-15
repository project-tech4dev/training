# FROM golang

# RUN mkdir -p /home/go/app
# WORKDIR /home/go/app
# COPY . /home/go/app
# RUN go build main.go
# ENTRYPOINT ./main

FROM node:10.16.0
RUN mkdir -p /home/node/app
WORKDIR /home/node/app
COPY nodeserver /home/node/app
RUN npm install


EXPOSE 9765

