# unable to use scratch as there were certificate issues when authenticating with Google Secrets Manager
FROM google/cloud-sdk:282.0.0-slim

ENTRYPOINT ["gsm"]
COPY ./build/linux /usr/bin/