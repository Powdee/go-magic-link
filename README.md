# Go Magic Link

![Go](./Go.png)

`go-magic-link` is a Go implementation of magic links authentication using the Echo framework with JWT utilities.

## Features

- Magic link authentication: Users receive a special link via email that logs them in when clicked.
- JWT token generation and verification for secure authentication.
- Built on the Echo web framework for robustness and scalability.

## How to run it

```shell
go run api/server.go
```

## Endpoints

`/auth/login`: Endpoint for generating magic links.
`/auth/verify`: Endpoint for verifying the magic link and authenticating the user.

## Configuration

`SecretKey`: The secret key used for JWT token signing.
`TokenExpiration`: Expiration time for the JWT token.
`EmailTemplate`: Template for the magic link email.

## License

MIT

Feel free to copy and paste this content into your README file! Let me know if there's anything else you need.
