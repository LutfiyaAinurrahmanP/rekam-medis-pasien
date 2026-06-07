package room

import (
	"errors"
	"math"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type RoomService interface {
	ListRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error)
	DeleteListRooms(query *dto.RoomPaginationQuery) (*dto.RoomDeletedListResponse, error)
	GetAvailableRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error)
	GetOccupiedRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error)
	GetActiveRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error)
	GetInactiveRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error)
	GetRoomByID(id uint) (*dto.RoomResponse, error)
	CreateRoom(req *dto.CreateRoomRequest) (*dto.RoomResponse, error)
	UpdateRoom(id uint, req *dto.UpdateRoomRequest) (*dto.RoomResponse, error)
	ActivateRoom(id uint) (*dto.RoomResponse, error)
	DeactivateRoom(id uint) (*dto.RoomResponse, error)
	OccupyRoom(id uint, beds int) (*dto.RoomResponse, error)
	ReleaseRoom(id uint, beds int) (*dto.RoomResponse, error)
	SoftDeleteRoom(id uint) error
	RestoreRoom(id uint) error
	HardDeleteRoom(id uint) error
}

type roomService struct {
	repo   repository.RoomRepository
	config *config.Config
}

func NewRoomService(repo repository.RoomRepository, config *config.Config) RoomService {
	return &roomService{
		repo:   repo,
		config: config,
	}
}

func (s roomService) ListRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}

	if query.PageSize < 1 {
		query.PageSize = s.config.Pagination.DefaultPageSize
	}

	if query.PageSize > s.config.Pagination.MaxPageSize {
		query.PageSize = s.config.Pagination.MaxPageSize
	}

	if query.SortBy == "" {
		query.SortBy = "created_at"
	}

	if query.SortDir == "" {
		query.SortDir = "desc"
	}

	rooms, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	roomResponses := make([]dto.RoomResponse, len(rooms))
	for i, room := range rooms {
		roomResponses[i] = *s.toRoomResponse(&room)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.RoomListResponse{
		Data: roomResponses,
		Meta: dto.RoomPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s roomService) DeleteListRooms(query *dto.RoomPaginationQuery) (*dto.RoomDeletedListResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}

	if query.PageSize < 1 {
		query.PageSize = s.config.Pagination.DefaultPageSize
	}

	if query.PageSize > s.config.Pagination.MaxPageSize {
		query.PageSize = s.config.Pagination.MaxPageSize
	}

	if query.SortBy == "" {
		query.SortBy = "deleted_at"
	}

	if query.SortDir == "" {
		query.SortDir = "desc"
	}

	rooms, total, err := s.repo.DeleteList(query)
	if err != nil {
		return nil, err
	}

	roomResponses := make([]dto.DeletedRoomResponse, len(rooms))
	for i, room := range rooms {
		roomResponses[i] = *s.toDeleteRoomResponse(&room)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.RoomDeletedListResponse{
		Data: roomResponses,
		Meta: dto.RoomPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s roomService) GetAvailableRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}

	if query.PageSize < 1 {
		query.PageSize = s.config.Pagination.DefaultPageSize
	}

	if query.PageSize > s.config.Pagination.MaxPageSize {
		query.PageSize = s.config.Pagination.MaxPageSize
	}

	if query.SortBy == "" {
		query.SortBy = "available_beds"
	}

	if query.SortDir == "" {
		query.SortDir = "desc"
	}

	rooms, total, err := s.repo.FindAvailableRooms(query)
	if err != nil {
		return nil, err
	}

	roomResponses := make([]dto.RoomResponse, len(rooms))
	for i, room := range rooms {
		roomResponses[i] = *s.toRoomResponse(&room)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.RoomListResponse{
		Data: roomResponses,
		Meta: dto.RoomPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s roomService) GetOccupiedRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}

	if query.PageSize < 1 {
		query.PageSize = s.config.Pagination.DefaultPageSize
	}

	if query.PageSize > s.config.Pagination.MaxPageSize {
		query.PageSize = s.config.Pagination.MaxPageSize
	}

	if query.SortBy == "" {
		query.SortBy = "available_beds"
	}

	if query.SortDir == "" {
		query.SortDir = "asc"
	}

	rooms, total, err := s.repo.FindOccupiedRooms(query)
	if err != nil {
		return nil, err
	}

	roomResponses := make([]dto.RoomResponse, len(rooms))
	for i, room := range rooms {
		roomResponses[i] = *s.toRoomResponse(&room)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.RoomListResponse{
		Data: roomResponses,
		Meta: dto.RoomPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s roomService) GetActiveRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}

	if query.PageSize < 1 {
		query.PageSize = s.config.Pagination.DefaultPageSize
	}

	if query.PageSize > s.config.Pagination.MaxPageSize {
		query.PageSize = s.config.Pagination.MaxPageSize
	}

	if query.SortBy == "" {
		query.SortBy = "created_at"
	}

	if query.SortDir == "" {
		query.SortDir = "desc"
	}

	rooms, total, err := s.repo.FindActiveRooms(query)
	if err != nil {
		return nil, err
	}

	roomResponses := make([]dto.RoomResponse, len(rooms))
	for i, room := range rooms {
		roomResponses[i] = *s.toRoomResponse(&room)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.RoomListResponse{
		Data: roomResponses,
		Meta: dto.RoomPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s roomService) GetInactiveRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}

	if query.PageSize < 1 {
		query.PageSize = s.config.Pagination.DefaultPageSize
	}

	if query.PageSize > s.config.Pagination.MaxPageSize {
		query.PageSize = s.config.Pagination.MaxPageSize
	}

	if query.SortBy == "" {
		query.SortBy = "created_at"
	}

	if query.SortDir == "" {
		query.SortDir = "desc"
	}

	rooms, total, err := s.repo.FindInactiveRooms(query)
	if err != nil {
		return nil, err
	}

	roomResponses := make([]dto.RoomResponse, len(rooms))
	for i, room := range rooms {
		roomResponses[i] = *s.toRoomResponse(&room)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.RoomListResponse{
		Data: roomResponses,
		Meta: dto.RoomPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s roomService) GetRoomByID(id uint) (*dto.RoomResponse, error) {
	room, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toRoomResponse(room), nil
}

func (s roomService) CreateRoom(req *dto.CreateRoomRequest) (*dto.RoomResponse, error) {
	exists, err := s.repo.IsRoomNumberExists(req.RoomNumber)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("room number already exists")
	}

	// Set defaults for optional fields
	availableBeds := req.BedCapacity // Default to bed_capacity
	if req.AvailableBeds != nil {
		availableBeds = *req.AvailableBeds
	}

	pricePerDay := 0.0 // Default to 0
	if req.PricePerDay != nil {
		pricePerDay = *req.PricePerDay
	}

	isActive := true // Default to true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	room := &models.Room{
		RoomNumber:    req.RoomNumber,
		RoomTypeID:    req.RoomTypeID,
		DepartmentID:  req.DepartmentID,
		BedCapacity:   req.BedCapacity,
		AvailableBeds: availableBeds,
		PricePerDay:   pricePerDay,
		IsActive:      isActive,
	}

	if err := s.repo.Create(room); err != nil {
		return nil, err
	}

	return s.toRoomResponse(room), nil
}

func (s roomService) UpdateRoom(id uint, req *dto.UpdateRoomRequest) (*dto.RoomResponse, error) {
	room, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.RoomNumber != nil && *req.RoomNumber != room.RoomNumber {
		exists, err := s.repo.IsRoomNumberExists(*req.RoomNumber, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("room number already exists")
		}
		room.RoomNumber = *req.RoomNumber
	}

	if req.RoomTypeID != nil {
		room.RoomTypeID = req.RoomTypeID
	}

	if req.DepartmentID != nil {
		room.DepartmentID = req.DepartmentID
	}

	if req.BedCapacity != nil {
		room.BedCapacity = *req.BedCapacity
	}

	if req.PricePerDay != nil {
		room.PricePerDay = *req.PricePerDay
	}

	if req.IsActive != nil {
		room.IsActive = *req.IsActive
	}

	if err := s.repo.Update(room); err != nil {
		return nil, err
	}

	return s.toRoomResponse(room), nil
}

func (s roomService) ActivateRoom(id uint) (*dto.RoomResponse, error) {
	room, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	room.IsActive = true

	if err := s.repo.Update(room); err != nil {
		return nil, err
	}

	return s.toRoomResponse(room), nil
}
func (s roomService) DeactivateRoom(id uint) (*dto.RoomResponse, error) {
	room, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	room.IsActive = false

	if err := s.repo.Update(room); err != nil {
		return nil, err
	}

	return s.toRoomResponse(room), nil
}

func (s roomService) OccupyRoom(id uint, beds int) (*dto.RoomResponse, error) {
	room, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if beds <= 0 {
		return nil, errors.New("number of beds must be greater than 0")
	}

	if room.AvailableBeds < beds {
		return nil, errors.New("not enough available beds")
	}

	room.AvailableBeds -= beds

	if err := s.repo.Update(room); err != nil {
		return nil, err
	}

	return s.toRoomResponse(room), nil
}

func (s roomService) ReleaseRoom(id uint, beds int) (*dto.RoomResponse, error) {
	room, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if beds <= 0 {
		return nil, errors.New("number of beds must be greater than 0")
	}

	newAvailableBeds := room.AvailableBeds + beds
	if newAvailableBeds > room.BedCapacity {
		return nil, errors.New("cannot release more beds than capacity")
	}

	room.AvailableBeds = newAvailableBeds

	if err := s.repo.Update(room); err != nil {
		return nil, err
	}

	return s.toRoomResponse(room), nil
}

func (s roomService) SoftDeleteRoom(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.SoftDelete(id)
}

func (s roomService) RestoreRoom(id uint) error {
	return s.repo.Restore(id)
}

func (s roomService) HardDeleteRoom(id uint) error {
	return s.repo.HardDelete(id)
}

func (s *roomService) toRoomResponse(room *models.Room) *dto.RoomResponse {
	return &dto.RoomResponse{
		ID:            room.ID,
		RoomNumber:    room.RoomNumber,
		RoomTypeID:    room.RoomTypeID,
		DepartmentID:  room.DepartmentID,
		BedCapacity:   room.BedCapacity,
		AvailableBeds: room.AvailableBeds,
		PricePerDay:   room.PricePerDay,
		IsActive:      room.IsActive,
		CreatedAt:     room.CreatedAt,
		UpdatedAt:     room.UpdatedAt,
	}
}

func (s *roomService) toDeleteRoomResponse(room *models.Room) *dto.DeletedRoomResponse {

	return &dto.DeletedRoomResponse{
		ID:            room.ID,
		RoomNumber:    room.RoomNumber,
		RoomTypeID:    room.RoomTypeID,
		DepartmentID:  room.DepartmentID,
		BedCapacity:   room.BedCapacity,
		AvailableBeds: room.AvailableBeds,
		PricePerDay:   room.PricePerDay,
		IsActive:      room.IsActive,
		CreatedAt:     room.CreatedAt,
		UpdatedAt:     room.UpdatedAt,
		DeletedAt:     &room.DeletedAt.Time,
	}
}
