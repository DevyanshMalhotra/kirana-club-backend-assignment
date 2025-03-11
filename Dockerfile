# build section
FROM golang:1.24 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o kirana-club .

# run section
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/kirana-club .
COPY StoreMasterAssignment.csv .
EXPOSE 8080
CMD ["./kirana-club"]
