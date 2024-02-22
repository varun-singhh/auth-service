FROM alpine:edge
ADD main /main
COPY /configs /configs
RUN apk add --no-cache tzdata ca-certificates
EXPOSE 9999
CMD ["/main"]
