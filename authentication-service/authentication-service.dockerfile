FROM alpine:latest

RUN mkdir /app

COPY authenticationApp /app

COPY app.env /app

CMD [ "/app/authenticationApp" ]