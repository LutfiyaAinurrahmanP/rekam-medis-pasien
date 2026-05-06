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

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByPhone(phone string) (*models.User, error)
	FindByUsernameOrEmail(usernameOrEmail string) (*models.User, error)
	List(query *dto.UserPaginationQuery) ([]models.User, int64, error)
	DeleteList(query *dto.UserPaginationQuery) ([]models.User, int64, error)
	Update(user *models.User) error
	SoftDelete(id uint) error
	HardDelete(id uint) error
	Restore(id uint) error
	IsUsernameExists(username string, excludeID ...uint) (bool, error)
	IsEmailExists(email string, excludeID ...uint) (bool, error)
	IsPhoneExists(phone string, excludeID ...uint) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsernameOrEmail(usernameOrEmail string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func escapeSearchPattern(s string) string {
    // Escape karakter spesial SQL LIKE: % _ \
    s = strings.ReplaceAll(s, `\`, `\\`)
    s = strings.ReplaceAll(s, `%`, `\%`)
    s = strings.ReplaceAll(s, `_`, `\_`)
    return s
}

func (r *userRepository) List(query *dto.UserPaginationQuery) ([]models.User, int64, error) {
    var users []models.User
    var total int64

    db := r.db.Model(&models.User{})

    if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
        // ✅ Escape dulu, baru wrap dengan %
        escaped := escapeSearchPattern(searchTerm)
        searchPattern := "%" + escaped + "%"

        // ✅ Tambahkan ESCAPE '\' agar PostgreSQL tahu karakter escape-nya
        db = db.Where(
            "username ILIKE ? ESCAPE '\\' OR email ILIKE ? ESCAPE '\\' OR phone ILIKE ? ESCAPE '\\'",
            searchPattern, searchPattern, searchPattern,
        )
    }

    if query.Role != "" {
        db = db.Where("role = ?", query.Role)
    }

    if query.IsActive != nil {
        db = db.Where("is_active = ?", *query.IsActive)
    }

    db = applyUserListOrder(db, query.SortBy, query.SortDir)

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
        findErr = findDB.Offset(offset).Limit(query.PageSize).Find(&users).Error
    }()

    wg.Wait()

    if countErr != nil {
        return nil, 0, countErr
    }
    if findErr != nil {
        return nil, 0, findErr
    }

    return users, total, nil
}

func (r *userRepository) DeleteList(query *dto.UserPaginationQuery) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := r.db.Unscoped().Model(&models.User{}).Where("deleted_at IS NOT NULL")

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		searchPattern := "%" + searchTerm + "%"
		db = db.Where("username ILIKE ? OR email ILIKE ? OR phone ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if query.Role != "" {
		db = db.Where("role = ?", query.Role)
	}

	if query.IsActive != nil {
		db = db.Where("is_active = ?", query.IsActive)
	}

	db = applyUserListOrder(db, query.SortBy, query.SortDir)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func applyUserListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
	column := "created_at"
	switch sortBy {
	case "username", "email", "created_at":
		column = sortBy
	}

	direction := "DESC"
	if strings.EqualFold(sortDir, "asc") {
		direction = "ASC"
	}

	return db.Order(fmt.Sprintf("%s %s", column, direction))
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) SoftDelete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *userRepository) Restore(id uint) error {
	result := r.db.Model(&models.User{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *userRepository) IsUsernameExists(username string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.User{}).Where("username = ?", username)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *userRepository) IsEmailExists(email string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.User{}).Where("email = ?", email)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *userRepository) IsPhoneExists(phone string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.User{}).Where("phone = ?", phone)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}
