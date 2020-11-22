package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// Returns the JWT key.
func getJwtKey() []byte {
	return []byte(viper.GetString("jwtKey"))
}

// Creates a new JWT that contains the provided Claim signed by the key from the configuration.
func CreateJwt(claim Claim) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, createClaim(claim))
	tokenStr, err := token.SignedString(getJwtKey())
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// Converts a token string into a jwt.Token. Verifies that a JWT is valid and returns a jwt.Token when this is the case.
func jwtToToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Verify that HMAC is used
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected singing method: %v", token.Header["alg"])
		}
		return getJwtKey(), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Parses a JWT by first converting the supplied token string into a jwt.Token and then checking if the token is valid.
func ParseJwt(tokenStr string) (jwt.Claims, error) {
	token, err := jwtToToken(tokenStr)
	if err != nil {
		return nil, err
	}

	// TODO: Do some type checking here. If we add more JWTs in the future, a JWT with multiple claims can be valid for
	//       another JWT which is not a good case since it can lead to security issues. We can forget this until we have
	//       more claims.
	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}
	return token.Claims, nil
}
