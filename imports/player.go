package imports

type Player struct {
	id        int
	email     string
	firstName string
	lastName  string
	phone     string
	usaNum    string
	level     string
	teamNames string
}

type PlayerImport struct {
	Id        int      `json:"id"`
	Email     string   `json:"email"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Phone     string   `json:"phone"`
	UsaNum    string   `json:"useNum"`
	Level     string   `json:"level"`
	TeamNames []string `json:"teamNames"`
}
