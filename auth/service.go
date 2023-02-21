package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

// generate token
// validate token
type Service interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
}

func NewJWTService() *jwtService {
	return &jwtService{}
}

var SECRET_KEY = []byte("SeCRET_kEY_cro0wdFund1ng")

func (j *jwtService) GenerateToken(userID int) (string, error) {
	payload := jwt.MapClaims{}
	payload["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString(SECRET_KEY)

	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func (j *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(
		encodedToken,
		func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)

			if !ok {
				return nil, errors.New("invalid token")
			}

			return []byte(SECRET_KEY), nil
		})

	if err != nil {
		return token, err
	}

	return token, nil
}
