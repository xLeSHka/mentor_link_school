package usersRoute

type reqLoginDto struct {
	Name string `json:"name" binding:"required"`
}
type respLoginDto struct {
	Token string `json:"token"`
}

type resGetProfile struct {
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatar_url,omitempty"`
	BIO       *string `json:"bio,omitempty"`
}

type reqGetHelpRequest struct {
	MentorID string `json:"mentor_id" binding:"required"`
	Goal     string `json:"goal" binding:"required"`
}
type respGetHelp struct {
	MentorName string  `json:"mentor_name"`
	Goal       string  `json:"goal"`
	BIO        *string `json:"bio,omitempty"`
	Status     string  `json:"status"`
}
type respUploadAvatarDto struct {
	Url string `json:"url"`
}
