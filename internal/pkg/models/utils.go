package models

import (
	"github.com/badoux/checkmail"
	"regexp"
)

var (
	reSlugValidator     = regexp.MustCompile(`^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$`)
	reNicknameValidator = regexp.MustCompile(`([a-zA-Z0-9]|_|\.)*`)
)

func validateEmail(email string) error {
	if email != "" {
		if err := checkmail.ValidateFormat(email); err != nil {
			return ErrValidation
		}
	}
	return nil
}

func validateNickname(nickname string) error {
	if !reNicknameValidator.MatchString(nickname) {
		return ErrValidation
	}
	return nil
}

func validateSlug(slug string) error {
	if !reSlugValidator.MatchString(slug) {
		return ErrValidation
	}
	return nil
}
