FROM alpine:3.13

COPY drb /

ENTRYPOINT ["/drb"]
