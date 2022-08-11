package api

import db "github.com/korzepadawid/qr-codes-analyzer/db/sqlc"

func mapUserToResponse(user db.User) userResponse {
	return userResponse{
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
	}
}
