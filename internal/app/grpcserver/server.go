package grpcserver

import (
	"context"
	"errors"

	"github.com/shekshuev/shortener/internal/app/models"
	"github.com/shekshuev/shortener/internal/app/proto"
	"github.com/shekshuev/shortener/internal/app/service"
	"github.com/shekshuev/shortener/internal/app/store"
)

// Server реализует gRPC-сервис URLShortenerServer.
type Server struct {
	proto.UnimplementedURLShortenerServer
	service service.Service
}

// NewServer создаёт новый экземпляр gRPC-сервера с переданным сервисом.
func NewServer(service service.Service) *Server {
	return &Server{
		service: service,
	}
}

// Shorten обрабатывает сокращение одного URL.
// Запрос: ShortenRequest { url, user_id }.
// Ответ: ShortenResponse { result: короткий URL } или ошибка.
func (s *Server) Shorten(ctx context.Context, req *proto.ShortenRequest) (*proto.ShortenResponse, error) {
	shortURL, err := s.service.CreateShortURL(ctx, req.Url, req.UserId)
	if err != nil {
		return nil, err
	}
	return &proto.ShortenResponse{Result: shortURL}, nil
}

// BatchShorten обрабатывает сокращение нескольких URL за один запрос.
// Запрос: BatchShortenRequest с массивом URL.
// Ответ: BatchShortenResponse с массивом результатов или ошибка.
func (s *Server) BatchShorten(ctx context.Context, req *proto.BatchShortenRequest) (*proto.BatchShortenResponse, error) {
	createDTOs := make([]models.BatchShortURLCreateDTO, len(req.Items))
	for i, item := range req.Items {
		createDTOs[i] = models.BatchShortURLCreateDTO{
			CorrelationID: item.CorrelationId,
			OriginalURL:   item.OriginalUrl,
		}
	}
	readDTOs, err := s.service.BatchCreateShortURL(ctx, createDTOs, req.UserId)
	if err != nil {
		if !errors.Is(err, store.ErrAlreadyExists) {
			return nil, err
		}
	}
	items := make([]*proto.BatchShortenResponseItem, len(readDTOs))
	for i, dto := range readDTOs {
		items[i] = &proto.BatchShortenResponseItem{
			CorrelationId: dto.CorrelationID,
			ShortUrl:      dto.ShortURL,
		}
	}
	return &proto.BatchShortenResponse{Items: items}, nil
}

// GetUserURLs возвращает все сокращённые URL пользователя.
// Запрос: UserURLsRequest { user_id }.
// Ответ: UserURLsResponse с массивом ссылок или ошибка.
func (s *Server) GetUserURLs(ctx context.Context, req *proto.UserURLsRequest) (*proto.UserURLsResponse, error) {
	readDTO, err := s.service.GetUserURLs(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	items := make([]*proto.UserURLItem, len(readDTO))
	for i, dto := range readDTO {
		items[i] = &proto.UserURLItem{
			ShortUrl:    dto.ShortURL,
			OriginalUrl: dto.OriginalURL,
		}
	}
	return &proto.UserURLsResponse{Urls: items}, nil
}

// DeleteUserURLs удаляет список URL пользователя асинхронно.
// Запрос: DeleteURLsRequest { user_id, short_urls }.
// Ответ: DeleteURLsResponse без тела.
func (s *Server) DeleteUserURLs(ctx context.Context, req *proto.DeleteURLsRequest) (*proto.DeleteURLsResponse, error) {
	go s.service.DeleteURLs(ctx, req.UserId, req.ShortUrls)
	return &proto.DeleteURLsResponse{}, nil
}

// Ping проверяет доступность базы данных.
// Запрос: PingRequest.
// Ответ: PingResponse при успешном подключении или ошибка.
func (s *Server) Ping(ctx context.Context, _ *proto.PingRequest) (*proto.PingResponse, error) {
	err := s.service.CheckDBConnection(ctx)
	if err != nil {
		return nil, err
	}
	return &proto.PingResponse{}, nil
}

// GetStats возвращает статистику сервиса: количество URL и пользователей.
// Запрос: StatsRequest.
// Ответ: StatsResponse с количеством URL и пользователей или ошибка.
func (s *Server) GetStats(ctx context.Context, _ *proto.StatsRequest) (*proto.StatsResponse, error) {
	stats, err := s.service.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	return &proto.StatsResponse{
		Urls:  int32(stats.URLs),
		Users: int32(stats.Users),
	}, nil
}

// GetOriginalURL возвращает оригинальный URL по его сокращённой форме.
// Запрос: GetOriginalURLRequest { short_url }.
// Ответ: GetOriginalURLResponse с оригинальной ссылкой или ошибка.
func (s *Server) GetOriginalURL(ctx context.Context, req *proto.GetOriginalURLRequest) (*proto.GetOriginalURLResponse, error) {
	longURL, err := s.service.GetLongURL(ctx, req.ShortUrl)
	if err != nil {
		if err == store.ErrAlreadyDeleted {
			return nil, store.ErrAlreadyDeleted
		}
		return nil, err
	}
	return &proto.GetOriginalURLResponse{OriginalUrl: longURL}, nil
}
