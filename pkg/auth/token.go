package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/treeverse/lakefs/pkg/db"
)

type TokenVerifier interface {
	VerifyWithAudience(ctx context.Context, token, audience string) (*jwt.StandardClaims, error)
}

type DBTokenVerifier struct {
	secret []byte
	db     db.Database
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
	if err != nil || !claims.VerifyAudience(audience, true) {
		return nil, ErrInvalidToken
	}
	expired, err := d.IsTokenExpired(ctx, claims.Id)
	if expired || err != nil {
		return nil, ErrInvalidToken
	}
	err = d.ExpireToken(ctx, claims.Id, time.Unix(claims.ExpiresAt, 0))
	if err != nil {
		return nil, err
	}
	err = d.RemoveExpiredByTime(ctx)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func NewDBTokenVerifier(db db.Database, secret []byte) TokenVerifier {
	return &DBTokenVerifier{db: db, secret: secret}
}

// RemoveExpiredByTime removes all the tokens that the token_expires_at passed
func (d *DBTokenVerifier) RemoveExpiredByTime(ctx context.Context) error {
	_, err := d.db.Exec(ctx, `DELETE FROM expired_tokens Where token_expires at < $1`, time.Now())
	return err
}

func (d *DBTokenVerifier) ExpireToken(ctx context.Context, tokenID string, tokenExpiresAt time.Time) error {
	_, err := d.db.Exec(ctx, `INSERT INTO expired_tokens (token_id, token_expires_at)
			VALUES ($1,$2)`, tokenID, tokenExpiresAt)
	return err
}

func (d *DBTokenVerifier) IsTokenExpired(ctx context.Context, tokenID string) (bool, error) {
	tokens := 0
	err := d.db.Get(ctx, &tokens, `SELECT COUNT(*) FROM expired_tokens where token_id=$1`, tokenID)
	return tokens >= 0, err
}
