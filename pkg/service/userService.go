package service

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"ozinshe/pkg/entity"
	"ozinshe/pkg/helpers"
	"regexp"
	"time"
)

type UserService interface {
	SignUp(*entity.User) error
	PasswordValidator(string) error
	VerifyAccount(string) error
	SigIn(*entity.Credentials) (*entity.User, error)
	TokenGenerator(int, string, string) (string, error)
	ChangePasswordByUserId(int, string, string) error
	ConfirmPasswordValidator(string, string) error
	PasswordRecover(email string) error
}

func (s *Service) SignUp(user *entity.User) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Printf("error while creating hash of password in SignUp(Service): %s", err.Error())
		return err
	}
	user.Password = string(hashedPass)
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !regex.MatchString(user.Email) {
		return errors.New(entity.InvalidEmail)
	}
	ExistedUser, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		if err.Error() != entity.DidNotFind {
			return err
		}
		if err = s.repo.CreateUser(user); err != nil {
			return err
		}
	} else if !ExistedUser.IsEmailVerified {
		// if user registered but not verified and try to register again
		err = s.repo.DeleteVerificationEmailByUserId(ExistedUser.Id)
		if err != nil {
			s.log.Printf("error during deleting verification email in VerifyAccount(Service):", err.Error())
			return err
		}
		user.Id = ExistedUser.Id
		if err = s.repo.UpdateUserByID(user); err != nil {
			return err
		}
	} else {
		return errors.New(entity.AlreadyExist)
	}
	emailContent, secretCode, err := s.VerificationEmailGenerator(user.Email)
	if err != nil {
		s.log.Printf("error during verification email creation in SignUp(Service):", err.Error())
		return err
	}
	if err = s.SendVerificationEmail(user.Email, emailContent); err != nil {
		s.log.Printf("error during verification email sending in SignUp(Service):", err.Error())
		return errors.New("invalid email")
	}
	if err = s.CreateVerificationEmail(user.Id, secretCode); err != nil {
		s.log.Printf("error during verification email creating in SignUp(Service):", err.Error())
		return err
	}
	return nil
}

func (s *Service) PasswordValidator(password string) error {
	regex := regexp.MustCompile(`[A-Za-z0-9]*[@$!&]*.{7,}$`)
	if !regex.MatchString(password) {
		return errors.New("A password should be alphanumeric.\nFirst letter of the password should be capital.\nPassword must contain a special character (@, $, !, &, etc).\nPassword length must be greater than 8 characters.")
	}
	return nil
}

func (s *Service) ConfirmPasswordValidator(password, confirmPassword string) error {
	if confirmPassword != password {
		return errors.New(entity.InvalidConfirmPassword)
	}
	return nil
}
func (s *Service) VerifyAccount(secretCode string) error {
	verificationEmail, err := s.repo.GetVerificationEmailStatusBySecretCode(secretCode)
	if err != nil {
		s.log.Printf("error during getting verification email in VerifyAccount(Service):", err.Error())
		return err
	}
	if verificationEmail.ExpTime.Before(time.Now()) {
		s.log.Printf("error in VerifyAccount(Service) %s", err.Error())
		err = s.repo.DeleteVerificationEmailByUserId(verificationEmail.UserId)
		if err != nil {
			s.log.Printf("error during deleting verification email in VerifyAccount(Service):", err.Error())
			return err
		}
		return errors.New("497")
	}
	err = s.repo.UpdateUsersEmailStatus(verificationEmail.UserId)
	if err != nil {
		s.log.Printf("error during updating user email status in VerifyAccount(Service):", err.Error())
		return err
	}
	err = s.repo.DeleteVerificationEmailByUserId(verificationEmail.UserId)
	if err != nil {
		s.log.Printf("error during deleting verification email in VerifyAccount(Service):", err.Error())
		return err
	}
	return nil
}

func (s *Service) SigIn(credentials *entity.Credentials) (*entity.User, error) {
	user, err := s.repo.GetUserByEmail(credentials.Email)
	if err != nil {
		// can be now rows in result set
		return nil, err
	} else if !user.IsEmailVerified {
		//not verified email
		return nil, errors.New("email not verified")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		s.log.Printf("given password  is incorrect: %s", credentials.Password)
		return nil, fmt.Errorf(entity.InvalidPassword)
	}
	return user, nil
}

//func (s *Service) GetPasswordByUserId(userId int) (string,error) {
//	return s.repo.GetPasswordByUserId(userId)
//}

func (s *Service) ChangePasswordByUserId(userId int, oldPassword, newPassword string) error {
	currentPassword, err := s.repo.GetPasswordByUserId(userId)
	if err != nil {
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(currentPassword), []byte(oldPassword)); err != nil {
		return fmt.Errorf(entity.InvalidPassword)
	}
	newHashedPass, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		s.log.Printf("error while creating hash of newPassword in ChangePasswordByUserId(Service): %s", err.Error())
		return err
	}
	return s.repo.ChangePasswordByUserId(userId, string(newHashedPass))
}

func (s *Service) PasswordRecover(email string) error {
	tempPassword := helpers.GeneratePassword()
	tempPasswordHash, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	emailContent := "Your new password is - " + tempPassword + " Do not share it and change immediately!"
	err = s.repo.ChangePasswordByEmail(email, string(tempPasswordHash))
	if err != nil {
		return err
	}
	err = s.SendVerificationEmail(email, emailContent)
	return err
}
