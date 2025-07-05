/*
# VladOS Security

VladOS Security (shortly VOS) is a package used for all security related stuff inside VladOS the project.
Below are all features this package provides.

# Key Chain

2 cryptografy keys used for JWT and thus for session encryption.
Automatic keys switch after [EncryptionKeysSwitchingTime] minutes ensures keys are hard to crack.

# Json Web Token functional

Fully working JWT with a nice functions to create and validate tokens.
Used in combination with sessions and in [auth.AuthMiddleware] functional.

For more details see [JWTheader], [NewJWT], [ValidateJWT].

# Password and Salt logic

Password is generated with [golang.org/x/crypto/argon2] package.
Salt is generated via [crypto/rand.Read] to ensure cryptografic salt.

For more details see [GenerateSalt], [NewPSH], [ValidatePSH].

# Authentication

Implements [AuthMiddlware] to use in http server that handles whole authefication with cookie handlers.
Saves result in the [net/http.Request.Context] for further use in the request handlers.
*/
package vos

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md
