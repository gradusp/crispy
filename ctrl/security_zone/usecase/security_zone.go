package usecase

import (
	"github.com/gradusp/crispy/ctrl/security_zone"
)

type SecurityZoneUsecase struct {
	securityZoneRepo security_zone.Repository
}

func NewSecurityZoneUseCase(securityZoneRepo security_zone.Repository) *SecurityZoneUsecase {
	return &SecurityZoneUsecase{
		securityZoneRepo: securityZoneRepo,
	}
}
