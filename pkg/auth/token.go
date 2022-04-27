package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/treeverse/lakefs/pkg/db"
	"github.com/treeverse/lakefs/pkg/logging"
)

type TokenVerifier interface {
	VerifyWithAudience(ctx context.Context, token, audience string) (*jwt.StandardClaims, error)
}

type DBTokenVerifier struct {
	db     db.Database
	log    logging.Logger
	secret []byte
}

func VerifyToken(secret []byte, tokenString string) (*jwt.StandardClaims, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (d *DBTokenVerifier) VerifyWithAudience(ctx context.Context, token, audience string) (*jwt.StandardClaims, error) {
	claims, err := VerifyToken(d.secret, token)
	if err != nil {
		return nil, err
	}
	if !claims.VerifyAudience(audience, true) {
		return nil, ErrInvalidToken
	}
	tokenID := claims.Id
	tokenExpiresAt := time.Unix(claims.ExpiresAt, 0)
	canUseToken, err := d.markTokenSingleUse(ctx, tokenID, tokenExpiresAt)
	if err != nil {
		return nil, err
	}
	if !canUseToken {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func NewDBTokenVerifier(db db.Database, log logging.Logger, secret []byte) TokenVerifier {
	return &DBTokenVerifier{
		db:     db,
		log:    log,
		secret: secret,
	}
}

// markTokenSingleUse returns true if token is valid for single use
func (d *DBTokenVerifier) markTokenSingleUse(ctx context.Context, tokenID string, tokenExpiresAt time.Time) (bool, error) {
	res, err := d.db.Exec(ctx, `INSERT INTO expired_tokens (token_id, token_expires_at) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
		tokenID, tokenExpiresAt)
	if err != nil {
		return false, err
	}
	canUseToken := res.RowsAffected() == 1
	// cleanup old tokens
	_, err = d.db.Exec(ctx, `DELETE FROM expired_tokens WHERE token_expires at < $1`, time.Now())
	if err != nil {
		d.log.WithError(err).Error("delete expired tokens")
	}
	return canUseToken, nil
}
