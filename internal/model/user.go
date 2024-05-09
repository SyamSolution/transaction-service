package model

type User struct {
	UserID      int    `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	City        string `json:"city"`
	Country     string `json:"country"`
	PostalCode  string `json:"postal_code"`
	NIK         string `json:"nik"`
}

type ResponseUser struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}
