package usersRoute

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/httpError"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"log"
	"net/http"
)

func (h *Route) getGroups(c *gin.Context) {
	personId, err := jwt.Parse(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	var req reqGetRole
	if err := h.validator.ShouldBindQuery(c, &req); err != nil {
		httpError.New(http.StatusBadRequest, err.Error()).SendError(c)
		c.Abort()
		return
	}
	groups, err := h.usersService.GetGroups(c.Request.Context(), personId, req.Role)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		return
	}
	resp := make([]*respGetGroupDto, 0, len(groups))
	for _, g := range groups {
		if g.AvatarURL != nil {
			avatarURL, err := h.minioRepository.GetImage(*g.AvatarURL)
			if err != nil {
				httpError.New(http.StatusInternalServerError, err.Error()).SendError(c)
				c.Abort()
				return
			}
			g.AvatarURL = &avatarURL
		}
		resp = append(resp, mapGroup(g, req.Role))
	}
	c.JSON(http.StatusOK, resp)
}

var upgrader = websocket.Upgrader{
	// Solve cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options
func (r *Route) Websocket(c *gin.Context) {
	_, err := jwt.Parse(c)
	if err != nil {
		err.(*httpError.HTTPError).SendError(c)
		c.Abort()
		return
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	//groups, err := r.usersService.GetGroups(context.Background(), personID, "student")
	//if err != nil {
	//	err.(*httpError.HTTPError).SendError(c)
	//	return
	//}
	//for _, g := range groups {
	//	clients[g.ID] = append(clients[g.ID], ws)
	//}
	for {
		err := ws.WriteMessage(websocket.TextMessage, []byte("hello world"))
		if err != nil {
			log.Println("write:", err)
			ws.Close()
			return
		}
	}
}

//var clients = make(map[uuid.UUID][]*websocket.Conn)
//var mentors = make(chan *models.Role)
//
//func (r *Route) echo() {
//	for {
//		m := <-mentors
//		mentor, err := r.usersService.GetByID(context.Background(), m.UserID)
//		if err != nil {
//			log.Printf("Websocket error: %s", err)
//			continue
//		}
//		m.User = mentor
//		for i, client := range clients[m.GroupID] {
//			err := client.WriteJSON(mapMentor(m))
//			if err != nil {
//				log.Printf("Websocket error: %s", err)
//				client.Close()
//				clients[m.GroupID] = clients[m.GroupID][i : i+1]
//			}
//		}
//	}
//}
//func SendMentor(mentor *models.User) {
//	mentors <- mentor
//}
