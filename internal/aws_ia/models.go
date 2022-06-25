package aws_ia

type User struct {
	Names          string `json:"names"`
	Surname        string `json:"surname"`
	SecondSurname  string `json:"second_surname"`
	IdentityNumber string `json:"identity_number"`
}
