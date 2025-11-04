package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
	"vigovia-task/models"
	"vigovia-task/storage"

	"golang.org/x/crypto/bcrypt"
)

// AuthService handles user authentication operations
type AuthService struct {
	store *storage.MemoryStore
}

// NewAuthService creates a new authentication service
func NewAuthService(store *storage.MemoryStore) *AuthService {
	return &AuthService{
		store: store,
	}
}

// Signup creates a new user account
func (as *AuthService) Signup(req *models.SignupRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	_, err := as.store.GetUserByEmail(req.Email)
	if err == nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate unique user ID
	userID := generateID("user")

	// Create user
	user := &models.User{
		ID:        userID,
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(hashedPassword),
		FullName:  req.FullName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store user
	if err := as.store.CreateUser(user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Store token
	if err := as.store.StoreToken(token, user.ID); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// Login authenticates a user and returns a token
func (as *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	// Get user by email
	user, err := as.store.GetUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate token
	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Store token
	if err := as.store.StoreToken(token, user.ID); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// ValidateToken validates a token and returns the user ID
func (as *AuthService) ValidateToken(token string) (string, error) {
	return as.store.GetUserIDByToken(token)
}

// GetUserByID retrieves a user by ID
func (as *AuthService) GetUserByID(userID string) (*models.User, error) {
	return as.store.GetUserByID(userID)
}

// Logout removes the token
func (as *AuthService) Logout(token string) error {
	return as.store.DeleteToken(token)
}

// generateID generates a unique ID with a prefix
func generateID(prefix string) string {
	timestamp := time.Now().Format("20060102150405")
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("%s-%s-%s", prefix, timestamp, randomStr)
}

// generateToken generates a random token
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
