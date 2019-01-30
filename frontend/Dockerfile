FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o frontend .

FROM debian
EXPOSE 8080
RUN apt-get update && \
    apt-get install -y curl dnsutils && \
    rm -rf /var/lib/apt/lists/*
COPY --from=builder /build/frontend /app/
WORKDIR /app
CMD ["./frontend"]
