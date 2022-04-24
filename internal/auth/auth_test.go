package auth_test

import (
	"context"
	"fmt"
	"github.com/c0llinn/ebook-store/internal/auth"
	mocks "github.com/c0llinn/ebook-store/mocks/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

const (
	newIdMethod                  = "NewID"
	saveMethod                   = "Save"
	findByEmail                  = "FindByEmail"
	updateMethod                 = "Update"
	generateTokenMethod          = "GenerateTokenForUser"
	newPasswordMethod            = "NewPassword"
	sendEmailMethod              = "SendPasswordResetEmail"
	hashPasswordMethod           = "HashPassword"
	compareHashAndPasswordMethod = "CompareHashAndPassword"
)

type AuthenticatorTestSuite struct {
	suite.Suite
	token             *mocks.TokenHandler
	repo              *mocks.Repository
	emailClient       *mocks.EmailClient
	passwordGenerator *mocks.PasswordGenerator
	hash              *mocks.HashHandler
	idGenerator       *mocks.IDGenerator
	authenticator     *auth.Authenticator
}

func (s *AuthenticatorTestSuite) SetupTest() {
	s.token = new(mocks.TokenHandler)
	s.repo = new(mocks.Repository)
	s.emailClient = new(mocks.EmailClient)
	s.passwordGenerator = new(mocks.PasswordGenerator)
	s.hash = new(mocks.HashHandler)
	s.idGenerator = new(mocks.IDGenerator)

	config := auth.Config{
		Repository:        s.repo,
		Tokener:           s.token,
		Hasher:            s.hash,
		EmailClient:       s.emailClient,
		PasswordGenerator: s.passwordGenerator,
		IDGenerator:       s.idGenerator,
	}

	s.authenticator = auth.New(config)
}

func TestAuthenticator(t *testing.T) {
	suite.Run(t, new(AuthenticatorTestSuite))
}

func (s *AuthenticatorTestSuite) TestRegister_WhenPasswordHashingFails() {
	s.idGenerator.On(newIdMethod).Return("user-id")

	request := auth.RegisterRequest{
		FirstName:            "Raphael",
		LastName:             "Collin",
		Email:                "raphael@test.com",
		Password:             "123456",
		PasswordConfirmation: "123456",
	}

	user := request.User("user-id")
	s.hash.On(hashPasswordMethod, user.Password).Return("", fmt.Errorf("some-error"))

	_, err := s.authenticator.Register(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some-error"), err)

	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 1)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNotCalled(s.T(), saveMethod)
	s.token.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthenticatorTestSuite) TestRegister_WhenRepositoryFails() {
	s.idGenerator.On(newIdMethod).Return("user-id")

	request := auth.RegisterRequest{
		FirstName:            "Raphael",
		LastName:             "Collin",
		Email:                "raphael@test.com",
		Password:             "123456",
		PasswordConfirmation: "123456",
	}

	user := request.User("user-id")
	s.hash.On(hashPasswordMethod, user.Password).Return("hashed-password", nil)

	updatedUser := user
	updatedUser.Password = "hashed-password"
	s.repo.On(saveMethod, context.TODO(), &updatedUser).Return(fmt.Errorf("some error"))

	_, err := s.authenticator.Register(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 1)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), saveMethod, 1)
	s.token.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthenticatorTestSuite) TestRegister_WhenTokenGenerationFails() {
	s.idGenerator.On(newIdMethod).Return("user-id")

	request := auth.RegisterRequest{
		FirstName:            "Raphael",
		LastName:             "Collin",
		Email:                "raphael@test.com",
		Password:             "123456",
		PasswordConfirmation: "123456",
	}

	user := request.User("user-id")
	s.hash.On(hashPasswordMethod, user.Password).Return("hashed-password", nil)

	updatedUser := user
	updatedUser.Password = "hashed-password"
	s.repo.On(saveMethod, context.TODO(), &updatedUser).Return(nil)

	s.token.On(generateTokenMethod, updatedUser).Return("", fmt.Errorf("some error"))

	_, err := s.authenticator.Register(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 1)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), saveMethod, 1)
	s.token.AssertNumberOfCalls(s.T(), generateTokenMethod, 1)
}

func (s *AuthenticatorTestSuite) TestRegister_Successfully() {
	s.idGenerator.On(newIdMethod).Return("user-id")

	request := auth.RegisterRequest{
		FirstName:            "Raphael",
		LastName:             "Collin",
		Email:                "raphael@test.com",
		Password:             "123456",
		PasswordConfirmation: "123456",
	}

	user := request.User("user-id")
	s.hash.On(hashPasswordMethod, user.Password).Return("hashed-password", nil)

	updatedUser := user
	updatedUser.Password = "hashed-password"
	s.repo.On(saveMethod, context.TODO(), &updatedUser).Return(nil)
	s.token.On(generateTokenMethod, updatedUser).Return("token", nil)

	response, err := s.authenticator.Register(context.TODO(), request)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), auth.FromCredentials(auth.Credentials{Token: "token"}), response)

	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 1)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), saveMethod, 1)
	s.token.AssertNumberOfCalls(s.T(), generateTokenMethod, 1)
}

func (s *AuthenticatorTestSuite) TestLogin_WhenUserWasNotFound() {
	request := auth.LoginRequest{
		Email:    "email@test.com",
		Password: "12345678",
	}
	s.repo.On(findByEmail, context.TODO(), request.Email).Return(auth.User{}, fmt.Errorf("some error"))

	_, err := s.authenticator.Login(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.hash.AssertNotCalled(s.T(), compareHashAndPasswordMethod)
	s.token.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthenticatorTestSuite) TestLogin_WhenPasswordsDontMatch() {
	request := auth.LoginRequest{
		Email:    "email@test.com",
		Password: "12345678",
	}

	user := auth.User{ID: "some-id", Email: request.Email, Password: "some-password"}
	s.repo.On(findByEmail, context.TODO(), user.Email).Return(user, nil)
	s.hash.On(compareHashAndPasswordMethod, user.Password, request.Password).Return(auth.ErrWrongPassword)

	_, err := s.authenticator.Login(context.TODO(), request)

	assert.Equal(s.T(), auth.ErrWrongPassword, err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.hash.AssertNumberOfCalls(s.T(), compareHashAndPasswordMethod, 1)
	s.token.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthenticatorTestSuite) TestLogin_Successfully() {
	request := auth.LoginRequest{
		Email:    "email@test.com",
		Password: "12345678",
	}

	user := auth.User{ID: "some-id", Email: request.Email, Password: "some-password"}

	s.repo.On(findByEmail, context.TODO(), user.Email).Return(user, nil)
	s.hash.On(compareHashAndPasswordMethod, user.Password, request.Password).Return(nil)
	s.token.On(generateTokenMethod, user).Return("token", nil)

	response, err := s.authenticator.Login(context.TODO(), request)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), auth.FromCredentials(auth.Credentials{Token: "token"}), response)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.hash.AssertNumberOfCalls(s.T(), compareHashAndPasswordMethod, 1)
	s.token.AssertNumberOfCalls(s.T(), generateTokenMethod, 1)
}

func (s AuthenticatorTestSuite) TestResetPassword_WhenUserWasNotFound() {
	request := auth.PasswordResetRequest{Email: "some email"}
	s.repo.On(findByEmail, context.TODO(), request.Email).Return(auth.User{}, fmt.Errorf("some error"))

	err := s.authenticator.ResetPassword(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.repo.AssertNotCalled(s.T(), updateMethod)
	s.hash.AssertNotCalled(s.T(), hashPasswordMethod)
	s.passwordGenerator.AssertNotCalled(s.T(), newPasswordMethod)
	s.emailClient.AssertNotCalled(s.T(), sendEmailMethod)
}

func (s AuthenticatorTestSuite) TestResetPassword_WhenPasswordHashingFails() {
	request := auth.PasswordResetRequest{Email: "some email"}

	s.repo.On(findByEmail, context.TODO(), request.Email).Return(auth.User{}, nil)
	s.passwordGenerator.On(newPasswordMethod).Return("new-password")
	s.hash.On(hashPasswordMethod, "new-password").Return("", fmt.Errorf("some error"))

	err := s.authenticator.ResetPassword(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.passwordGenerator.AssertNumberOfCalls(s.T(), newPasswordMethod, 1)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNotCalled(s.T(), updateMethod)
	s.emailClient.AssertNotCalled(s.T(), sendEmailMethod)
}

func (s AuthenticatorTestSuite) TestResetPassword_WhenUpdateFails() {
	request := auth.PasswordResetRequest{Email: "some email"}
	user := auth.User{Email: request.Email, Password: "another-password"}
	newPassword := "password"

	s.repo.On(findByEmail, context.TODO(), user.Email).Return(user, nil)
	s.passwordGenerator.On(newPasswordMethod).Return(newPassword)
	s.hash.On(hashPasswordMethod, newPassword).Return("new-hashed-password", nil)
	user.Password = "new-hashed-password"

	s.repo.On(updateMethod, context.TODO(), &user).Return(fmt.Errorf("some error"))

	err := s.authenticator.ResetPassword(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), updateMethod, 1)
	s.passwordGenerator.AssertNumberOfCalls(s.T(), newPasswordMethod, 1)
	s.emailClient.AssertNotCalled(s.T(), sendEmailMethod)
}

func (s AuthenticatorTestSuite) TestResetPassword_WhenEmailSendingFails() {
	request := auth.PasswordResetRequest{Email: "some email"}
	user := auth.User{Email: request.Email, Password: "another-password"}
	newPassword := "password"

	s.repo.On(findByEmail, context.TODO(), user.Email).Return(user, nil)
	s.passwordGenerator.On(newPasswordMethod).Return(newPassword)
	s.hash.On(hashPasswordMethod, newPassword).Return("new-hashed-password", nil)
	user.Password = "new-hashed-password"

	s.repo.On(updateMethod, context.TODO(), &user).Return(nil)
	s.emailClient.On(sendEmailMethod, context.TODO(), user, newPassword).Return(fmt.Errorf("some error"))

	err := s.authenticator.ResetPassword(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), updateMethod, 1)
	s.passwordGenerator.AssertNumberOfCalls(s.T(), newPasswordMethod, 1)
	s.emailClient.AssertNumberOfCalls(s.T(), sendEmailMethod, 1)
}

func (s AuthenticatorTestSuite) TestResetPassword_Successfully() {
	request := auth.PasswordResetRequest{Email: "some email"}
	user := auth.User{Email: request.Email, Password: "another-password"}
	newPassword := "password"

	s.repo.On(findByEmail, context.TODO(), user.Email).Return(user, nil)
	s.passwordGenerator.On(newPasswordMethod).Return(newPassword)
	s.hash.On(hashPasswordMethod, newPassword).Return("new-hashed-password", nil)
	user.Password = "new-hashed-password"

	s.repo.On(updateMethod, context.TODO(), &user).Return(nil)
	s.emailClient.On(sendEmailMethod, context.TODO(), user, newPassword).Return(nil)

	err := s.authenticator.ResetPassword(context.TODO(), request)

	assert.Nil(s.T(), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), updateMethod, 1)
	s.passwordGenerator.AssertNumberOfCalls(s.T(), newPasswordMethod, 1)
	s.emailClient.AssertNumberOfCalls(s.T(), sendEmailMethod, 1)
}
