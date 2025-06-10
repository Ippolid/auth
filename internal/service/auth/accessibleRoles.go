package auth

import (
	"context"
	"fmt"
)

var accessibleRoles map[string]string

func (s *serv) accessibleRoles(ctx context.Context) (map[string]string, error) {
	fmt.Println(accessibleRoles)
	if accessibleRoles == nil {
		accessibleRoles = make(map[string]string)

		cfg := s.access.CFG()

		// Заполняем мапу для эндпоинтов админа
		for endpoint, _ := range cfg {
			accessibleRoles[endpoint] = "admin" //nolint:goconst
		}
	}

	return accessibleRoles, nil
}
