FROM golang AS builder
WORKDIR /go/src/github.com/nkex606/chatroom-user
COPY . .
RUN go get -d -v -t ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
ENV addr="127.0.0.1"
ENV port="8080"
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/nkex606/chatroom-user/app .
CMD [ "./app" ] 
