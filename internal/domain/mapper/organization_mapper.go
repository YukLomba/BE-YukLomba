package mapper

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
	"github.com/google/uuid"
)

func ToOrganizationShort(org *entity.Organization) *dto.OrganizationShort {
	if org == nil {
		return nil
	}
	return &dto.OrganizationShort{
		ID:   org.ID,
		Name: org.Name,
	}
}

func ToOrganizationResponse(org *entity.Organization) *dto.OrganizationResponse {
	if org == nil {
		return nil
	}
	return &dto.OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Logo:        org.Logo,
		Description: org.Description,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}
}

func ToOrganizationsResponse(orgs []*entity.Organization) []*dto.OrganizationResponse {
	if orgs == nil {
		return nil
	}
	responses := make([]*dto.OrganizationResponse, len(orgs))
	for i, org := range orgs {
		responses[i] = ToOrganizationResponse(org)
	}
	return responses
}

func ToOrganizationFromCreate(req *dto.OrganizationCreateRequest) *entity.Organization {
	return &entity.Organization{
		Name:        req.Name,
		Logo:        req.Logo,
		Description: req.Description,
	}
}

func ToOrganizationFromUpdate(req *dto.OrganizationUpdateRequest, id uuid.UUID) *entity.Organization {
	return &entity.Organization{
		ID:          id,
		Name:        req.Name,
		Logo:        req.Logo,
		Description: req.Description,
	}
}
