package mapper

import (
	"time"

	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
)

func ToCompetitionFromCreate(dto *dto.CompetitionCreateRequest) *entity.Competition {
	competition := entity.Competition{
		Title:       dto.Title,
		Type:        dto.Type,
		Description: dto.Description,
		Deadline:    time.Time(dto.Deadline),
		Category:    dto.Category,
	}
	if dto.OrganizerID != nil {
		competition.OrganizerID = *dto.OrganizerID
	}
	return &competition
}

func ToCompetitionsFromCreate(dto []*dto.CompetitionCreateRequest) []*entity.Competition {
	var comps []*entity.Competition
	for _, com := range dto {
		comps = append(comps, ToCompetitionFromCreate(com))
	}
	return comps
}

func ToCompetitionFromUpdate(dto *dto.CompetitionUpdateRequest) *map[string]interface{} {
	data := make(map[string]interface{})
	if dto.Title != nil {
		data["title"] = *dto.Title
	}
	if dto.Type != nil {
		data["type"] = *dto.Type
	}
	if dto.Description != nil {
		data["description"] = *dto.Description
	}
	if dto.Image != nil {
		data["image"] = dto.Image
	}
	if dto.Deadline != nil {
		data["image"] = *dto.Deadline
	}
	if dto.Category != nil {
		data["category"] = *dto.Category
	}
	if dto.EventLink != nil {
		data["event_link"] = *dto.EventLink
	}
	if dto.Results != nil {
		data["result"] = *dto.Results
	}
	return &data
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
