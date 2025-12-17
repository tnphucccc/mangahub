package grpc

import (
	"context"

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
		// Query: req.Query,
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
