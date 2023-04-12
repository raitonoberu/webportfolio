# build
FROM golang:alpine AS builder
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o main . 

# deploy
FROM scratch
WORKDIR /

COPY --from=builder /app/main main
ENTRYPOINT [ "/main" ]
