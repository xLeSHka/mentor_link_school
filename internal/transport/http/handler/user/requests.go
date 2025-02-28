package usersRoute

import "mime/multipart"

type reqSigninDto struct {
	Password string `json:"password" binding:"required,min=8,max=60"`
	Email    string `json:"email" binding:"required,min=8,max=60,email"`
}
type resSigninDto struct {
	Token string `json:"token"`
}

type resGetProfile struct {
	Name      string  `json:"name"`
	Surname   string  `json:"surname"`
	Email     string  `json:"email"`
	AvatarUrl *string `json:"avatar_url,omitempty"`
}

type reqSignupDto struct {
	Password  string  `json:"password" binding:"required,min=8,max=60"`
	Name      string  `json:"name" binding:"required,min=1,max=100"`
	Surname   string  `json:"surname" binding:"required,min=1,max=120"`
	Email     string  `json:"email" binding:"required,min=8,max=120,email"`
	Age       int     `json:"age" binding:"required,min=1,max=120"`
	AvatarUrl *string `json:"avatar_url" binding:"omitempty,max=350"`
}
type resSignupDto struct {
	Token string `json:"token"`
}
type reqUploadAvatarDto struct {
	Image *multipart.FileHeader `form:"image"`
}
type respUploadAvatarDto struct {
	Url string `json:"url"`
}
