# build backend
FROM golang:alpine AS back
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o main . 

# build frontend
FROM node:18 AS front
RUN mkdir /app
COPY frontend/ /app/
WORKDIR /app
RUN npm install
RUN ls public
RUN npm run build

# deploy
FROM scratch
WORKDIR /

COPY --from=back /app/main main
COPY --from=front /app/build/ static/
RUN ls 
ENTRYPOINT [ "/main" ]
