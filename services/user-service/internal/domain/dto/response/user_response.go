package response

import (
	"user-service/internal/services/user"
	"uuid"
)

type GetProfileResponse struct {
	RoleName []string  `json:"role"`
	ID       int64     `json:"id"`
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Phone    *string   `json:"phone"`
	Lat      *string   `json:"lat"`
	Lng      *string   `json:"lng"`
	Address  *string   `json:"address"`
	PhotoURL *string   `json:"photo"`
}

type UploadImageProfileResponse struct {
	ImageURL *string `json:"image_url"`
}

func ToGetProfileResponse(result *user.GetProfileResult) *GetProfileResponse {
	return &GetProfileResponse{
		RoleName: result.RoleName,
		ID:       result.ID,
		UUID:     result.UUID,
		Name:     result.Name,
		Email:    result.Email,
		Phone:    result.Phone,
		Lat:      result.Lat,
		Lng:      result.Lng,
		Address:  result.Address,
		PhotoURL: result.PhotoURL,
	}
}
