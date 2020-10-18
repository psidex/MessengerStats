FROM golang:latest AS builder
WORKDIR /msbuild
COPY cmd .
COPY internal .
COPY go.mod .
COPY go.sum .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./msserver ./cmd/server/main.go

FROM alpine:latest
WORKDIR /messengerstats
COPY static static
COPY --from=builder /psmbuild/msserver .
ENV INSIDE_DOCKER "True"
EXPOSE 8080/tcp
CMD ["./msserver"]

# docker run -d --name msserver \
#     --network proxynet \
#     -v /home/user/messengerstatsdockervolume:/data \
#     psidex/messengerstats