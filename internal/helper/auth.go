package helper

import (
	"errors"
	"fmt"
	"go-ecommerce-app/internal/domain"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Secret string
}

func SetupAuth(s string) Auth {
	return Auth{
		Secret: s,
	}
}

func (a Auth) CreateHashedPassword(p string) (string, error) {
	if len(p) < 6 {
		return "", errors.New("password length should be at least 6 characters long")
	}

	hashP, err := bcrypt.GenerateFromPassword([]byte(p), 10)

	if err != nil {
		// log actual error and report to loggin tool
		return "", errors.New("password hash failed")
	}

	return string(hashP), nil
}

func (a Auth) GenerateToken(id uint, email string, role string) (string, error) {
	if id == 0 || email == "" || role == "" {
		return "", errors.New("required inputs are missing to generate token")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 dias de token
	})

	tokenStr, err := token.SignedString([]byte(a.Secret))

	if err != nil {
		return "", errors.New("unable to signed the token")
	}

	return tokenStr, nil
}

func (a Auth) VerifyPassword(plainPass string, hashPass string) error {
	if len(plainPass) < 6 {
		return errors.New("password length should be at least 6 characters long")
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(plainPass))

	if err != nil {
		return errors.New("password does not match")
	}

	return nil
}

func (a Auth) VerifyToken(token string) (domain.User, error) {
	tokenArray := strings.Split(token, " ")
	if len(tokenArray) != 2 {
		return domain.User{}, nil
	}

	if tokenArray[0] != "Bearer" {
		return domain.User{}, errors.New("invalid token")
	}

	tokenStr := tokenArray[1]

	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unknown siging mehtod %v", token.Header)
		}

		return []byte(a.Secret), nil
	})

	if err != nil {
		return domain.User{}, errors.New("invalid signing method")
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return domain.User{}, errors.New("token is expired")
		}

		user := domain.User{}
		user.ID = uint(claims["user_id"].(float64))
		user.Email = claims["email"].(string)
		user.UserType = claims["role"].(string)
		return user, nil
	}

	return domain.User{}, errors.New("token verification failed")
}

func (a Auth) Authorize(ctx *fiber.Ctx) error {
	authHeader := ctx.GetReqHeaders()["Authorization"]

	if len(authHeader) == 0 {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  "missing header authorization",
		})
	}

	user, err := a.VerifyToken(authHeader[0])

	if err == nil && user.ID > 0 {
		ctx.Locals("user", user)
		return ctx.Next()
	} else {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  err,
		})
	}
}

func (a Auth) AuthorizeSeller(ctx *fiber.Ctx) error {
	authHeader := ctx.GetReqHeaders()["Authorization"]

	if len(authHeader) == 0 {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  "missing header authorization",
		})
	}

	user, err := a.VerifyToken(authHeader[0])

	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  err,
		})
	} else if user.ID > 0 && user.UserType == domain.SELLER {
		ctx.Locals("user", user)
		return ctx.Next()
	} else {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  errors.New("please join seller program to manage products"),
		})
	}
}

func (a Auth) GetCurrentUser(ctx *fiber.Ctx) domain.User {
	user := ctx.Locals("user")
	return user.(domain.User)
}

func (a Auth) GenerateCode() (string, error) {
	return RandomNumbers(6)
}
