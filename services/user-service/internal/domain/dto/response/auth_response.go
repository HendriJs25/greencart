package response

import (
	authservice "user-service/internal/services/auth"
	"uuid"
)

type SignInResponse struct {
	CSRFToken string   `json:"csrf_token"`
	User      UserInfo `json:"user"`
}

type UserInfo struct {
	UUID  uuid.UUID `json:"uuid"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Phone *string   `json:"phone"`
	Lat   *string   `json:"lat"`
	Lng   *string   `json:"lng"`
}

func ToSignInResponse(result *authservice.SignInResult) SignInResponse {
	return SignInResponse{
		CSRFToken: result.CSRFToken,
		User: UserInfo{
			UUID:  result.UUID,
			Name:  result.Name,
			Email: result.Email,
			Phone: result.Phone,
			Lat:   result.Lat,
			Lng:   result.Lng,
		},
	}
}
