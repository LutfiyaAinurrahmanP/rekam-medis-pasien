package dto

import "time"

type RoomListResponse struct {
	Data []RoomResponse     `json:"data"`
	Meta RoomPaginationMeta `json:"meta"`
}

type RoomDeletedListResponse struct {
	Data []DeletedRoomResponse `json:"data"`
	Meta RoomPaginationMeta    `json:"meta"`
}

type CreateRoomRequest struct {
	RoomNumber    string   `json:"room_number" binding:"required,min=1,max=20"`
	RoomType      string   `json:"room_type" binding:"required,oneof=vip class_1 class_2 class_3 icu emergency"`
	DepartmentID  *uint    `json:"department_id" binding:"omitempty"`
	BedCapacity   int      `json:"bed_capacity" binding:"required,min=1"`
	AvailableBeds *int     `json:"available_beds" binding:"omitempty,min=0"`
	PricePerDay   *float64 `json:"price_per_day" binding:"omitempty,min=0"`
	IsActive      *bool    `json:"is_active" binding:"omitempty"`
}

type UpdateRoomRequest struct {
	RoomNumber   *string  `json:"room_number" binding:"omitempty,min=1,max=20"`
	RoomType     *string  `json:"room_type" binding:"omitempty,oneof=vip class_1 class_2 class_3 icu emergency"`
	DepartmentID *uint    `json:"department_id" binding:"omitempty"`
	BedCapacity  *int     `json:"bed_capacity" binding:"omitempty,min=1"`
	PricePerDay  *float64 `json:"price_per_day" binding:"omitempty,min=0"`
	IsActive     *bool    `json:"is_active" binding:"omitempty"`
}

type OccupyRoomRequest struct {
	Beds int `json:"beds" binding:"required,min=1"`
}

type ReleaseRoomRequest struct {
	Beds int `json:"beds" binding:"required,min=1"`
}

type RoomResponse struct {
	ID            uint            `json:"id"`
	RoomNumber    string          `json:"room_number"`
	RoomType      string          `json:"room_type"`
	DepartmentID  *uint           `json:"department_id,omitempty"`
	Department    *DepartmentInfo `json:"department,omitempty"`
	BedCapacity   int             `json:"bed_capacity"`
	AvailableBeds int             `json:"available_beds"`
	PricePerDay   float64         `json:"price_per_day"`
	IsActive      bool            `json:"is_active"`
	OccupancyRate float64         `json:"occupancy_rate"`
	IsAvailable   bool            `json:"is_available,omitempty"`
	IsFull        bool            `json:"is_full,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at,omitempty"`
}

type DeletedRoomResponse struct {
	ID            uint            `json:"id"`
	RoomNumber    string          `json:"room_number"`
	RoomType      string          `json:"room_type"`
	DepartmentID  *uint           `json:"department_id,omitempty"`
	Department    *DepartmentInfo `json:"department,omitempty"`
	BedCapacity   int             `json:"bed_capacity"`
	AvailableBeds int             `json:"available_beds"`
	PricePerDay   float64         `json:"price_per_day"`
	IsActive      bool            `json:"is_active"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at,omitempty"`
	DeletedAt     *time.Time      `json:"deleted_at"`
}

type DepartmentInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code,omitempty"`
}

type RoomPaginationQuery struct {
	Page     int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search" binding:"omitempty"`
	SortBy   string `form:"sort_by,default=created_at" binding:"omitempty,oneof=created_at room_number room_type available_beds"`
	SortDir  string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type RoomPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type AvailableRoomsResponse struct {
	TotalAvailableRooms int                `json:"total_available_rooms"`
	TotalAvailableBeds  int                `json:"total_available_beds"`
	Data                []RoomResponse     `json:"data"`
	Meta                RoomPaginationMeta `json:"meta,omitempty"`
}
