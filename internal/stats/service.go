package stats

// Service handles statistics business logic
type Service struct {
	repo *Repository
}

// NewService creates a new stats service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetUserStats retrieves statistics for a user
func (s *Service) GetUserStats(userID string) (*UserStats, error) {
	return s.repo.GetUserStats(userID)
}
