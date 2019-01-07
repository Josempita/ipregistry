package model

type Config struct {
	Templates string `json:"templates"`
	Root      string `json:"root"`
	Password  string `json:"password"`
	Username  string `json:"username"`
}

type MessagerService struct {
	Messager chan Details
}

type curve interface {
	GetConfig() Config
}

type Model struct {
	curve
}

type Details struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func New(curve curve) *Model {
	return &Model{
		curve: curve,
	}
}
