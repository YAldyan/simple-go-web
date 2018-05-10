package controller

import (
	"example-go-web/util"
	"crypto/md5"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string
	Email          string
	HashedPassword string
	Username       string
}

const (
	hashCost       = 10
	passwordLength = 6
	userIDLength   = 16
)

func (user *User) AvatarURL() string {
	return fmt.Sprintf(
		"//www.gravatar.com/avatar/%x",
		md5.Sum([]byte(user.Email)),
	)
}

func (user *User) AvatarUser() string {
	return "/assets/images/orang.png"
}

func (user *User) ImagesRoute() string {
	return "/user/" + user.ID
}

func NewUser(username, email, password string) (User, error) {
	user := User{
		Email:    email,
		Username: username,
	}
	if username == "" {
		return user, util.ErrNoUsername
	}

	if email == "" {
		return user, util.ErrNoEmail
	}

	if password == "" {
		return user, util.ErrNoPassword
	}

	if len(password) < passwordLength {
		return user, util.ErrPasswordTooShort
	}

	// Check if the username exists
	existingUser, err := GlobalUserStore.FindByUsername(username)
	if err != nil {
		return user, err
	}
	if existingUser != nil {
		return user, util.ErrUsernameExists
	}

	// Check if the email exists
	existingUser, err = GlobalUserStore.FindByEmail(email)
	if err != nil {
		return user, err
	}
	if existingUser != nil {
		return user, util.ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)

	user.HashedPassword = string(hashedPassword)
	user.ID = util.GenerateID("usr", userIDLength)
	return user, err
}

func FindUser(username, password string) (*User, error) {
	out := &User{
		Username: username,
	}

	existingUser, err := GlobalUserStore.FindByUsername(username)
	if err != nil {
		return out, err
	}
	if existingUser == nil {
		return out, util.ErrCredentialsIncorrect
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(existingUser.HashedPassword),
		[]byte(password),
	) != nil {
		return out, util.ErrCredentialsIncorrect
	}

	return existingUser, nil
}

func UpdateUser(user *User, email, currentPassword, newPassword string) (User, error) {
	out := *user
	out.Email = email

	// Check if the email exists
	existingUser, err := GlobalUserStore.FindByEmail(email)
	if err != nil {
		return out, err
	}
	if existingUser != nil && existingUser.ID != user.ID {
		return out, util.ErrEmailExists
	}

	// At this point, we can update the email address
	user.Email = email

	// No current password? Don't try update the password.
	if currentPassword == "" {
		return out, nil
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(user.HashedPassword),
		[]byte(currentPassword),
	) != nil {
		return out, util.ErrPasswordIncorrect
	}

	if newPassword == "" {
		return out, util.ErrNoPassword
	}

	if len(newPassword) < passwordLength {
		return out, util.ErrPasswordTooShort
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), hashCost)
	user.HashedPassword = string(hashedPassword)
	return out, err
}
