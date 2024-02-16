package location

import (
	"merchants.sidooh/api/presenter"
	"merchants.sidooh/pkg/datastore"
	"merchants.sidooh/pkg/entities"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	ReadCounties() (*[]presenter.County, error)
	ReadSubCounties(county int) (*[]presenter.SubCounty, error)
	ReadWards(subCounty int) (*[]presenter.Ward, error)
	ReadLandmarks(ward int) (*[]presenter.Landmark, error)
	GetLandmark(landmarkId string) (*entities.Location, error)
}
type repository struct {
}

func (r *repository) ReadCounties() (counties *[]presenter.County, err error) {
	err = datastore.DB.Model(&entities.Location{}).
		Select("county_id", "county").
		Group("county_id").
		Group("county"). // needed for mysql aggregation in sql_mode=only_full_group_by
		Find(&counties).
		Error
	return
}

func (r *repository) ReadSubCounties(county int) (subCounties *[]presenter.SubCounty, err error) {
	err = datastore.DB.Model(&entities.Location{}).
		Where("county_id", county).
		Select("sub_county_id", "sub_county").
		Group("sub_county_id").
		Group("sub_county").
		Find(&subCounties).
		Error
	return
}

func (r *repository) ReadWards(subCounty int) (wards *[]presenter.Ward, err error) {
	err = datastore.DB.Model(&entities.Location{}).
		Where("sub_county_id", subCounty).
		Select("ward_id", "ward").
		Group("ward_id").
		Group("ward").
		Find(&wards).
		Error
	return
}

func (r *repository) ReadLandmarks(ward int) (landmarks *[]presenter.Landmark, err error) {
	err = datastore.DB.Model(&entities.Location{}).
		Where("ward_id", ward).
		Select("landmark_id", "landmark").
		Group("landmark_id").
		Group("landmark").
		Find(&landmarks).
		Error
	return
}

func (r *repository) GetLandmark(landmarkId string) (landmark *entities.Location, err error) {
	err = datastore.DB.Model(&entities.Location{}).
		Where("landmark_id", landmarkId).
		First(&landmark).
		Error
	return
}

// NewRepo is the single instance repo that is being created.
func NewRepo() Repository {
	return &repository{}
}
