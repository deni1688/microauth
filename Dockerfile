FROM alpine:latest
RUN apk add libc6-compat
WORKDIR /app
COPY microauth .
EXPOSE 9876
CMD ./microauth
