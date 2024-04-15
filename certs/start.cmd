mkdir .\ca

openssl genrsa -out .\ca\ca.key 2048

openssl req -x509 -new -key .\ca\ca.key -days 36500 -out .\ca\ca.crt -subj "/C=CN/ST=GuangDong/L=ShenZhen/O=EMQX/CN=server"


mkdir .\server
openssl genrsa -out .\server\server.key 2048

openssl req -new -key .\server\server.key -config openssl.cnf -out .\server\server.csr

openssl x509 -req -in .\server\server.csr -CA .\ca\ca.crt -CAkey .\ca\ca.key -CAcreateserial -out .\server\server.crt -days 36500 -sha256 -extensions v3_req -extfile openssl.cnf


mkdir .\client
openssl genrsa -out .\client\client.key 2048

openssl req -new -key .\client\client.key -out .\client\client.csr -subj "/C=CN/ST=GuangDong/L=ShenZhen/O=EMQX/CN=client"

openssl x509 -req -days 36500 -in .\client\client.csr -CA .\ca\ca.crt -CAkey .\ca\ca.key -CAcreateserial -out .\client\client.crt
