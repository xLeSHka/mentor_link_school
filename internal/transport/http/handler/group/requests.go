package groupsRoute

type reqGetMentorDto struct {
	GroupEmail string `json:"group_email" binding:"required,min=8,max=120,email"`
}
type GetGroupID struct {
	ID string `uri:"groupId" binding:"required,uuid"`
}
type reqCreateGroupDto struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}
