package contacts

type Contact struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name" validate:"required,min=1,max=50"`
	LastName    string `json:"last_name" validate:"required,min=1,max=50"`
	PhoneNumber string `json:"phone_number" validate:"required,len=10"`
	Address     string `json:"address" validate:"required,min=2,max=100"`
}
