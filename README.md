# GoWSServe

Just testing out a go server with accounts and sessions feeding a canvas.

## Self Signing SSL For Dev

```
openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt
```

Will need to trust the cert in system/browsers.