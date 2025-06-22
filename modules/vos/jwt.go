package vos

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"strings"
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

// Encodes session into a new JWT.
//
// Should never return an error, since:
//   - [hash.Write] never returns an error;
//   - [encoding/json.Marshal] for session should not return an error.
func (ses Session) NewJWT() (string, error) {
	sesJSON, err := json.Marshal(&ses)
	if err != nil {
		return "", err
	}
	sesB64 := base64.URLEncoding.EncodeToString(sesJSON)
	tokenSignature, err := SignJWT(jwtHeaderB64, sesB64, getCurrentEncryptionKey())
	if err != nil {
		return "", err
	}
	return jwtHeaderB64 + "." + sesB64 + "." + tokenSignature, nil
}

// Returns a signature/secret part of a JWT.
//
// Should never return an error, since [hash.Write] never returns an error.
func SignJWT(header, body string, encryptionKey []byte) (string, error) {
	hash := hmac.New(crypto.SHA512.New, encryptionKey)
	_, err := hash.Write([]byte(header + body))
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(hash.Sum(nil)), nil
}

// Returns if given JWT is valid or not.
//
// Should never return an error, since:
//   - [hash.Write] never returns an error;
//   - [encoding/json.Unmarshal] works with valid marshalled session thus should not return an error.
func (ses *Session) ValidateJWT(token string) (bool, error) {
	tokenParts := strings.SplitN(token, ".", 4)
	if len(tokenParts) != 3 {
		return false, nil // invalid amount of token parts
	}

	if tokenParts[0] != jwtHeaderB64 {
		return false, nil // unknown header is used
	}

	tokenStart := tokenParts[0] + "." + tokenParts[1] + "."

	tokenSignature, err := SignJWT(tokenParts[0], tokenParts[1], getCurrentEncryptionKey())
	if err != nil {
		return false, err
	}
	reconstructed := tokenStart + tokenSignature
	if reconstructed != token {
		tokenSignature, err = SignJWT(tokenParts[0], tokenParts[1], getPreviousEncryptionKey())
		if err != nil {
			return false, err
		}
		reconstructed = tokenStart + tokenSignature
		if reconstructed != token {
			return false, nil // with both keys jwt was invalid
		}
	}

	return true, json.Unmarshal([]byte(tokenParts[1]), ses)
}
