package repository

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	List(query *dto.DepartmentPaginationQuery) ([]models.Department, int64, error)
	DeleteList(query *dto.DepartmentPaginationQuery) ([]models.Department, int64, error)
	FindById(id uint) (*models.Department, error)
	FindByName(name string) (*models.Department, error)
	FindByCode(code string) (*models.Department, error)
	Create(department *models.Department) error
	Update(department *models.Department) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	IsCodeExists(code string, excludeID ...uint) (bool, error)
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{
		db: db,
	}
}

func escapeSearchDepartmentPattern(s string) string {
    s = strings.ReplaceAll(s, `\`, `\\`)
    s = strings.ReplaceAll(s, `%`, `\%`)
    s = strings.ReplaceAll(s, `_`, `\_`)
    return s
}

func applyDepartmentListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
    column := "created_at"
    switch sortBy {
    case "name", "code", "created_at":
        column = sortBy
    }

    direction := "DESC"
    if strings.EqualFold(sortDir, "asc") {
        direction = "ASC"
    }

    return db.Order(fmt.Sprintf("%s %s", column, direction))
}

func (r *departmentRepository) List(query *dto.DepartmentPaginationQuery) ([]models.Department, int64, error) {
    var (
        departments []models.Department
        total       int64
    )

    db := r.db.Model(&models.Department{})

    if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
        escaped := escapeSearchDepartmentPattern(searchTerm)
        searchPattern := "%" + escaped + "%"

        db = db.Where(
            "name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\'",
            searchPattern, searchPattern,
        )
    }

    db = applyDepartmentListOrder(db, query.SortBy, query.SortDir)

    countDB := db.Session(&gorm.Session{})
    findDB := db.Session(&gorm.Session{})

    var countErr, findErr error
    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        countErr = countDB.Count(&total).Error
    }()

    go func() {
        defer wg.Done()
        offset := (query.Page - 1) * query.PageSize
        findErr = findDB.Offset(offset).Limit(query.PageSize).Find(&departments).Error
    }()

    wg.Wait()

    if countErr != nil {
        return nil, 0, countErr
    }
    if findErr != nil {
        return nil, 0, findErr
    }

    return departments, total, nil
}

func (r *departmentRepository) DeleteList(query *dto.DepartmentPaginationQuery) ([]models.Department, int64, error) {
	var (
        departments []models.Department
        total       int64
    )

    db := r.db.Unscoped().Model(&models.Department{}).Where("deleted_at IS NOT NULL")

    if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
        escaped := escapeSearchDepartmentPattern(searchTerm)
        searchPattern := "%" + escaped + "%"

        db = db.Where(
            "name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\'",
            searchPattern, searchPattern,
        )
    }

    db = applyDepartmentListOrder(db, query.SortBy, query.SortDir)

    countDB := db.Session(&gorm.Session{})
    findDB := db.Session(&gorm.Session{})

    var countErr, findErr error
    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        countErr = countDB.Count(&total).Error
    }()

    go func() {
        defer wg.Done()
        offset := (query.Page - 1) * query.PageSize
        findErr = findDB.Offset(offset).Limit(query.PageSize).Find(&departments).Error
    }()

    wg.Wait()

    if countErr != nil {
        return nil, 0, countErr
    }
    if findErr != nil {
        return nil, 0, findErr
    }

    return departments, total, nil
}

func (r *departmentRepository) FindById(id uint) (*models.Department, error) {
	var department models.Department
	err := r.db.First(&department, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("department not found")
		}
		return nil, err
	}
	return &department, nil
}

func (r *departmentRepository) FindByName(name string) (*models.Department, error) {
	var department models.Department
	err := r.db.Where("name = ?", name).First(&department).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("name not found")
		}
		return nil, err
	}

	return &department, nil
}

func (r *departmentRepository) FindByCode(code string) (*models.Department, error) {
	var department models.Department
	err := r.db.Where("code = ?", code).First(&department).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("code not found")
		}
		return nil, err
	}
	return &department, nil
}

func (r *departmentRepository) Create(department *models.Department) error {
	return r.db.Create(department).Error
}

func (r *departmentRepository) Update(department *models.Department) error {
	return r.db.Save(department).Error
}

func (r *departmentRepository) SoftDelete(id uint) error {
	return r.db.Delete(&models.Department{}, id).Error
}

func (r *departmentRepository) Restore(id uint) error {
	result := r.db.Model(&models.Department{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("department not found")
	}
	return nil
}

func (r *departmentRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.Department{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("department not found")
	}

	return nil
}

func (r *departmentRepository) IsCodeExists(code string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Unscoped().Model(&models.Department{}).Where("code = ?", code)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	err := query.Count(&count).Error
	return count > 0, err
}
