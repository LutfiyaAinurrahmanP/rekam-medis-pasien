package service

import (
"errors"
"math"

"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type TypeTestService interface {
List(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error)
ListActive(query *dto.TypeTestPaginationQuery) (*dto.ActiveTypeTestListResponse, error)
ListInactive(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error)
DeleteList(query *dto.TypeTestPaginationQuery) (*dto.TypeTestDeletedListResponse, error)
Search(query *dto.TypeTestSearchQuery) (*dto.TypeTestSearchResponse, error)
FindByID(id uint) (*dto.TypeTestResponse, error)
FindByCode(code string) (*dto.TypeTestResponse, error)
FindByCategory(category string, query *dto.TypeTestPaginationQuery) (*dto.TypeTestCategoryListResponse, error)
Create(req *dto.CreateTypeTestRequest) (*dto.TypeTestResponse, error)
Update(id uint, req *dto.UpdateTypeTestRequest) (*dto.TypeTestResponse, error)
Activate(id uint) error
Deactivate(id uint) error
SoftDelete(id uint) error
Restore(id uint) error
HardDelete(id uint) error
}

type typeTestService struct {
repo   repository.TypeTestRepository
config *config.Config
}

func NewTypeTestService(repo repository.TypeTestRepository, config *config.Config) TypeTestService {
return &typeTestService{repo: repo, config: config}
}

func (s *typeTestService) normalizeQuery(query *dto.TypeTestPaginationQuery, defaultSortBy, defaultSortDir string) {
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
query.SortBy = defaultSortBy
}
if query.SortDir == "" {
query.SortDir = defaultSortDir
}
}

func (s *typeTestService) List(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error) {
s.normalizeQuery(query, "name", "asc")

typeTests, total, err := s.repo.List(query)
if err != nil {
return nil, err
}

responses := make([]dto.TypeTestResponse, len(typeTests))
for i, t := range typeTests {
responses[i] = *s.toResponse(&t)
}

totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
return &dto.TypeTestListResponse{
Data: responses,
Meta: dto.TypeTestPaginationMeta{
Page:       query.Page,
PageSize:   query.PageSize,
TotalItems: total,
TotalPages: totalPages,
},
}, nil
}

func (s *typeTestService) ListActive(query *dto.TypeTestPaginationQuery) (*dto.ActiveTypeTestListResponse, error) {
s.normalizeQuery(query, "name", "asc")

typeTests, total, err := s.repo.ListActive(query)
if err != nil {
return nil, err
}

categories, err := s.repo.CountActiveByCategory()
if err != nil {
return nil, err
}

responses := make([]dto.TypeTestResponse, len(typeTests))
for i, t := range typeTests {
responses[i] = *s.toResponse(&t)
}

totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
return &dto.ActiveTypeTestListResponse{
TotalActiveTests: total,
Categories:       categories,
Data:             responses,
Meta: dto.TypeTestPaginationMeta{
Page:       query.Page,
PageSize:   query.PageSize,
TotalItems: total,
TotalPages: totalPages,
},
}, nil
}

func (s *typeTestService) ListInactive(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error) {
s.normalizeQuery(query, "name", "asc")

typeTests, total, err := s.repo.ListInactive(query)
if err != nil {
return nil, err
}

responses := make([]dto.TypeTestResponse, len(typeTests))
for i, t := range typeTests {
responses[i] = *s.toResponse(&t)
}

totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
return &dto.TypeTestListResponse{
Data: responses,
Meta: dto.TypeTestPaginationMeta{
Page:       query.Page,
PageSize:   query.PageSize,
TotalItems: total,
TotalPages: totalPages,
},
}, nil
}

func (s *typeTestService) DeleteList(query *dto.TypeTestPaginationQuery) (*dto.TypeTestDeletedListResponse, error) {
s.normalizeQuery(query, "deleted_at", "desc")

typeTests, total, err := s.repo.DeleteList(query)
if err != nil {
return nil, err
}

responses := make([]dto.DeletedTypeTestResponse, len(typeTests))
for i, t := range typeTests {
responses[i] = *s.toDeletedResponse(&t)
}

totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
return &dto.TypeTestDeletedListResponse{
Data: responses,
Meta: dto.TypeTestPaginationMeta{
Page:       query.Page,
PageSize:   query.PageSize,
TotalItems: total,
TotalPages: totalPages,
},
}, nil
}

func (s *typeTestService) Search(query *dto.TypeTestSearchQuery) (*dto.TypeTestSearchResponse, error) {
if query.Page < 1 {
query.Page = 1
}
if query.PageSize < 1 {
query.PageSize = s.config.Pagination.DefaultPageSize
}
if query.PageSize > s.config.Pagination.MaxPageSize {
query.PageSize = s.config.Pagination.MaxPageSize
}

typeTests, total, err := s.repo.Search(query)
if err != nil {
return nil, err
}

responses := make([]dto.TypeTestResponse, len(typeTests))
for i, t := range typeTests {
responses[i] = *s.toResponse(&t)
}

criteria := map[string]interface{}{}
if query.Keyword != "" {
criteria["keyword"] = query.Keyword
}
if query.Category != "" {
criteria["category"] = query.Category
}
if query.MinPrice != nil {
criteria["min_price"] = *query.MinPrice
}
if query.MaxPrice != nil {
criteria["max_price"] = *query.MaxPrice
}
if query.IsActive != nil {
criteria["is_active"] = *query.IsActive
}

totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
return &dto.TypeTestSearchResponse{
SearchCriteria: criteria,
ResultsFound:   total,
Data:           responses,
Meta: dto.TypeTestPaginationMeta{
Page:       query.Page,
PageSize:   query.PageSize,
TotalItems: total,
TotalPages: totalPages,
},
}, nil
}

func (s *typeTestService) FindByID(id uint) (*dto.TypeTestResponse, error) {
t, err := s.repo.FindByID(id)
if err != nil {
return nil, err
}
return s.toResponse(t), nil
}

func (s *typeTestService) FindByCode(code string) (*dto.TypeTestResponse, error) {
t, err := s.repo.FindByCode(code)
if err != nil {
return nil, err
}
return s.toResponse(t), nil
}

func (s *typeTestService) FindByCategory(category string, query *dto.TypeTestPaginationQuery) (*dto.TypeTestCategoryListResponse, error) {
s.normalizeQuery(query, "name", "asc")

typeTests, total, err := s.repo.FindByCategory(category, query)
if err != nil {
return nil, err
}

responses := make([]dto.TypeTestResponse, len(typeTests))
for i, t := range typeTests {
responses[i] = *s.toResponse(&t)
}

totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
return &dto.TypeTestCategoryListResponse{
Category:   category,
TotalTests: total,
Data:       responses,
Meta: dto.TypeTestPaginationMeta{
Page:       query.Page,
PageSize:   query.PageSize,
TotalItems: total,
TotalPages: totalPages,
},
}, nil
}

func (s *typeTestService) Create(req *dto.CreateTypeTestRequest) (*dto.TypeTestResponse, error) {
exists, err := s.repo.IsCodeExists(req.Code)
if err != nil {
return nil, err
}
if exists {
return nil, errors.New("code already exists")
}

t := &models.TypeTest{
Name:     req.Name,
Code:     req.Code,
IsActive: true,
}
if req.Category != nil {
t.Category = *req.Category
}
if req.Description != nil {
t.Description = *req.Description
}
if req.Price != nil {
t.Price = *req.Price
}
if req.IsActive != nil {
t.IsActive = *req.IsActive
}

if err := s.repo.Create(t); err != nil {
return nil, err
}
return s.toResponse(t), nil
}

func (s *typeTestService) Update(id uint, req *dto.UpdateTypeTestRequest) (*dto.TypeTestResponse, error) {
t, err := s.repo.FindByID(id)
if err != nil {
return nil, err
}

if req.Name != nil {
t.Name = *req.Name
}
if req.Code != nil {
exists, err := s.repo.IsCodeExists(*req.Code, id)
if err != nil {
return nil, err
}
if exists {
return nil, errors.New("code already exists")
}
t.Code = *req.Code
}
if req.Category != nil {
t.Category = *req.Category
}
if req.Description != nil {
t.Description = *req.Description
}
if req.Price != nil {
t.Price = *req.Price
}
if req.IsActive != nil {
t.IsActive = *req.IsActive
}

if err := s.repo.Update(t); err != nil {
return nil, err
}
return s.toResponse(t), nil
}

func (s *typeTestService) Activate(id uint) error {
return s.repo.Activate(id)
}

func (s *typeTestService) Deactivate(id uint) error {
return s.repo.Deactivate(id)
}

func (s *typeTestService) SoftDelete(id uint) error {
return s.repo.SoftDelete(id)
}

func (s *typeTestService) Restore(id uint) error {
return s.repo.Restore(id)
}

func (s *typeTestService) HardDelete(id uint) error {
return s.repo.HardDelete(id)
}

func (s *typeTestService) toResponse(t *models.TypeTest) *dto.TypeTestResponse {
return &dto.TypeTestResponse{
ID:          t.ID,
Name:        t.Name,
Code:        t.Code,
Category:    t.Category,
Description: t.Description,
Price:       t.Price,
IsActive:    t.IsActive,
CreatedAt:   t.CreatedAt,
UpdatedAt:   t.UpdatedAt,
}
}

func (s *typeTestService) toDeletedResponse(t *models.TypeTest) *dto.DeletedTypeTestResponse {
r := &dto.DeletedTypeTestResponse{
ID:          t.ID,
Name:        t.Name,
Code:        t.Code,
Category:    t.Category,
Description: t.Description,
Price:       t.Price,
IsActive:    t.IsActive,
CreatedAt:   t.CreatedAt,
UpdatedAt:   t.UpdatedAt,
}
if t.DeletedAt.Valid {
deletedAt := t.DeletedAt.Time
r.DeletedAt = &deletedAt
}
return r
}
