package users

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *UserRepository
	jwtSecret string
	jwtExpire time.Duration
	logger    *logrus.Logger
}

func NewUserService(repo *UserRepository, jwtSecret string, jwtExpire time.Duration, logger *logrus.Logger) *UserService {
	return &UserService{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtExpire: jwtExpire,
		logger:    logger,
	}
}

func (us *UserService) RegisterUser(req UserRegister) (*UserLoginResponse, error) {
	exists, err := us.repo.EmailExists(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	exists, err = us.repo.UsernameExists(req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("username already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(passwordHash),
	}

	createdUser, err := us.repo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := us.generateToken(createdUser.ID, createdUser.Username, createdUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	us.logger.WithField("user_id", createdUser.ID).Info("User registered successfully")

	return &UserLoginResponse{
		AccessToken: token,
		User:        createdUser.ToResponse(),
		ExpiresAt:   time.Now().Add(us.jwtExpire),
	}, nil
}

func (us *UserService) LoginUser(req UserLogin) (*UserLoginResponse, error) {
	user, err := us.repo.GetByEmail(req.Email)
	if err != nil {
		us.logger.WithError(err).WithField("email", req.Email).Warn("Failed login attempt")
		return nil, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		us.logger.WithError(err).WithField("user_id", user.ID).Warn("Failed login attempt - invalid password")
		return nil, fmt.Errorf("invalid credentials")
	}

	go func() {
		if err := us.repo.UpdateLastSeen(user.ID); err != nil {
			us.logger.WithError(err).WithField("user_id", user.ID).Warn("Failed to update last seen")
		}
	}()

	token, err := us.generateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	us.logger.WithField("user_id", user.ID).Info("User logged in successfully")

	return &UserLoginResponse{
		AccessToken: token,
		User:        user.ToResponse(),
		ExpiresAt:   time.Now().Add(us.jwtExpire),
	}, nil
}

func (us *UserService) GetUserByID(id int) (*UserResponse, error) {
	user, err := us.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	response := user.ToResponse()
	return &response, nil
}

func (us *UserService) UpdateUser(id int, update UserUpdate) (*UserResponse, error) {
	if update.Username != nil {
		exists, err := us.repo.UsernameExists(*update.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to check username existence: %w", err)
		}
		if exists {
			currentUser, err := us.repo.GetByID(id)
			if err != nil {
				return nil, fmt.Errorf("failed to get current user: %w", err)
			}
			if currentUser.Username != *update.Username {
				return nil, fmt.Errorf("username already exists")
			}
		}
	}

	user, err := us.repo.UpdateUser(id, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	us.logger.WithField("user_id", id).Info("User updated successfully")

	response := user.ToResponse()
	return &response, nil
}

func (us *UserService) DeleteUser(id int) error {
	err := us.repo.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	us.logger.WithField("user_id", id).Info("User deleted successfully")

	return nil
}

func (us *UserService) generateToken(userID int, username, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(us.jwtExpire).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(us.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (us *UserService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(us.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
