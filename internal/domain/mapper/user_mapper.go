package mapper

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
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

func MapUserUpdate(dto *dto.UserProfileUpdate) *map[string]interface{} {
	userMap := make(map[string]interface{})
	if dto.Username != nil {
		userMap["username"] = dto.Username
	}
	if dto.Email != nil {
		userMap["email"] = dto.Email
	}
	if dto.University != nil {
		userMap["university"] = dto.University
	}
	if dto.Interests != nil {
		userMap["interests"] = dto.Interests
	}
	if dto.Password != nil {
		userMap["password"] = dto.Password
	}
	return &userMap
}
