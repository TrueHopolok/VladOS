package vos

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// JWT header consisting of type itself and algorithm used.
// Currently [crypto/hmac] combined with [crypto.SHA512] are being used.
const JWTheader = `
{
  "alg": "HS512",
  "typ": "JWT"
}`

// Parsed header into base64 string to be used in encryption and comparisons
var jwtHeaderB64 = base64.URLEncoding.EncodeToString([]byte(JWTheader))

// Returns a secret part of a JWT in base64 by signing the message with given encryption key.
// Use encryption key from key chain pair.
//
// Should never return an error, since [hash.Write] never returns an error.
func signJWT(message, encryptionKey []byte) (string, error) {
	hash := hmac.New(sha512.New, encryptionKey)
	_, err := hash.Write(message)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(hash.Sum(nil)), nil
}

// Encodes session into a new JWT.
//
// Should never return an error, since:
//   - [hash.Write] never returns an error;
//   - [encoding/json.Marshal] for [Session] should not return an error.
func NewJWT(ses Session) (string, error) {
	sesJSON, err := json.Marshal(&ses)
	if err != nil {
		return "", err
	}
	sesB64 := base64.URLEncoding.EncodeToString(sesJSON)
	tokenSignature, err := signJWT([]byte(jwtHeaderB64+sesB64), getCurrentEncryptionKey())
	if err != nil {
		return "", err
	}
	return jwtHeaderB64 + "." + sesB64 + "." + tokenSignature, nil
}

// Returns if given JWT is valid or not.
// In case of being valid, parses token and returns it as [Session] struct.
//
// Should never return an error, since:
//   - [hash.Write] never returns an error;
//   - [encoding/json.Unmarshal] works with valid marshalled [Session] thus should not return an error;
//   - [encoding/base64.URLEncoding.DecodeString] works with valid encoded [Session] thus should not return an error.
func ValidateJWT(token string) (Session, bool, error) {
	ses := Session{}
	tokenParts := strings.SplitN(token, ".", 4)
	if len(tokenParts) != 3 {
		return ses, false, nil // invalid amount of token parts
	}

	if tokenParts[0] != jwtHeaderB64 {
		return ses, false, nil // unknown header is used
	}

	tokenStart := tokenParts[0] + "." + tokenParts[1] + "."

	tokenSignature, err := signJWT([]byte(tokenParts[0]+tokenParts[1]), getCurrentEncryptionKey())
	if err != nil {
		return ses, false, fmt.Errorf("signing error: %w", err)
	}
	reconstructed := tokenStart + tokenSignature
	if reconstructed != token {
		tokenSignature, err = signJWT([]byte(tokenParts[0]+tokenParts[1]), getPreviousEncryptionKey())
		if err != nil {
			return ses, false, fmt.Errorf("signing error: %w", err)
		}
		reconstructed = tokenStart + tokenSignature
		if reconstructed != token {
			return ses, false, nil // with both keys jwt was invalid
		}
	}

	sessionJSON, err := base64.URLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return ses, true, fmt.Errorf("base64 error: %w", err)
	}
	if err = json.Unmarshal(sessionJSON, &ses); err != nil {
		return ses, true, fmt.Errorf("json error: %w", err)
	}
	return ses, true, nil
}
