FROM golang:latest AS go-builder
ADD . /app
WORKDIR /app/server
RUN bash generate_ent.sh
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=go-builder /main ./
RUN chmod +x ./main
EXPOSE 8080
CMD ./main
