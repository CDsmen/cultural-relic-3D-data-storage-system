FROM golang:1.21

ENV GOPROXY='https://proxy.byted.org|https://goproxy.cn|direct'

COPY ./exe/gltfpack /usr/local/bin/
RUN chmod +x /usr/local/bin/gltfpack

RUN mkdir -p /app/save_3d_file

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o artifact_svr

EXPOSE 8080

CMD ["./artifact_svr"]