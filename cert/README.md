## Enabling TLS layer for gRPC input:

* In order to make TLS work, we need to provide two files for server:

    1. /cert/serv.cert
    1. /cert/priv.key

* And for client:

    1. /cert/serv.cert

## For development, there is created self-signed certificate for server, dont do it for production environment cause it is security risk
