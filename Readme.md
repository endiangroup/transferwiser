# Transferwiser

This service acts as a read-only interface for users to see the list of transactions of a transferwise account.

# Config

All configuration uses environment variables:

- `TRANSFERWISER_PORT` is the http port where the service will be served. Default `8080`
- `TRANSFERWISER_TWHOST` is the transferwise host. Default `sandbox.transferwise.tech`
- `TRANSFERWISER_TWLOGINREDIRECT` is the url to return when linking with a transferwise account.
- `TRANSFERWISER_TWCLIENTID` is the client_id sent to transferwise authentication requests. Default `endiangroup/transferwiser`
- `TRANSFERWISER_REDISADDR` the addr to find redis. Default `localhost:6379`


# Mutual tls authentication reference

https://blog.codeship.com/how-to-set-up-mutual-tls-authentication/
https://medium.com/@itseranga/tls-mutual-authentication-with-golang-and-nginx-937f0da22a0e
https://fale.io/blog/2017/06/05/create-a-pki-in-golang/
https://kb.op5.com/pages/viewpage.action?pageId=19073746#sthash.9kVrUqUY.dpbs
https://stackoverflow.com/questions/42643048/signing-certificate-request-with-certificate-authority
https://gist.github.com/ncw/9253562
https://gist.github.com/sdorra/1c95de8cb80da31610d2ad767cd6f251


