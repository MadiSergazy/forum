package service

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/repository"
	"github.com/zhayt/clean-arch-tmp-forum/logger"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

const salt = "Sfasfasfasfas"

var EmailRX = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

var empty model.User

type AuthService struct {
	repo repository.Authorization
	l    *logger.Logger
}

func NewAuthService(repo repository.Authorization, l *logger.Logger) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(user model.User) error {
	if !checkLogin(user.Login) {
		fmt.Println("checkLogin")
		return InvalidData
	}
	if !checkPassword(user.Password) {
		fmt.Println("checkPassword")
		return InvalidData
	}
	if !checkUsername(user.Username) {
		fmt.Println("checkUsername")
		return InvalidData
	}

	user.Password = generatePassword(user.Password)

	err := s.repo.CreateUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed:") {
			return fmt.Errorf("not unique data: %w", InvalidData)
		}

		return err
	}

	return nil
}

func (s *AuthService) GenerateToken(login, password string) (model.User, error) {
	user, err := s.repo.GetUser(login)
	if err != nil {
		return empty, fmt.Errorf("%s: %w", err.Error(), InvalidData)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+salt)); err != nil {
		return empty, fmt.Errorf("%s: %w", err.Error(), InvalidData)
	}

	token, err := uuid.NewV4()
	if err != nil {
		return empty, fmt.Errorf("couldn't generate token: %w", err)
	}

	user.Token = token.String()
	user.TokenDuration = time.Now().Add(12 * time.Hour)
	if err := s.repo.SaveTokens(user.Login, user.TokenDuration, user.Token); err != nil {
		return empty, nil
	}
	return user, nil
}

func generatePassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	return string(bytes)
}

func (s *AuthService) GetUserByToken(token string) (model.User, error) {
	return s.repo.GetUserByToken(token)
}

func (s *AuthService) DeleteToken(token string) error {
	return s.repo.DeleteToken(token)
}

func checkUsername(username string) bool {
	for _, w := range username {
		if w >= 33 && w <= 126 {
			return true
		} else {
			return false
		}
	}
	return true
}

func checkPassword(password string) bool {
	for _, rune := range password {
		if rune > 32 && rune <= 126 && len(password) > 8 {
			return true
		} else {
			return false
		}
	}

	return true
}

func checkLogin(email string) bool {
	_, err := mail.ParseAddress(email)

	//if !EmailRX.MatchString(email) {
	//	return false
	//}

	if err != nil {
		return false
	}
	return true
}
