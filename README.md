# auth-proxy

[![Build Status](https://travis-ci.org/asedge/auth-proxy.png)](https://travis-ci.org/asedge/auth-proxy)

This proxy allows unauthenticated access to a web resource that normally requires basic authentication by setting the basic auth header using the USERNAME/PASSWORD provided.

### Usage:
```
USERNAME=sean PASSWORD=secret ./auth-proxy
```

The proxy listens on port `8989` so to use it you can just set the `http_proxy` environment variable.

```
export http_proxy=http://localhost:8989
curl http://site.to.proxy/
```

### Build & Run Docker image:
```
docker build -t auth-proxy .
docker run -d --name auth_proxy -p 8989:8989 -e USERNAME=<basic_auth_username> -e PASSWORD=<basic_auth_password> auth-proxy
```
