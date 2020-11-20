package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Contains options for the JWT creation.
type ClaimOptions struct {
	Claim jwt.MapClaims
	Exp   time.Duration
}

// Something that can be encoded into JWT claims.
type Claim interface {
	ToClaim() ClaimOptions
}

// Returns a claim.
func createClaim(claim Claim) jwt.MapClaims {
	claimOptions := claim.ToClaim()
	claims := claimOptions.Claim
	// Add expiry to the claims if desired
	if claimOptions.Exp != 0 {
		claims["exp"] = time.Now().Add(claimOptions.Exp).Unix()
	}
	return claims
}

type UserClaim struct {
	UserId int32
}

func (uc *UserClaim) ToClaim() ClaimOptions {
	claims := jwt.MapClaims{
		"userId": uc.UserId,
	}

	return ClaimOptions{
		Claim: claims,
		Exp:   time.Hour * 24, // A UserClaim should be valid for 24 hours/1 day
	}
}

// Returns a new UserClaim instance.
func NewUserClaim(userId int32) UserClaim {
	return UserClaim{
		UserId: userId,
	}
}
