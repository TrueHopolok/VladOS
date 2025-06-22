package vos

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/TrueHopolok/VladOS/modules/vos/auth"
)

// JWT header consisting of type itself and algorithm used.
// Currently [crypto/hmac] combined with [crypto.SHA512] are being used.
const JWTheaderJSON string = `
{
  "alg": "HS512",
  "typ": "JWT"
}`

// Parsed header into base64 string to be used in encryption and comparisons
var jwtHeaderB64 string = base64.URLEncoding.EncodeToString([]byte(JWTheaderJSON))

// Returns a signature/secret part of a JWT.
//
// Should never return an error, since [hash.Write] never returns an error.
func signJWT(header, body string, encryptionKey []byte) (string, error) {
	hash := hmac.New(crypto.SHA512.New, encryptionKey)
	_, err := hash.Write([]byte(header + body))
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(hash.Sum(nil)), nil
}

// Encodes session into a new JWT.
//
// Should never return an error, since:
//   - [hash.Write] never returns an error;
//   - [encoding/json.Marshal] for [github.com/TrueHopolok/VladOS/modules/vos/auth.Session] should not return an error.
func NewJWT(ses auth.Session) (string, error) {
	sesJSON, err := json.Marshal(&ses)
	if err != nil {
		return "", err
	}
	sesB64 := base64.URLEncoding.EncodeToString(sesJSON)
	tokenSignature, err := signJWT(jwtHeaderB64, sesB64, getCurrentEncryptionKey())
	if err != nil {
		return "", err
	}
	return jwtHeaderB64 + "." + sesB64 + "." + tokenSignature, nil
}

// Returns if given JWT is valid or not.
//
// Should never return an error, since:
//   - [hash.Write] never returns an error;
//   - [encoding/json.Unmarshal] works with valid marshalled [github.com/TrueHopolok/VladOS/modules/vos/auth.Session] thus should not return an error.
func ValidateJWT(token string) (auth.Session, bool, error) {
	ses := auth.Session{}
	tokenParts := strings.SplitN(token, ".", 4)
	if len(tokenParts) != 3 {
		return ses, false, nil // invalid amount of token parts
	}

	if tokenParts[0] != jwtHeaderB64 {
		return ses, false, nil // unknown header is used
	}

	tokenStart := tokenParts[0] + "." + tokenParts[1] + "."

	tokenSignature, err := signJWT(tokenParts[0], tokenParts[1], getCurrentEncryptionKey())
	if err != nil {
		return ses, false, err
	}
	reconstructed := tokenStart + tokenSignature
	if reconstructed != token {
		tokenSignature, err = signJWT(tokenParts[0], tokenParts[1], getPreviousEncryptionKey())
		if err != nil {
			return ses, false, err
		}
		reconstructed = tokenStart + tokenSignature
		if reconstructed != token {
			return ses, false, nil // with both keys jwt was invalid
		}
	}

	err = json.Unmarshal([]byte(tokenParts[1]), &ses)
	return ses, true, err
}
