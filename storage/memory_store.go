package storage

import (
	"fmt"
	"sync"

	"vigovia-task/models"
)

// MemoryStore is an in-memory storage for itineraries and users
type MemoryStore struct {
	itineraries map[string]*models.Itinerary
	users       map[string]*models.User      // key: user ID
	usersByEmail map[string]*models.User     // key: email for quick lookup
	tokens      map[string]string            // key: token, value: user ID
	mu          sync.RWMutex
}

// NewMemoryStore creates a new instance of MemoryStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		itineraries:  make(map[string]*models.Itinerary),
		users:        make(map[string]*models.User),
		usersByEmail: make(map[string]*models.User),
		tokens:       make(map[string]string),
	}
}

// Create stores a new itinerary
func (ms *MemoryStore) Create(itinerary *models.Itinerary) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.itineraries[itinerary.ID]; exists {
		return fmt.Errorf("itinerary with id %s already exists", itinerary.ID)
	}

	ms.itineraries[itinerary.ID] = itinerary
	return nil
}

// GetByID retrieves an itinerary by ID
func (ms *MemoryStore) GetByID(id string) (*models.Itinerary, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	itinerary, exists := ms.itineraries[id]
	if !exists {
		return nil, fmt.Errorf("itinerary with id %s not found", id)
	}

	return itinerary, nil
}

// GetAll retrieves all itineraries
func (ms *MemoryStore) GetAll() []*models.Itinerary {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	itineraries := make([]*models.Itinerary, 0, len(ms.itineraries))
	for _, itinerary := range ms.itineraries {
		itineraries = append(itineraries, itinerary)
	}

	return itineraries
}

// Update updates an existing itinerary
func (ms *MemoryStore) Update(id string, itinerary *models.Itinerary) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.itineraries[id]; !exists {
		return fmt.Errorf("itinerary with id %s not found", id)
	}

	ms.itineraries[id] = itinerary
	return nil
}

// Delete removes an itinerary
func (ms *MemoryStore) Delete(id string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.itineraries[id]; !exists {
		return fmt.Errorf("itinerary with id %s not found", id)
	}

	delete(ms.itineraries, id)
	return nil
}

// User-related methods

// CreateUser stores a new user
func (ms *MemoryStore) CreateUser(user *models.User) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.users[user.ID]; exists {
		return fmt.Errorf("user with id %s already exists", user.ID)
	}

	if _, exists := ms.usersByEmail[user.Email]; exists {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	ms.users[user.ID] = user
	ms.usersByEmail[user.Email] = user
	return nil
}

// GetUserByEmail retrieves a user by email
func (ms *MemoryStore) GetUserByEmail(email string) (*models.User, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	user, exists := ms.usersByEmail[email]
	if !exists {
		return nil, fmt.Errorf("user with email %s not found", email)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (ms *MemoryStore) GetUserByID(id string) (*models.User, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	user, exists := ms.users[id]
	if !exists {
		return nil, fmt.Errorf("user with id %s not found", id)
	}

	return user, nil
}

// StoreToken stores an authentication token
func (ms *MemoryStore) StoreToken(token, userID string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.tokens[token] = userID
	return nil
}

// GetUserIDByToken retrieves a user ID by token
func (ms *MemoryStore) GetUserIDByToken(token string) (string, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	userID, exists := ms.tokens[token]
	if !exists {
		return "", fmt.Errorf("invalid token")
	}

	return userID, nil
}

// DeleteToken removes a token
func (ms *MemoryStore) DeleteToken(token string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	delete(ms.tokens, token)
	return nil
}
