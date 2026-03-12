package service

import (
	"errors"
	"time"

	"github.com/devsherkhane/trello-clone/internal/models"
	"github.com/devsherkhane/trello-clone/internal/repository"
	"github.com/devsherkhane/trello-clone/internal/utils"
)

type AuthService interface {
	Register(username, email, password string) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
	ForgotPassword(email string) error
	ResetPassword(token, newPassword string) error
	UpdateProfile(userID int, username, newEmail string) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{userRepo: repo}
}

func (s *authService) Register(username, email, password string) (*models.User, error) {
	// Check if exists
	existing, _ := s.userRepo.GetByEmail(email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	id, err := s.userRepo.Create(username, email, hashed)
	if err != nil {
		return nil, err
	}
	
	// Create struct and remove sensitive data
	u := &models.User{ID: int(id), Username: username, Email: email}
	return u, nil
}

func (s *authService) Login(email, password string) (string, *models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	user.PasswordHash = "" // Safety clearing
	return token, user, nil
}

func (s *authService) ForgotPassword(email string) error {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		// Prevent email enumeration
		return nil 
	}

	token, err := utils.GenerateResetToken()
	if err != nil {
		return err
	}

	user.ResetToken = token
	user.ResetExpires = time.Now().Add(1 * time.Hour)

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	go utils.SendResetEmail(user.Email, token)
	return nil
}

func (s *authService) ResetPassword(token, newPassword string) error {
	user, err := s.userRepo.GetByResetToken(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashed
	user.ResetToken = ""
	user.ResetExpires = time.Time{}

	return s.userRepo.Update(user)
}

func (s *authService) UpdateProfile(userID int, username, newEmail string) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Update logic safely
	if username != "" {
		user.Username = username
	}
	if newEmail != "" && newEmail != user.Email {
		_, err := s.userRepo.GetByEmail(newEmail)
		if err == nil {
			return nil, errors.New("email is already taken")
		}
		user.Email = newEmail
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}
