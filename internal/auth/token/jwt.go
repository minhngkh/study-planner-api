package token

import (
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

// Contains the serialized token and the token itself before serialization
type JwtToken struct {
	Value string
	JwtRegisteredClaims
}

// TODO: Wrap this to abstract from the library
type JwtRegisteredClaims = jwt.Claims
type JwtPrivateClaims = interface{}

// Signs a JWT token in JWS format
func SignJwtToken(
	key []byte,
	rc JwtRegisteredClaims,
	pcs ...JwtPrivateClaims,
) (JwtToken, error) {
	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.HS256,
			Key:       key,
		},
		nil,
	)
	if err != nil {
		return JwtToken{}, err
	}

	builder := jwt.Signed(signer).
		Claims(rc)
	for _, pc := range pcs {
		builder = builder.Claims(pc)
	}

	signedToken, err := builder.Serialize()
	if err != nil {
		return JwtToken{}, err
	}

	return JwtToken{Value: signedToken, JwtRegisteredClaims: rc}, nil
}

// Validates a signed JWT token in JWS format.
// "claims" will be populated and used for validation.
// Returns the registered claims.
func ValidateSignedJwtToken(
	key []byte,
	token string,
	claims ...JwtPrivateClaims,
) (JwtRegisteredClaims, error) {
	// Parse the token
	parsed, err := jwt.ParseSigned(
		token,
		[]jose.SignatureAlgorithm{jose.HS256},
	)
	if err != nil {
		return JwtRegisteredClaims{}, err
	}

	var rc JwtRegisteredClaims

	// Parse the registered claims
	err = parsed.Claims(key, &rc)
	if err != nil {
		return JwtRegisteredClaims{}, err
	}

	// Validate the token
	err = rc.Validate(jwt.Expected{})
	if err != nil {
		return JwtRegisteredClaims{}, err
	}

	// Parse the private claims
	err = parsed.Claims(key, claims...)
	if err != nil {
		return JwtRegisteredClaims{}, err
	}

	return rc, nil
}

// Encrypts a JWT token in JWE format
func EncryptJwtToken(
	key []byte,
	rc JwtRegisteredClaims,
	pcs ...JwtPrivateClaims,
) (JwtToken, error) {
	encrypter, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.DIRECT,
			Key:       key,
		},
		nil,
	)

	if err != nil {
		return JwtToken{}, err
	}

	builder := jwt.Encrypted(encrypter).
		Claims(rc)
	for _, pc := range pcs {
		builder = builder.Claims(pc)
	}

	serializedToken, err := builder.Serialize()
	if err != nil {
		return JwtToken{}, err
	}

	return JwtToken{Value: serializedToken, JwtRegisteredClaims: rc}, nil
}

// Validates an encrypted JWT token in JWE format.
// "claims" will be populated and used for validation.
// Returns the registered claims.
func ValidateEncryptedJwtToken(
	key []byte,
	token string,
	claims ...JwtPrivateClaims,
) (JwtRegisteredClaims, error) {
	// Decrypt & parse the token
	parsedToken, err := jwt.ParseEncrypted(
		token,
		[]jose.KeyAlgorithm{jose.DIRECT},
		[]jose.ContentEncryption{jose.A256GCM},
	)
	if err != nil {
		return JwtRegisteredClaims{}, err
	}

	// Parse the registered claims
	var rc JwtRegisteredClaims
	err = parsedToken.Claims(key, &rc)
	if err != nil {
		return JwtRegisteredClaims{}, err
	}

	// Validate the token
	err = rc.Validate(jwt.Expected{})
	if err != nil {
		return JwtRegisteredClaims{}, err
	}

	// Parse the private claims
	err = parsedToken.Claims(key, claims...)
	if err != nil {
		return JwtRegisteredClaims{}, err
	}

	return rc, nil
}
