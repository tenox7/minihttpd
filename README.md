# upload
Receive File Handler for multipart/form-data HTTP POST / HTML Form File Uploads

## usage
compile and place in cgi-bin directory

point your browsert to http://x.x.x.x/cgi-bin/upload

or curl:

```
curl -F "filename=@<yourfile.dat>" http://x.x.x.x/cgi-bin/upload
```
