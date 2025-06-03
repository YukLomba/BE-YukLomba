package mapper

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/google/uuid"
)

func ToUserResponse(user *entity.User) *dto.UserResponse {
	if user == nil {
		return nil
	}
	return &dto.UserResponse{
		ID:                 user.ID.String(),
		Username:           user.Username,
		Email:              user.Email,
		University:         user.University,
		Interests:          user.Interests,
		CreatedAt:          user.CreatedAt,
		Organization:       ToOrganizationShort(user.Organization),
		JoinedCompetitions: toCompetitionShorts(user.JoinedCompetitions),
	}
}

func ToUsersResponse(users []*entity.User) []*dto.UserResponse {
	var usersResponse []*dto.UserResponse
	for _, user := range users {
		usersResponse = append(usersResponse, ToUserResponse(user))
	}
	return usersResponse
}

func ToUserWithID(dto *dto.UserProfileUpdate, id uuid.UUID) *entity.User {
	user := new(entity.User)
	user.ID = id
	if dto.Username != nil {
		user.Username = *dto.Username
	}
	if dto.Email != nil {
		user.Email = *dto.Email
	}
	if dto.University != nil {
		user.University = *dto.University
	}
	if dto.Interests != nil {
		user.Interests = *dto.Interests
	}
	return user
}
