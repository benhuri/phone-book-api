package contacts

type Contact struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    PhoneNumber string `json:"phone_number"`
    Address   string `json:"address"`
}