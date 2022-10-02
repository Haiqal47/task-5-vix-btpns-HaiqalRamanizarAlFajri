package app

type UserRegister struct {
	Username string `json:"username"`
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"minstringlength(6)~Password must be longer than 6 characters"`
}

type UserLogin struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password"`
}

type UserUpdate struct {
	Username string `json:"username" valid:"optional"`
	Password string `json:"password" valid:"minstringlength(6)~Password must be longer than 6 characters,optional"`
}

type PhotoCreated struct {
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
	UserId   string `json:"user_id"`
}

type PhotoUpdate struct {
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
	UserId   string `json:"user_id"`
}
