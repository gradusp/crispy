package usecase

import (
	"github.com/gradusp/crispy/ctrl/securityzone"
)

type SecurityZoneUsecase struct {
	securityZoneRepo securityzone.Repository
}

func NewSecurityZoneUseCase(securityZoneRepo securityzone.Repository) *SecurityZoneUsecase {
	return &SecurityZoneUsecase{
		securityZoneRepo: securityZoneRepo,
	}
}
