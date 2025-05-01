package mitt

import (
	"context"
	"github.com/misshanya/mitter/internal/models"
	"testing"
)

// Tests
func TestMittService_CreateMitt(t *testing.T) {
	service := NewService(mockMittRepo{})
	ctx := context.Background()

	mitt, err := service.CreateMitt(ctx, mockUserID, &models.MittCreate{
		Content: mockMittModel.Content,
	})
	if err != nil {
		t.Fatal(err)
	}

	if mitt != mockMittModel {
		t.Fatal("mitt does not match")
	}
}

func TestMittService_GetMitt(t *testing.T) {
	service := NewService(mockMittRepo{})
	ctx := context.Background()

	mitt, err := service.GetMitt(ctx, mockMittModel.ID)
	if err != nil {
		t.Fatal(err)
	}

	if mitt != mockMittModel {
		t.Fatal("mitt does not match")
	}
}

func TestMittService_GetAllUserMitts(t *testing.T) {
	service := NewService(mockMittRepo{})
	ctx := context.Background()

	mitts, err := service.GetAllUserMitts(ctx, mockUserID, 1, 0)
	if err != nil {
		t.Fatal(err)
	}

	excepted := []*models.Mitt{mockMittModel}
	if mitts[0] != excepted[0] {
		t.Fatal("mitts does not match")
	}
}

func TestMittService_UpdateMitt(t *testing.T) {
	service := NewService(mockMittRepo{})
	ctx := context.Background()

	mitt, err := service.UpdateMitt(ctx, mockUserID, mockMittModel.ID, &models.MittUpdate{
		Content: mockMittModel.Content + " (updated)",
	})
	if err != nil {
		t.Fatal(err)
	}

	if mitt != mockMittModel {
		t.Fatal("mitt does not match")
	}
}

func TestMittService_DeleteMitt(t *testing.T) {
	service := NewService(mockMittRepo{})
	ctx := context.Background()

	err := service.DeleteMitt(ctx, mockUserID, mockMittModel.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMittService_SwitchLike(t *testing.T) {
	service := NewService(mockMittRepo{})
	ctx := context.Background()

	// Like mitt
	newState, err := service.SwitchLike(ctx, mockUserID, mockMittModel.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !newState {
		t.Fatal("liked state should be true")
	}

	// Remove like
	newState, err = service.SwitchLike(ctx, mockUserID, mockMittModel.ID)
	if err != nil {
		t.Fatal(err)
	}

	if newState {
		t.Fatal("liked state should be false")
	}
}
