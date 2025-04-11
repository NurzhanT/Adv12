package domain

import "golang.org/x/crypto/bcrypt"

type User struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

func (u *User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type UserRepository interface {
	Create(user User) error
	GetByUsername(username string) (User, error)
}
