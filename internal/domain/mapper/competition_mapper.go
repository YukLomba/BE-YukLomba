package mapper

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/google/uuid"
)

func ToCompetitionFromCreate(dto *dto.CompetitionCreateRequest) *entity.Competition {
	return &entity.Competition{
		Title:       dto.Title,
		Type:        dto.Type,
		Description: dto.Description,
		Deadline:    dto.Deadline,
		Category:    dto.Category,
		OrganizerID: *dto.OrganizerID,
	}
}

func ToCompetitionsFromCreate(dto []*dto.CompetitionCreateRequest) []*entity.Competition {
	var comps []*entity.Competition
	for _, com := range dto {
		comps = append(comps, ToCompetitionFromCreate(com))
	}
	return comps
}

func ToCompetitionFromUpdate(dto *dto.CompetitionUpdateRequest, id uuid.UUID) *entity.Competition {
	comps := new(entity.Competition)
	comps.ID = id
	if dto.Title != nil {
		comps.Title = *dto.Title
	}
	if dto.Type != nil {
		comps.Type = *dto.Type
	}
	if dto.Description != nil {
		comps.Description = *dto.Description
	}
	if dto.Image != nil {
		comps.Image = dto.Image
	}
	if dto.Deadline != nil {
		comps.Deadline = *dto.Deadline
	}
	if dto.Category != nil {
		comps.Category = *dto.Category
	}
	if dto.EventLink != nil {
		comps.EventLink = *dto.EventLink
	}
	if dto.Results != nil {
		comps.Results = *dto.Results
	}
	return comps
}

func ToCompetitionResponse(competition *entity.Competition) *dto.CompetitionResponse {
	if competition == nil {
		return nil
	}
	return &dto.CompetitionResponse{
		ID:          competition.ID,
		Title:       competition.Title,
		Type:        competition.Type,
		Description: competition.Description,
		Deadline:    competition.Deadline,
		Category:    competition.Category,
		Organizer: dto.OrganizationShort{
			ID:   competition.Organizer.ID,
			Name: competition.Organizer.Name,
		},
		CreatedAt: competition.CreatedAt,
		UpdatedAt: competition.UpdatedAt,
	}
}

func ToCompetitionsResponse(competitions []*entity.Competition) []*dto.CompetitionResponse {
	var comps []*dto.CompetitionResponse
	for _, comp := range competitions {
		comps = append(comps, ToCompetitionResponse(comp))
	}
	return comps
}

func toCompetitionShort(comp *entity.Competition) *dto.CompetitionShort {
	if comp == nil {
		return nil
	}
	return &dto.CompetitionShort{
		ID:    comp.ID,
		Title: comp.Title,
		Type:  comp.Type,
		Organizer: dto.OrganizationShort{
			ID:   comp.Organizer.ID,
			Name: comp.Organizer.Name,
		},
		Description: comp.Description,
		Deadline:    comp.Deadline,
		Category:    comp.Category,
		CreatedAt:   comp.CreatedAt,
		UpdatedAt:   comp.UpdatedAt,
	}
}

func toCompetitionShorts(comps []*entity.Competition) []*dto.CompetitionShort {
	var compsShort []*dto.CompetitionShort
	for _, comp := range comps {
		compsShort = append(compsShort, toCompetitionShort(comp))
	}
	return compsShort
}
