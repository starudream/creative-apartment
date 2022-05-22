# Creative Apartment （城投宽庭）

![Golang](https://img.shields.io/github/workflow/status/starudream/creative-apartment/Golang/master?label=Golang&style=for-the-badge)
![Docker](https://img.shields.io/github/workflow/status/starudream/creative-apartment/Docker/master?label=Docker&style=for-the-badge)
![Release](https://img.shields.io/github/v/release/starudream/creative-apartment?include_prereleases&style=for-the-badge)
![License](https://img.shields.io/badge/License-Apache%20License%202.0-blue?style=for-the-badge)

每日运行收集电费、水费相关数据。

## Usage

```shell
Usage:
  creative-apartment [flags]

Flags:
      --debug           (env: SCA_DEBUG) show debug information
  -h, --help            help for creative-apartment
      --path string     (env: SCA_PATH) configuration file path
      --port int        (env: SCA_PORT) http server port (default 8089)
      --secret string   (env: SCA_SECRET) http server login secret
      --startup         (env: SCA_STARTUP) execute jobs at startup
  -v, --version         version for creative-apartment
```

### Docker

![Version](https://img.shields.io/docker/v/starudream/creative-apartment?style=for-the-badge)
![Size](https://img.shields.io/docker/image-size/starudream/creative-apartment/latest?style=for-the-badge)
![Pull](https://img.shields.io/docker/pulls/starudream/creative-apartment?style=for-the-badge)

```bash
docker pull starudream/creative-apartment
```

```bash
mkdir -p /opt/docker/creative-apartment

docker run -d \
    --name creative-apartment \
    --restart always \
    -p 8089:8089 \
    -e SCA_DEBUG=true \
    -e SCA_PATH=/data/creative-apartment.yaml \
    -v /opt/docker/creative-apartment:/data \
    starudream/creative-apartment
```

## Configuration

### Before

如何获取登录用户的 `token`，首先可以使用登录接口，使用短信验证码、手机号、密码进行登录。

但是宽庭是单点登录，所以无法同时运行该程序和手机客户端。

所以最好采用的是抓包，使用 `BlackBox` `TrustMeAlready` `VNET` 工具，抓包获取 `access_token`。

可以参考文章 [无 Root 抓包 HTTPS 请求](https://blog.starudream.cn/2022/05/09/android-packet-capture-without-root/)

### Path

The configuration file is read sequentially from the following paths:

- `${EXECUTED_PATH}/creative-apartment.yaml`
- `${HOME}/creative-apartment.yaml`
- `${HOME}/.config/starudream/creative-apartment.yaml`
- `${SCA_PATH}`

### Environment Variables

Each variable is preceded by a `SCA_` prefix

| Variable        | Type   | Default | Description              |
|-----------------|--------|---------|--------------------------|
| LOG_LEVEL       | STRING | INFO    | log level                |
| DEBUG           | BOOL   | FALSE   | show debug information   |
| PATH            | STRING | -       | configuration file path  |
| PORT            | INT    | 8089    | http server port         |
| SECRET          | STRING | -       | http server login secret |
| STARTUP         | BOOL   | FALSE   | execute jobs at startup  |
| DINGTALK_TOKEN  | STRING | -       | dingtalk robot token     |
| DINGTALK_SECRET | STRING | -       | dingtalk robot secret    |

- `LOG_LEVEL`: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`

### Example

```yaml
customers:
  - phone: "${PHONE}"
    token: "${ACCESS_TOKEN}"
dingtalk:
  secret: "${DINGTALK_SECRET}"
  token: "${DINGTALK_TOKEN}"
secret: "${SECRET}"
```

## Screenshot

![dingtalk](./docs/dingtalk.jpg)

## License

[Apache License 2.0](./LICENSE)
