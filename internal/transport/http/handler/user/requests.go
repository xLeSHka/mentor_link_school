package usersRoute

type reqSigninDto struct {
	Email    string `json:"email" binding:"required,min=8,max=60,email"`
	Password string `json:"password" binding:"required,min=8,max=60,c-password"`
}
type resSigninDto struct {
	Token string `json:"token"`
}

type resGetProfile struct {
	Name      string  `json:"name"`
	Surname   string  `json:"surname"`
	Email     string  `json:"email"`
	AvatarUrl *string `json:"avatar_url,omitempty"`
	BIO       *string `json:"bio,omitempty"`
}

type reqSignupDto struct {
	FirstName  string `form:"first_name" binding:"required,max=60"`
	SecondName string `form:"second_name" binding:"required,max=60"`
	Password   string `json:"password" binding:"required,min=8,max=60,c-password"`
	Email      string `json:"email" binding:"required,min=8,max=120,email"`
}
type resSignupDto struct {
	Token string `json:"token"`
}
type respUploadAvatarDto struct {
	Url string `json:"url"`
}
type reqCreateGroupDto struct {
	Name string `json:"name"`
}
