FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY aip_food_lookup /usr/local/bin/aip_food_lookup

EXPOSE 8080

CMD ["/usr/local/bin/aip_food_lookup"]
