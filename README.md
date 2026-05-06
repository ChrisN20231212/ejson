# ejson

轻量 JSON 在线解析服务，提供浏览器工具页和可脚本调用的 HTTP API。

## 运行

```bash
go run . -listen 0.0.0.0:8080
```

打开 `http://<你的服务器IP>:8080/`。

如果设置了 `PORT`，服务会优先监听 `0.0.0.0:$PORT`：

```bash
PORT=9090 go run .
```

## Docker

构建镜像：

```bash
docker build -t ejson .
```

运行容器：

```bash
docker run --rm -p 8080:8080 -e PORT=8080 ejson
```

打开 `http://<你的服务器IP>:8080/`。

## Docker Compose

直接使用已发布镜像启动：

```bash
docker compose up -d
```

如果需要先拉取最新镜像：

```bash
docker compose pull
docker compose up -d
```

停止服务：

```bash
docker compose down
```

## API

所有接口均为 `POST`，请求体格式：

```json
{"input":"你的内容"}
```

接口列表：

- `/api/validate`
- `/api/minify`
- `/api/escape`
- `/api/unescape`
- `/api/unicode/encode`
- `/api/unicode/decode`
