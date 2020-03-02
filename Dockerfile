FROM scratch
EXPOSE 8080
ENTRYPOINT ["/gsm-controller"]
COPY ./build/linux /