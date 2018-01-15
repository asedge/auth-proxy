# auth-proxy

This proxy allows unauthenticated access to a web resource that normally requires basic authentication by setting the basic auth header using the USERNAME/PASSWORD provided.

### Usage:
```
USERNAME=sean PASSWORD=secret PROXY_BASE=https://my.local.artifact.repo/ ./go-auth-proxy
```

### Build & Run Docker image:
```
docker build -t auth-proxy .
docker run -d --name auth_proxy -p 8989:8989 -e USERNAME=<basic_auth_username> -e PASSWORD=<basic_auth_password> -e PROXY_BASE=<base_url_to_proxy_auth_to> auth-proxy
```
