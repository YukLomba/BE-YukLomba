package mapper

import (
	"github.com/YukLomba/BE-YukLomba/internal/domain/dto"
	"github.com/YukLomba/BE-YukLomba/internal/domain/entity"
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

func ToOrganizationFromUpdate(req *dto.OrganizationUpdateRequest) *map[string]interface{} {
	data := make(map[string]interface{})
	if req.Description != nil {
		data["description"] = *req.Description
	}
	if req.Logo != nil {
		data["logo"] = *req.Logo
	}
	if req.Name != nil {
		data["name"] = *req.Name
	}
	return &data
}
