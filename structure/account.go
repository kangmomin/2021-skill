package structure

type Account struct {
	Id        int    `json:"id"`
	AccountId string `json:"accountId"`
	Password  string `json:"accountPassword"`
	Name      string `json:"name"`
	AuthImg   string `json:"authImg"`
	StudentId int    `json:"studentId"`
}
