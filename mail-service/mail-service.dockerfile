FROM alpine:latest
RUN mkdir /app
RUN mkdir /templates

COPY mailerApp /app
COPY ./templates /templates

# Run the server executable
CMD [ "/app/mailerApp" ]