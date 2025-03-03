package userService

import "github.com/google/uuid"

func (s *UsersService) GetCommonGroups(userID, mentorID uuid.UUID) ([]uuid.UUID, error) {
	return s.usersRepository.GetCommonGroups(userID, mentorID)
}
