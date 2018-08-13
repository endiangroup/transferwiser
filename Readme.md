# Transferwiser

This service acts as a read-only interface for users to see the list of transactions of a transferwise account.

# Config

All configuration uses environment variables:

- `TRANSFERWISER_PORT` is the http port where the service will be served. Default `8080`
- `TRANSFERWISER_TWHOST` is the transferwise host. Default `sandbox.transferwise.tech`
- `TRANSFERWISER_TWLOGINREDIRECT` is the url to return when linking with a transferwise account.
