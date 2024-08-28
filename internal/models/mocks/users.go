package mocks

import "github.com/joshdcuneo/go-ui/internal/models"

type UserModel struct{}

const (
	DuplicateUserEmail = "dupe@example.com"
	ValidUserEmail     = "valid@example.com"
	ValidUserPassword  = "validPa$$word"
)

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case DuplicateUserEmail:
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == ValidUserEmail && password == ValidUserPassword {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
