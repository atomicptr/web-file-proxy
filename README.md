# web-file-proxy

Simple proxy for web resources with minimalistic web ui.

## How to use

Create a SHA3-256 string prefixed with "wfp_" this will be
the password with which you can log into the web ui.

```bash
$ docker run \
    -e SECRET_HASH="$YOUR_SECRET_HASH" \
    -p 8081:8081 \
    -v ./proxy.db:/data/proxy.db
    atomicptr/web-file-proxy
```

## License

MIT