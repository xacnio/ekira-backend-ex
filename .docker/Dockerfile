FROM golang:1.19.6-alpine

RUN apk add --no-cache --virtual .deps \
    						pkgconfig \
    						gcc \
                musl-dev \
                zlib \
                imagemagick \
                imagemagick-libs \
                imagemagick-dev \
    						make

WORKDIR /app

COPY go.mod go.mod

COPY go.sum go.sum

RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12

COPY . .

RUN make swag

RUN make build

RUN mv build/ekira-backend .

EXPOSE 7560

CMD ["./ekira-backend"]
