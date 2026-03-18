package converter

import (
	"database/sql"
	"testing"
	"time"

	"github.com/lm379/hmx/config"
	"github.com/lm379/hmx/internal/models"
	"github.com/pgvector/pgvector-go"
)

func setupConverterTestConfig() {
	config.AppConfig = config.Config{
		S3CustomDomain: "https://cdn.example.com",
		S3Endpoint:     "https://s3.amazonaws.com",
		S3Bucket:       "my-bucket",
	}
}

func TestToOperaResponse(t *testing.T) {
	setupConverterTestConfig()

	now := time.Now()

	tests := []struct {
		name  string
		opera models.Opera
	}{
		{
			name: "Complete opera with all fields",
			opera: models.Opera{
				OperaID:     1,
				OperaTitle:  "Test Opera",
				VideoPath:   "videos/test.mp4",
				Description: "Test description",
				AiSummary:   "AI summary",
				IsHidden:    false,
				ReleaseDate: sql.NullTime{Time: now, Valid: true},
				Duration:    sql.NullString{String: "01:30:00", Valid: true},
				MusicPath:   sql.NullString{String: "music/test.mp3", Valid: true},
				SrtPath:     sql.NullString{String: "srt/test.srt", Valid: true},
				Avatar:      sql.NullString{String: "avatars/test.jpg", Valid: true},
				CreatedAt:   now,
				UpdatedAt:   now,
				Embedding:   pgvector.NewVector(make([]float32, 4096)),
				Artists: []*models.Artist{
					{
						ArtistID:     1,
						Name:         "Artist 1",
						ArtistAvatar: sql.NullString{String: "avatars/artist1.jpg", Valid: true},
					},
				},
			},
		},
		{
			name: "Opera with minimal fields",
			opera: models.Opera{
				OperaID:    2,
				OperaTitle: "Simple Opera",
				VideoPath:  "videos/simple.mp4",
				Description: "",
				AiSummary:  "",
				IsHidden:   false,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		},
		{
			name: "Hidden opera",
			opera: models.Opera{
				OperaID:    3,
				OperaTitle: "Hidden Opera",
				VideoPath:  "videos/hidden.mp4",
				IsHidden:   true,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		},
		{
			name: "Opera with multiple artists",
			opera: models.Opera{
				OperaID:    4,
				OperaTitle: "Multi-Artist Opera",
				VideoPath:  "videos/multi.mp4",
				CreatedAt:  now,
				UpdatedAt:  now,
				Artists: []*models.Artist{
					{ArtistID: 1, Name: "Artist 1"},
					{ArtistID: 2, Name: "Artist 2"},
					{ArtistID: 3, Name: "Artist 3"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToOperaResponse(tt.opera)

			if result == nil {
				t.Fatal("ToOperaResponse() returned nil")
			}

			if result.OperaID != tt.opera.OperaID {
				t.Errorf("OperaID mismatch: got %v, want %v", result.OperaID, tt.opera.OperaID)
			}
			if result.OperaTitle != tt.opera.OperaTitle {
				t.Errorf("OperaTitle mismatch: got %v, want %v", result.OperaTitle, tt.opera.OperaTitle)
			}

			// Check that URL is converted to full URL
			if !isFullURL(result.VideoPath) {
				t.Errorf("VideoPath should be a full URL, got %v", result.VideoPath)
			}
		})
	}
}

func TestToOperaResponseList(t *testing.T) {
	setupConverterTestConfig()

	operas := []models.Opera{
		{
			OperaID:    1,
			OperaTitle: "Opera 1",
			VideoPath:  "videos/1.mp4",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			OperaID:    2,
			OperaTitle: "Opera 2",
			VideoPath:  "videos/2.mp4",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			OperaID:    3,
			OperaTitle: "Opera 3",
			VideoPath:  "videos/3.mp4",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	results := ToOperaResponseList(operas)

	if len(results) != len(operas) {
		t.Errorf("ToOperaResponseList() returned %d results, want %d", len(results), len(operas))
	}

	for i, result := range results {
		if result.OperaID != operas[i].OperaID {
			t.Errorf("OperaID at index %d mismatch: got %v, want %v", i, result.OperaID, operas[i].OperaID)
		}
	}
}

func TestToArtistResponse(t *testing.T) {
	setupConverterTestConfig()

	now := time.Now()

	tests := []struct {
		name   string
		artist models.Artist
	}{
		{
			name: "Complete artist",
			artist: models.Artist{
				ArtistID:     1,
				Name:         "Test Artist",
				Bio:          "Test bio",
				ArtistAvatar: sql.NullString{String: "avatars/artist.jpg", Valid: true},
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
		{
			name: "Artist without avatar",
			artist: models.Artist{
				ArtistID:  2,
				Name:      "Artist No Avatar",
				Bio:       "Bio",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "Artist with operas",
			artist: models.Artist{
				ArtistID: 3,
				Name:     "Artist With Operas",
				Bio:      "Bio",
				CreatedAt: now,
				UpdatedAt: now,
				Operas: []*models.Opera{
					{
						OperaID:    1,
						OperaTitle: "Opera 1",
						VideoPath:  "videos/1.mp4",
						CreatedAt:  now,
						UpdatedAt:  now,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToArtistResponse(tt.artist)

			if result == nil {
				t.Fatal("ToArtistResponse() returned nil")
			}

			if result.ArtistID != tt.artist.ArtistID {
				t.Errorf("ArtistID mismatch: got %v, want %v", result.ArtistID, tt.artist.ArtistID)
			}
			if result.Name != tt.artist.Name {
				t.Errorf("Name mismatch: got %v, want %v", result.Name, tt.artist.Name)
			}
			if result.Bio != tt.artist.Bio {
				t.Errorf("Bio mismatch: got %v, want %v", result.Bio, tt.artist.Bio)
			}

			// Check avatar URL
			if tt.artist.ArtistAvatar.Valid && result.Avatar != nil {
				if !isFullURL(*result.Avatar) {
					t.Errorf("Avatar should be a full URL, got %v", *result.Avatar)
				}
			}
		})
	}
}

func TestToArtistResponseList(t *testing.T) {
	setupConverterTestConfig()

	artists := []models.Artist{
		{
			ArtistID:  1,
			Name:      "Artist 1",
			Bio:       "Bio 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ArtistID:  2,
			Name:      "Artist 2",
			Bio:       "Bio 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	results := ToArtistResponseList(artists)

	if len(results) != len(artists) {
		t.Errorf("ToArtistResponseList() returned %d results, want %d", len(results), len(artists))
	}

	for i, result := range results {
		if result.ArtistID != artists[i].ArtistID {
			t.Errorf("ArtistID at index %d mismatch: got %v, want %v", i, result.ArtistID, artists[i].ArtistID)
		}
	}
}

func TestToUserResponse(t *testing.T) {
	setupConverterTestConfig()

	now := time.Now()
	email := "test@example.com"

	tests := []struct {
		name string
		user *models.Users
	}{
		{
			name: "Complete user",
			user: &models.Users{
				UserID:      1,
				Username:    "testuser",
				Phone:       "13800138000",
				Email:       sql.NullString{String: email, Valid: true},
				Sex:         models.Male,
				Icon:        sql.NullString{String: "avatars/user.jpg", Valid: true},
				Role:        models.User,
				LastLoginAt: &now,
				LastIp:      "192.168.1.1",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
		{
			name: "User with minimal fields",
			user: &models.Users{
				UserID:    2,
				Username:  "simpleuser",
				Phone:     "13900139000",
				Sex:       models.Female,
				Role:      models.User,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToUserResponse(tt.user)

			if result == nil {
				t.Fatal("ToUserResponse() returned nil")
			}

			if result.UserID != tt.user.UserID {
				t.Errorf("UserID mismatch: got %v, want %v", result.UserID, tt.user.UserID)
			}
			if result.Username != tt.user.Username {
				t.Errorf("Username mismatch: got %v, want %v", result.Username, tt.user.Username)
			}

			// Check icon URL
			if tt.user.Icon.Valid && result.Icon != nil {
				if !isFullURL(*result.Icon) {
					t.Errorf("Icon should be a full URL, got %v", *result.Icon)
				}
			}
		})
	}
}

func TestToUserResponseNil(t *testing.T) {
	result := ToUserResponse(nil)
	if result != nil {
		t.Error("ToUserResponse() with nil user should return nil")
	}
}

func TestToOperaListResponse(t *testing.T) {
	setupConverterTestConfig()

	opera := models.Opera{
		OperaID:    1,
		OperaTitle: "Test Opera",
		VideoPath:  "videos/test.mp4",
		IsHidden:   false,
		CreatedAt:  time.Now(),
	}

	result := ToOperaListResponse(opera)

	if result == nil {
		t.Fatal("ToOperaListResponse() returned nil")
	}

	if result.OperaID != opera.OperaID {
		t.Errorf("OperaID mismatch: got %v, want %v", result.OperaID, opera.OperaID)
	}
}

func TestToOperaDetailResponse(t *testing.T) {
	setupConverterTestConfig()

	now := time.Now()

	opera := models.Opera{
		OperaID:     1,
		OperaTitle:  "Test Opera",
		VideoPath:   "videos/test.mp4",
		Description: "Description",
		AiSummary:   "Summary",
		IsHidden:    false,
		ReleaseDate: sql.NullTime{Time: now, Valid: true},
		Duration:    sql.NullString{String: "01:30:00", Valid: true},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result := ToOperaDetailResponse(opera)

	if result == nil {
		t.Fatal("ToOperaDetailResponse() returned nil")
	}

	if result.OperaID != opera.OperaID {
		t.Errorf("OperaID mismatch: got %v, want %v", result.OperaID, opera.OperaID)
	}
}

func TestToArtistDetailResponse(t *testing.T) {
	setupConverterTestConfig()

	now := time.Now()

	artist := models.Artist{
		ArtistID: 1,
		Name:     "Test Artist",
		Bio:      "Bio",
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ToArtistDetailResponse(artist)

	if result == nil {
		t.Fatal("ToArtistDetailResponse() returned nil")
	}

	if result.ArtistID != artist.ArtistID {
		t.Errorf("ArtistID mismatch: got %v, want %v", result.ArtistID, artist.ArtistID)
	}
}

func TestToSimpleOpera(t *testing.T) {
	setupConverterTestConfig()

	now := time.Now()

	opera := models.Opera{
		OperaID:    1,
		OperaTitle: "Test Opera",
		CreatedAt:  now,
	}

	result := ToSimpleOpera(opera)

	if result.OperaID != opera.OperaID {
		t.Errorf("OperaID mismatch: got %v, want %v", result.OperaID, opera.OperaID)
	}
	if result.OperaTitle != opera.OperaTitle {
		t.Errorf("OperaTitle mismatch: got %v, want %v", result.OperaTitle, opera.OperaTitle)
	}
}

func TestToSimpleOperaList(t *testing.T) {
	setupConverterTestConfig()

	operas := []models.Opera{
		{OperaID: 1, OperaTitle: "Opera 1", CreatedAt: time.Now()},
		{OperaID: 2, OperaTitle: "Opera 2", CreatedAt: time.Now()},
	}

	results := ToSimpleOperaList(operas)

	if len(results) != len(operas) {
		t.Errorf("ToSimpleOperaList() returned %d results, want %d", len(results), len(operas))
	}
}

func TestConverterURLHandling(t *testing.T) {
	setupConverterTestConfig()

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Relative path",
			path:     "videos/test.mp4",
			expected: "https://cdn.example.com/videos/test.mp4",
		},
		{
			name:     "Full URL",
			path:     "https://other-cdn.com/video.mp4",
			expected: "https://other-cdn.com/video.mp4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opera := models.Opera{
				OperaID:    1,
				OperaTitle: "Test",
				VideoPath:  tt.path,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			result := ToOperaResponse(opera)

			if result.VideoPath != tt.expected {
				t.Errorf("VideoPath = %v, want %v", result.VideoPath, tt.expected)
			}
		})
	}
}

// Helper function

func isFullURL(url string) bool {
	return len(url) > 0 && (url[:7] == "http://" || url[:8] == "https://")
}
