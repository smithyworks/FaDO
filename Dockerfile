# build react client
FROM node:16

WORKDIR /client

COPY ./client/package*.json ./

RUN npm ci --only=production

COPY ./client/ ./

RUN npm run build

# build go server
FROM golang:1.17-alpine

WORKDIR /server

COPY ./server/go.mod ./
COPY ./server/go.sum ./

RUN go mod download

COPY ./server/ ./

RUN mkdir build

COPY --from=0 /client/build/ ./build/

RUN go build -o /FaDO

RUN go install github.com/minio/mc@latest
RUN mv $GOPATH/bin/mc /mc

# Run the app
FROM alpine:3

WORKDIR /app

COPY --from=1 /FaDO ./

COPY --from=1 /mc /bin/

CMD [ "/app/FaDO" ]
