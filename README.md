# Mini HTTPD with Form Upload

A mini web server with added form based file upload. Serve single directory for in/out.

This time uses Caddy server and packaged as a Docker container.

To run:

```sh
docker run -d -v /some/directory:/www -p 80:80 tenox7/minihttpd:latest
```
