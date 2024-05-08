package imports

type Player struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	UsaNum    string `json:"useNum"`
	Level     string `json:"level"`
	TeamNames string `json:"teamNames"`
}
