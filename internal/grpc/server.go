package grpc

import (
	"context"
	"strconv"

	"github.com/tnphucccc/mangahub/internal/grpc/pb"
	"github.com/tnphucccc/mangahub/internal/manga"
	"github.com/tnphucccc/mangahub/pkg/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the gRPC MangaServiceServer interface.
type Server struct {
	pb.UnimplementedMangaServiceServer
	mangaService *manga.Service
}

// NewServer creates a new gRPC server.
func NewServer(mangaService *manga.Service) *Server {
	return &Server{
		mangaService: mangaService,
	}
}

// Convert response into MangaResponse
func toMangaResponse(m *models.Manga) *pb.MangaResponse {
	return &pb.MangaResponse{
		Id:            m.ID,
		Title:         m.Title,
		Author:        m.Author,
		Genres:        m.Genres,
		Status:        string(m.Status),
		TotalChapters: int32(m.TotalChapters),
		Description:   m.Description,
		CoverUrl:      m.CoverImageURL,
	}
}

// GetManga retrieves a manga by its ID.
func (s *Server) GetManga(ctx context.Context, req *pb.GetMangaRequest) (*pb.MangaResponse, error) {
	manga, err := s.mangaService.GetByID(req.MangaId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get manga: %v", err)
	}
	if manga == nil {
		return nil, status.Errorf(codes.NotFound, "Manga with ID %s not found", req.MangaId)
	}

	return toMangaResponse(manga), nil
}

// SearchManga searches for manga based on a query.
func (s *Server) SearchManga(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	query := models.MangaSearchQuery{
		Title:  req.GetTitle(),
		Author: req.GetAuthor(),
		Genre:  req.GetGenre(),
		Status: models.MangaStatus(req.GetStatus()),
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	}
	mangas, err := s.mangaService.Search(query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to search manga: %v", err)
	}

	var mangaResponses []*pb.MangaResponse
	for _, m := range mangas {
		mangaResponses = append(mangaResponses, toMangaResponse(&m))
	}

	return &pb.SearchResponse{Manga: mangaResponses}, nil
}

// UpdateProgress updates the user's reading progress for a manga.
func (s *Server) UpdateProgress(ctx context.Context, req *pb.UpdateProgressRequest) (*pb.UpdateProgressResponse, error) {
	statusVal := models.ReadingStatus(req.GetStatus())
	ratingVal := int(req.GetRating())

	request := models.ProgressUpdateRequest{
		Status:         &statusVal,
		Rating:         &ratingVal,
		CurrentChapter: int(req.GetChapter()),
	}

	err := s.mangaService.UpdateProgress(req.GetUserId(), req.GetMangaId(), request)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update progress: %v", err)
	}

	progress, err := s.mangaService.GetProgress(req.GetUserId(), req.GetMangaId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to retrieve updated progress: %v", err)
	}

	progressResponse := &pb.UserProgress{
		UserId:         progress.UserID,
		MangaId:        progress.MangaID,
		CurrentChapter: int32(progress.CurrentChapter),
		Status:         string(progress.Status),
		Rating:         int32(progress.GetRatingValue()),
		StartedAt:      strconv.FormatInt(progress.GetStartedAtValue().Unix(), 10),
		CompletedAt:    strconv.FormatInt(progress.GetCompletedAtValue().Unix(), 10),
		UpdatedAt:      strconv.FormatInt(progress.UpdatedAt.Unix(), 10),
	}

	return &pb.UpdateProgressResponse{
		Progress: progressResponse,
	}, nil
}
