package request

import "user-service/internal/services/user"

type UpdateUserDataRequest struct {
	Name    string   `json:"name" validate:"required,notblank"`
	Phone   int64    `json:"phone" validate:"required,number"`
	Lat     *float64 `json:"lat" validate:"required,gte=-90,lte=90"`
	Lng     *float64 `json:"lng" validate:"required,gte=-180,lte=180"`
	Address *string  `json:"address" validate:"required,notblank"`
}

type UpdatePasswordRequest struct {
	CurrentPassword      string `json:"current_password" validate:"required"`
	PasswordNew          string `json:"password_new" validate:"required,validpassword"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=PasswordNew"`
}

func ToUpdateUserDataInput(input UpdateUserDataRequest) user.UpdateUserDataInput {
	return user.UpdateUserDataInput{
		Name:    input.Name,
		Phone:   input.Phone,
		Lat:     input.Lat,
		Lng:     input.Lng,
		Address: input.Address,
	}
}
