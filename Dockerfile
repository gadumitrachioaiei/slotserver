FROM golang:alpine3.15
WORKDIR /app
COPY . .
RUN go build -o slotserver

FROM alpine
COPY --from=0 /app/slotserver .
ENTRYPOINT ["./slotserver"]
CMD ["--http", ":8080"]

