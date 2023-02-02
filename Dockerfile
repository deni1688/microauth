FROM golang:1.19-alpine
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -tags netgo -o microauth .

FROM alpine:latest
RUN apk add libc6-compat
WORKDIR /app
COPY --from=0 /build/microauth .
EXPOSE 9876
CMD ./microauth
