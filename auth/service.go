package auth

import "github.com/golang-jwt/jwt/v4"

// generate token
// validate token
type Service interface {
	GenerateToken(userID int) (string, error)
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
