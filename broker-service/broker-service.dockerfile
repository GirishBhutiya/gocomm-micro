FROM alpine:latest

RUN mkdir /app

COPY brokerApp /app

COPY app.env /app

CMD [ "/app/brokerApp" ]