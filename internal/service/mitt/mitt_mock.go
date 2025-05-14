package mitt

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/misshanya/mitter/internal/models"
)

var mockMittModel = &models.Mitt{
	ID:        uuid.New(),
	AuthorID:  mockUserID,
	Content:   "hello world",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Likes:     0,
}

var mockUserID = uuid.New()

// Mock mitt repo
type mockMittRepo struct{}

func (m mockMittRepo) CreateMitt(ctx context.Context, userID uuid.UUID, mitt *models.MittCreate) (*models.Mitt, error) {
	_ = ctx
	_ = userID
	_ = mitt

	return mockMittModel, nil
}

func (m mockMittRepo) GetMitt(ctx context.Context, id uuid.UUID) (*models.Mitt, error) {
	_ = ctx
	_ = id

	return mockMittModel, nil
}

func (m mockMittRepo) GetAllUserMitts(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*models.Mitt, error) {
	_ = ctx
	_ = userID
	_ = limit
	_ = offset

	return []*models.Mitt{mockMittModel}, nil
}

func (m mockMittRepo) UpdateMitt(ctx context.Context, mittID uuid.UUID, mitt *models.MittUpdate) (*models.Mitt, error) {
	_ = ctx
	_ = mittID

	// Update mock mitt content
	mockMittModel.Content = mitt.Content

	return mockMittModel, nil
}

func (m mockMittRepo) DeleteMitt(ctx context.Context, mittID uuid.UUID) error {
	_ = ctx
	_ = mittID

	// Nothing happens :)
	return nil
}

func (m mockMittRepo) LikeMitt(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) error {
	_ = ctx
	_ = userID
	_ = mittID

	// Add like to mock mitt
	mockMittModel.Likes++

	return nil
}

func (m mockMittRepo) IsMittLikedByUser(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) (bool, error) {
	_ = ctx
	_ = userID
	_ = mittID

	// Let's assume that if mock model has >=1 like, mitt is liked by user, if less, no
	return mockMittModel.Likes >= 1, nil
}

func (m mockMittRepo) DeleteMittLike(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) error {
	_ = ctx
	_ = userID
	_ = mittID

	// Delete like (oops, it can be negative cause you can run tests in different order)
	mockMittModel.Likes--

	return nil
}

func (m mockMittRepo) GetMittLikesCount(ctx context.Context, mittID uuid.UUID) (int64, error) {
	_ = ctx
	_ = mittID

	return mockMittModel.Likes, nil
}

func (m mockMittRepo) Feed(ctx context.Context, limit, offset int32) ([]*models.Mitt, error) {
	_ = ctx
	_ = limit
	_ = offset

	return []*models.Mitt{mockMittModel}, nil
}

// Mock User repo
type mockUserRepo struct{}

func (r *mockUserRepo) CreateUser(ctx context.Context, user *models.UserCreate) (uuid.UUID, error) {
	_ = ctx
	_ = user

	return mockUserID, nil
}

func (r *mockUserRepo) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	_ = ctx
	_ = login

	return &models.User{ID: mockUserID}, nil
}

func (r *mockUserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	_ = ctx
	_ = id

	return &models.User{ID: mockUserID}, nil
}

func (r *mockUserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_ = ctx
	_ = id

	return nil
}

func (r *mockUserRepo) UpdateUser(ctx context.Context, id uuid.UUID, user *models.UserUpdate) error {
	_ = ctx
	_ = id
	_ = user

	return nil
}

func (r *mockUserRepo) GetCurrentPasswordHash(ctx context.Context, id uuid.UUID) (string, error) {
	_ = ctx
	_ = id

	return "$2a$10$OW9yD0TyX0pOBO2MzJhtpeOC6O694OS37VJnnaJKFm.rUFt5fy4O6", nil
}

func (r *mockUserRepo) ChangePassword(ctx context.Context, id uuid.UUID, newHashedPassword string) error {
	_ = ctx
	_ = id
	_ = newHashedPassword

	return nil
}

func (r *mockUserRepo) FollowUser(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) error {
	_ = ctx
	_ = followerID
	_ = followeeID

	return nil
}

func (r *mockUserRepo) UnfollowUser(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) error {
	_ = ctx
	_ = followerID
	_ = followeeID

	return nil
}

func (r *mockUserRepo) GetUserFollows(ctx context.Context, followerID uuid.UUID, limit, offset int32) ([]uuid.UUID, error) {
	_ = ctx
	_ = followerID
	_ = limit
	_ = offset

	return []uuid.UUID{mockUserID}, nil
}

func (r *mockUserRepo) GetUserFollowers(ctx context.Context, followeeID uuid.UUID, limit, offset int32) ([]uuid.UUID, error) {
	_ = ctx
	_ = followeeID
	_ = limit
	_ = offset

	return []uuid.UUID{mockUserID}, nil
}

func (r *mockUserRepo) GetUserFriends(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]uuid.UUID, error) {
	_ = ctx
	_ = userID
	_ = limit
	_ = offset

	return []uuid.UUID{mockUserID}, nil
}

// Mock metrics
type mockMittMetrics struct {
	FakeTotalMitts   int
	FakeTotalLikes   int
	FakeViewedInFeed float64
}

func (m *mockMittMetrics) AddMitt() {
	m.FakeTotalMitts++
}

func (m *mockMittMetrics) DeleteMitt() {
	m.FakeTotalMitts--
}

func (m *mockMittMetrics) AddLike() {
	m.FakeTotalLikes++
}

func (m *mockMittMetrics) DeleteLike() {
	m.FakeTotalLikes--
}

func (m *mockMittMetrics) ViewInFeed(count float64) {
	m.FakeViewedInFeed += count
}
