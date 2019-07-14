FROM alpine
RUN apk add --no-cache ca-certificates
ADD ./imdb imdb
EXPOSE 8080
ENTRYPOINT ["/imdb"]