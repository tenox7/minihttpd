# mini http

simple http server with form based file upload

## usage

### server

```sh
./minihttpd -port :8080 -path /tmp
```

### client

```sh
$ curl -F "name=@<yourfile.dat>" http://.../upload
```
