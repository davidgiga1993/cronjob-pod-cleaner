FROM alpine

WORKDIR /app
COPY bin/cronjob-pod-cleaner ./
RUN chmod 555 /app/cronjob-pod-cleaner
ENTRYPOINT ["/bin/sh", "-c", "/app/cronjob-pod-cleaner"]

