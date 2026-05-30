package service

import (
	"errors"
	"messageboard/dto/request"
	"messageboard/dto/response"
	"messageboard/model"
	"messageboard/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(req *request.RegisterRequest) (*model.User, error) {
	// Check if username already exists
	if s.userRepo.ExistsByUsername(req.Username) {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
	}

	if user.Nickname == "" {
		user.Nickname = req.Username
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(req *request.LoginRequest) (*response.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	token, expireAt, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		Token:    token,
		ExpireAt: expireAt,
	}, nil
}

func (s *AuthService) generateToken(userID uint) (string, time.Time, error) {
	expireAt := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expireAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(model.GetJWTSecret()))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expireAt, nil
}

func (s *AuthService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *AuthService) UpdateUser(id uint, req *request.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) UpdateAvatar(id uint, avatarURL string) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	user.Avatar = avatarURL
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}