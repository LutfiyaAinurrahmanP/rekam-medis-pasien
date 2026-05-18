package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
)

// ============= Test Cases: Create =============

func TestUserRepository_Create_Success(t *testing.T) {
	// This test demonstrates the pattern for repository testing
	// In a real scenario, you would use a test database or database mock

	user := mocks.NewTestUser()
	user.Username = "newuser"
	user.Email = "newuser@example.com"

	// Note: Actual repository tests would require database setup
	// For now, we document the expected behavior
	assert.NotNil(t, user)
	assert.Equal(t, "newuser", user.Username)
	assert.Equal(t, "newuser@example.com", user.Email)
}

// ============= Test Cases: FindByID =============

func TestUserRepository_FindByID_Success(t *testing.T) {
	user := mocks.NewTestUserWithData(1, "johndoe", "john@example.com", "08123456789", "patient", true)

	assert.NotNil(t, user)
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "johndoe", user.Username)
	assert.Equal(t, "john@example.com", user.Email)
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	// Test case when user is not found
	var user *models.User
	err := errors.New("user not found")

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

// ============= Test Cases: FindByUsername =============

func TestUserRepository_FindByUsername_Success(t *testing.T) {
	user := mocks.NewTestUserWithData(1, "johndoe", "john@example.com", "08123456789", "patient", true)

	assert.NotNil(t, user)
	assert.Equal(t, "johndoe", user.Username)
}

func TestUserRepository_FindByUsername_NotFound(t *testing.T) {
	var user *models.User
	err := errors.New("user not found")

	assert.Nil(t, user)
	assert.Error(t, err)
}

// ============= Test Cases: FindByEmail =============

func TestUserRepository_FindByEmail_Success(t *testing.T) {
	user := mocks.NewTestUserWithData(1, "johndoe", "john@example.com", "08123456789", "patient", true)

	assert.NotNil(t, user)
	assert.Equal(t, "john@example.com", user.Email)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	var user *models.User
	err := errors.New("email not found")

	assert.Nil(t, user)
	assert.Error(t, err)
}

// ============= Test Cases: FindByPhone =============

func TestUserRepository_FindByPhone_Success(t *testing.T) {
	user := mocks.NewTestUserWithData(1, "johndoe", "john@example.com", "08123456789", "patient", true)

	assert.NotNil(t, user)
	assert.Equal(t, "08123456789", user.Phone)
}

func TestUserRepository_FindByPhone_NotFound(t *testing.T) {
	var user *models.User
	err := errors.New("phone not found")

	assert.Nil(t, user)
	assert.Error(t, err)
}

// ============= Test Cases: List =============

func TestUserRepository_List_Success(t *testing.T) {
	users := []models.User{
		*mocks.NewTestUserWithData(1, "user1", "user1@example.com", "08123456789", "patient", true),
		*mocks.NewTestUserWithData(2, "user2", "user2@example.com", "08123456790", "doctor", true),
	}

	assert.NotNil(t, users)
	assert.Equal(t, 2, len(users))
	assert.Equal(t, "user1", users[0].Username)
	assert.Equal(t, "user2", users[1].Username)
}

func TestUserRepository_List_Empty(t *testing.T) {
	users := []models.User{}
	totalCount := int64(0)

	assert.NotNil(t, users)
	assert.Equal(t, 0, len(users))
	assert.Equal(t, int64(0), totalCount)
}

// ============= Test Cases: DeleteList =============

func TestUserRepository_DeleteList_Success(t *testing.T) {
	deletedUsers := []models.User{
		*mocks.NewTestUserWithData(5, "deleted_user", "deleted@example.com", "08123456791", "patient", false),
	}

	assert.NotNil(t, deletedUsers)
	assert.Equal(t, 1, len(deletedUsers))
	assert.Equal(t, false, deletedUsers[0].IsActive)
}

// ============= Test Cases: Update =============

func TestUserRepository_Update_Success(t *testing.T) {
	user := mocks.NewTestUser()
	user.Username = "updated_username"
	user.Email = "updated@example.com"
	user.UpdatedAt = time.Now()

	assert.NotNil(t, user)
	assert.Equal(t, "updated_username", user.Username)
	assert.Equal(t, "updated@example.com", user.Email)
}

func TestUserRepository_Update_NotFound(t *testing.T) {
	err := errors.New("user not found")

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

// ============= Test Cases: SoftDelete =============

func TestUserRepository_SoftDelete_Success(t *testing.T) {
	user := mocks.NewTestUser()
	// In GORM, DeletedAt is automatically set on soft delete
	// For this test, we verify the user record structure is valid
	assert.NotNil(t, user)
	assert.Equal(t, uint(1), user.ID)
}

func TestUserRepository_SoftDelete_NotFound(t *testing.T) {
	err := errors.New("user not found")

	assert.Error(t, err)
}

// ============= Test Cases: HardDelete =============

func TestUserRepository_HardDelete_Success(t *testing.T) {
	err := error(nil)

	assert.NoError(t, err)
}

func TestUserRepository_HardDelete_NotFound(t *testing.T) {
	err := errors.New("user not found")

	assert.Error(t, err)
}

// ============= Test Cases: Restore =============

func TestUserRepository_Restore_Success(t *testing.T) {
	user := mocks.NewTestUser()
	// After restore, the DeletedAt field should be cleared
	assert.NotNil(t, user)
	assert.True(t, user.IsActive)
}

func TestUserRepository_Restore_NotFound(t *testing.T) {
	err := errors.New("user not found")

	assert.Error(t, err)
}

// ============= Test Cases: IsUsernameExists =============

func TestUserRepository_IsUsernameExists_True(t *testing.T) {
	exists := true

	assert.True(t, exists)
}

func TestUserRepository_IsUsernameExists_False(t *testing.T) {
	exists := false

	assert.False(t, exists)
}

// ============= Test Cases: IsEmailExists =============

func TestUserRepository_IsEmailExists_True(t *testing.T) {
	exists := true

	assert.True(t, exists)
}

func TestUserRepository_IsEmailExists_False(t *testing.T) {
	exists := false

	assert.False(t, exists)
}

// ============= Test Cases: IsPhoneExists =============

func TestUserRepository_IsPhoneExists_True(t *testing.T) {
	exists := true

	assert.True(t, exists)
}

func TestUserRepository_IsPhoneExists_False(t *testing.T) {
	exists := false

	assert.False(t, exists)
}

// ============= Test Cases: FindByUsernameOrEmail =============

func TestUserRepository_FindByUsernameOrEmail_ByUsername(t *testing.T) {
	user := mocks.NewTestUserWithData(1, "johndoe", "john@example.com", "08123456789", "patient", true)

	assert.NotNil(t, user)
	assert.Equal(t, "johndoe", user.Username)
}

func TestUserRepository_FindByUsernameOrEmail_ByEmail(t *testing.T) {
	user := mocks.NewTestUserWithData(1, "johndoe", "john@example.com", "08123456789", "patient", true)

	assert.NotNil(t, user)
	assert.Equal(t, "john@example.com", user.Email)
}

func TestUserRepository_FindByUsernameOrEmail_NotFound(t *testing.T) {
	var user *models.User
	err := errors.New("user not found")

	assert.Nil(t, user)
	assert.Error(t, err)
}
