package location

import "merchants.sidooh/api/presenter"

type Service interface {
	GetCounties() (*[]presenter.County, error)
	GetSubCounties(county int) (*[]presenter.SubCounty, error)
	GetWards(subCounty int) (*[]presenter.Ward, error)
	GetLandmarks(ward int) (*[]presenter.Landmark, error)
}

type service struct {
	locationRepository Repository
}

func (s *service) GetCounties() (*[]presenter.County, error) {
	return s.locationRepository.ReadCounties()
}

func (s *service) GetSubCounties(county int) (*[]presenter.SubCounty, error) {
	return s.locationRepository.ReadSubCounties(county)
}

func (s *service) GetWards(subCounty int) (*[]presenter.Ward, error) {
	return s.locationRepository.ReadWards(subCounty)
}

func (s *service) GetLandmarks(ward int) (*[]presenter.Landmark, error) {
	return s.locationRepository.ReadLandmarks(ward)
}

func NewService(location Repository) Service {
	return &service{
		locationRepository: location,
	}
}
