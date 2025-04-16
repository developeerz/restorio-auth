package mapper

import "github.com/developeerz/restorio-auth/internal/repository/postgres/models"

func UserAuthToIDAndAuth(userAuths []models.UserAuth) (int64, []string) {
	id := userAuths[0].UserTelegramID

	auths := make([]string, len(userAuths))
	for i, v := range userAuths {
		auths[i] = string(v.AuthID)
	}

	return id, auths
}
