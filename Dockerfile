FROM golang:latest AS builder
WORKDIR /msbuild
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./msserver ./cmd/server/main.go

FROM alpine:latest
WORKDIR /messengerstats
COPY static static
COPY views views
COPY --from=builder /msbuild/msserver .
ENV INSIDE_DOCKER "True"
EXPOSE 8080/tcp
CMD ["./msserver"]

# docker run -d --name msserver \
#     --network proxynet \
#     psidex/messengerstats