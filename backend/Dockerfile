FROM golang


RUN go get github.com/gorilla/mux \
           github.com/go-redis/redis \
           github.com/rs/cors

WORKDIR /go/src/app

COPY ./app .

RUN go install ./*
