package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"time"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (a *AuthRepo) CreateUser(user model.User) error {
	query := "INSERT INTO user(login,username,password) VALUES(?,?,?)"

	if _, err := a.db.Exec(query, user.Login, user.Username, user.Password); err != nil {
		return fmt.Errorf("could't create user: %w", err)
	}

	return nil
}

func (a *AuthRepo) GetUser(login string) (model.User, error) {
	query := "SELECT id,login,password FROM user WHERE login=$1"

	var fullUser model.User

	sqlRow := a.db.QueryRow(query, login)
	if err = sqlRow.Scan(&fullUser.Id, &fullUser.Login, &fullUser.Password); err != nil {
		return model.User{}, fmt.Errorf("couldn't get user: %w", err)
	}

	return fullUser, nil
}

func (a *AuthRepo) GetUserByToken(token string) (model.User, error) {
	query := "SELECT * FROM user WHERE token=$1"

	var fullUser model.User

	sqlRow := a.db.QueryRow(query, token)
	if err = sqlRow.Scan(&fullUser.Id, &fullUser.Password, &fullUser.Login, &fullUser.Username, &fullUser.Token, &fullUser.TokenDuration); err != nil {
		return model.User{}, fmt.Errorf("couldn't get user by token: %w", err)
	}

	return fullUser, nil
}

func (a *AuthRepo) SaveTokens(login string, tokenDuration time.Time, token string) error {
	query := "UPDATE user SET token=$1,tokenDuration=$2 WHERE login=$3"

	if _, err = a.db.Exec(query, token, tokenDuration, login); err != nil {
		return fmt.Errorf("couldn't save token: %w", err)
	}

	return nil
}

func (a *AuthRepo) DeleteToken(token string) error {
	query := "UPDATE user SET token=NULL, tokenDuration=NULL WHERE token=$1"

	if _, err := a.db.Exec(query, token); err != nil {
		return fmt.Errorf("couldn't delete user's token: %w", err)
	}

	return nil
}
