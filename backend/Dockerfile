FROM golang:1.17.2-buster

WORKDIR /bin

RUN wget https://dl.min.io/client/mc/release/linux-amd64/mc
RUN chmod +x mc
RUN export PATH=/bin:$PATH

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY . ./

RUN go build -o /FaDO

COPY ./build /build

EXPOSE 9090

CMD [ "/FaDO" ]
