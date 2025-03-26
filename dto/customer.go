package dto

import "go-final/model"

type CustomerResponse struct {
	ID          uint   `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

type UpdateCustomerRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func ConvertUserResponse(cusotomer model.Customer) CustomerResponse {
	return CustomerResponse{
		ID:          uint(cusotomer.CustomerID),
		FirstName:   cusotomer.FirstName,
		LastName:    cusotomer.LastName,
		Email:       cusotomer.Email,
		CreatedAt:   cusotomer.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   cusotomer.UpdatedAt.Format("2006-01-02 15:04:05"),
		PhoneNumber: cusotomer.PhoneNumber,
		Address:     cusotomer.Address,
	}
}
