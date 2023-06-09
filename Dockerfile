FROM golang:1.20

WORKDIR /usr/src/app
# COPY go.mod go.sum ./
# RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app .

FROM alpine:latest

RUN apk add --no-cache poppler poppler-utils
COPY --from=0 /usr/src/app/app ./

EXPOSE 3333
CMD ["./app"]
