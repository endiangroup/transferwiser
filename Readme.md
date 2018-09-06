# Transferwiser

This service acts as a read-only interface for users to see the list of transactions of a transferwise account.

# Config

All configuration uses environment variables:

- `TRANSFERWISER_ENV` is the environment to use. Set it to `production` to use the production CA. Default `dev`.
- `TRANSFERWISER_PORT` is the http port where the service will be served. Default `8080`
- `TRANSFERWISER_LETSENCRYPTPORT` is the http port where letsencrypt will run the challenges. Default `8081`
- `TRANSFERWISER_TWHOST` is the transferwise host to use. Default `api.sandbox.transferwise.tech`
- `TRANSFERWISER_TWAPITOKEN` is the transferwise API Token. Required.
- `TRANSFERWISER_TWPROFILE` is the profile ID of the account in transferwise. Required.


# Mutual tls authentication reference

https://blog.codeship.com/how-to-set-up-mutual-tls-authentication/
https://medium.com/@itseranga/tls-mutual-authentication-with-golang-and-nginx-937f0da22a0e
https://fale.io/blog/2017/06/05/create-a-pki-in-golang/
https://kb.op5.com/pages/viewpage.action?pageId=19073746#sthash.9kVrUqUY.dpbs
https://stackoverflow.com/questions/42643048/signing-certificate-request-with-certificate-authority
https://gist.github.com/ncw/9253562
https://gist.github.com/sdorra/1c95de8cb80da31610d2ad767cd6f251


