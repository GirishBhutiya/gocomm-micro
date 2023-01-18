#build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY DBApp /app

CMD [ "/app/DBApp" ]