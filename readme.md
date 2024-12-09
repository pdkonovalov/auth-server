# Jwt authorization service.

This is a simple rest api with two endpoints. The first endpoint accepts the guid as a parameter and issues a pair of access and update tokens. We can get as many pairs as we need.

```bash
curl http://localhost/api/v1/jwt/new?guid=d72a56aa-95b0-409e-b945-ee176a0a8e5b
```

Guid must be valid uuid ([RFC 9562](https://datatracker.ietf.org/doc/html/rfc9562)). 

```json
{
    "AccessToken":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpcCI6IjE3Mi4xOS4wLjE6NDk0MTIiLCJqdGkiOiI1YWQ3ZTdiMS0yN2ZhLTQ4ZjgtOTgzNC1mN2QzNDg3NDU0NWUifQ.B59YVqgDvrZUF2DHCJeP9Px_tOX8RZuwmnk036CUYcL__VnZwfD8QPmCRf0Nnbq9duSFfz2iSyOp_nsSbZfmqg",
    "RefreshToken":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiI1YWQ3ZTdiMS0yN2ZhLTQ4ZjgtOTgzNC1mN2QzNDg3NDU0NWUifQ.v-Rsz4JkFcws_eG8eDYwhdVxnf3LVkzUf_RSdoE6e2IdJfcbqxahac4-Y2WbK3ymZGIxejSN2aleZ7koD0zh_g"
}
```

The second endpoint refresh pair of tokens.

```bash
curl --request POST http://localhost/api/v1/jwt/refresh \
--header 'Content-Type: application/json' \
--data '{
    "AccessToken":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpcCI6IjE3Mi4xOS4wLjE6NDk0MTIiLCJqdGkiOiI1YWQ3ZTdiMS0yN2ZhLTQ4ZjgtOTgzNC1mN2QzNDg3NDU0NWUifQ.B59YVqgDvrZUF2DHCJeP9Px_tOX8RZuwmnk036CUYcL__VnZwfD8QPmCRf0Nnbq9duSFfz2iSyOp_nsSbZfmqg",
    "RefreshToken":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiI1YWQ3ZTdiMS0yN2ZhLTQ4ZjgtOTgzNC1mN2QzNDg3NDU0NWUifQ.v-Rsz4JkFcws_eG8eDYwhdVxnf3LVkzUf_RSdoE6e2IdJfcbqxahac4-Y2WbK3ymZGIxejSN2aleZ7koD0zh_g"
}'
```

```json
{
    "AccessToken":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpcCI6IjE3Mi4xOS4wLjE6NTI2NjIiLCJqdGkiOiJiMWNkN2U0NC01NGJlLTRmOTUtYWNmMC1jM2U1MjMwZDYzZmEifQ.YnIAdCvM0ayKgmDLxXFlSZpcauET7K_uWozMEgPrni03v97io50NZOuQeKf5c-jxE-b1P5E16v0f35ITtF5V-w",
    "RefreshToken":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJiMWNkN2U0NC01NGJlLTRmOTUtYWNmMC1jM2U1MjMwZDYzZmEifQ.YI7_2w8M8L7ty45cnjU85HtQjx6Km_Vxkl75RJ6_y3l1ake5hJHNVy1Z5LYsuVxtwpkvvPVFxzbxQSXyyCiNNQ"
}
```

The refresh operation for a single token pair can only be applied once. To refresh the tokens, we must use the refreshed pair.
