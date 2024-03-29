FROM golang:alpine as builder
ENV GO111MODULE="" \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN apk add --no-cache git
RUN go mod download
COPY . .
RUN go build -o main .
WORKDIR /dist
RUN cp /build/main .


FROM golang:alpine as deploy
COPY --from=builder /dist .
RUN apk update && apk add wget && apk add bash && apk add zip
ENV BASE_PATH="/opt/data"
VOLUME [ "/opt/data" ]
EXPOSE 8002
CMD ["./main"]