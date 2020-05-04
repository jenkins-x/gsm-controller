FROM gcr.io/distroless/base
COPY ./build/linux/gsm /gsm
ENTRYPOINT ["/gsm"]
