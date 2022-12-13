## generate private key
openssl genrsa -out rsa.private 4096

## generate public key
openssl rsa -in rsa.private -out rsa.public