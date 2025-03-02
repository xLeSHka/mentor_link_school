package usersRoute

// @Summary Получить список моих запросов
// @Schemes
// @Tags Users
// @Accept json
// @Produce json
// @Router !!!!!!/api/user/requests [get]
// @Success 200 {object} []respGetHelp
//func (h *Route) getGroups(c *gin.Context) {
//	personId, err := jwt.Parse(c)
//	if err != nil {
//		err.(*httpError.HTTPError).SendError(c)
//		c.Abort()
//		return
//	}
//	var req reqGetRole
//	if err := h.validator.ShouldBindQuery(c, &req); err != nil {
//		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
//		c.Abort()
//		return
//	}
//	groups, err := h.usersService.GetGroups(c.Request.Context(), personId)
//	if err != nil {
//		err.(*httpError.HTTPError).SendError(c)
//		return
//	}
//	resp := make([]*respGetGroupDto, 0, len(groups))
//	for _, g := range groups {
//		if g.Group.AvatarURL != nil {
//			avatarURL, err := h.minioRepository.GetImage(*g.Group.AvatarURL)
//			if err != nil {
//				httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
//				c.Abort()
//				return
//			}
//			g.Group.AvatarURL = &avatarURL
//		}
//		resp = append(resp, mapGroup(g, req.Role))
//	}
//	c.JSON(http.StatusOK, resp)
//}
