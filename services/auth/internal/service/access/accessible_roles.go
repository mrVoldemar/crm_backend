package access

import (
	"context"
	"github.com/mrVoldemar/crm_backend/services/auth/internal/model"
)

var accessibleRoles map[string]string

func (s *serv) accessibleRoles(ctx context.Context) (map[string]string, error) {
	if accessibleRoles == nil {
		accessibleRoles = make(map[string]string)

		// Лезем в базу за данными о доступных ролях для каждого эндпоинта
		// Можно кэшировать данные, чтобы не лезть в базу каждый раз

		// Например, для эндпоинта /note_v1.NoteV1/Get доступна только роль admin
		accessibleRoles[model.ExamplePath] = "user"
	}

	return accessibleRoles, nil
}
