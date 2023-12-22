FROM alpine:latest

RUN mkdir /app

COPY mailerApp /app

COPY app.env /app

CMD [ "/app/mailerApp" ]