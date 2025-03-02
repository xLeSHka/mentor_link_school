package groupsRoute

// func (h *Route) getGroups(c *gin.Context) {
// 	personId := uuid.MustParse(c.MustGet("personId").(string))

// 	groups, err := h.groupService.GetGroups(c.Request.Context(), personId)
// 	if err != nil {
// 		err.(*httpError.HTTPError).SendError(c)
// 		return
// 	}
// 	resp := make([]*respGetGroupDto, 0, len(groups))
// 	for _, g := range groups {
// 		if g.AvatarURL != nil {
// 			avatarURL, err := h.minioRepository.GetImage(*g.AvatarURL)
// 			if err != nil {
// 				httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
// 				c.Abort()
// 				return
// 			}
// 			g.AvatarURL = &avatarURL
// 		}
// 		resp = append(resp, mapGroup(g))
// 	}
// 	c.JSON(http.StatusOK, resp)
// }
