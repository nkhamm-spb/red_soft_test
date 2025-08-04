package schemas

type User struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Surname     string   `json:"surname"`
	Gender      string   `json:"gender"`
	Age         int      `json:"age"`
	Nationalize string   `json:"nationalize"`
	Emails      []string `json:"emails"`
}

type NewUser struct {
	Name    string   `json:"name"`
	Surname string   `json:"surname"`
	Emails  []string `json:"emails"`
}

type EditUser struct {
	Name        string   `json:"name"`
	Surname     string   `json:"surname"`
	Gender      string   `json:"gender"`
	Age         int      `json:"age"`
	Nationalize string   `json:"nationalize"`
	Emails      []string `json:"emails"`
}
