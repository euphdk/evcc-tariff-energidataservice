FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

ADD energidataservice-dk-chain.pem /etc/ssl/certs/

COPY . .
RUN go build -v -o /usr/local/bin/tariff main.go

CMD ["tariff"]
