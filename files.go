package ryft

import (
	"context"
	"net/http"
)

type FilesService struct {
	client *Client
}

type File struct {
	ID               string `json:"id,omitempty"`
	Category         string `json:"category,omitempty"`
	CreatedTimestamp int    `json:"createdTimestamp,omitempty"`
}

type FileList struct {
	Items []File `json:"items"`
}

type CreateFileRequest struct {
	FilePath string
	Category string
	Account  string
}

type FileListParams struct {
	ListParams
	Category string
}

func (s *FilesService) List(ctx context.Context, params FileListParams) (*FileList, error) {
	query := buildListQuery(params.ListParams)
	if params.Category != "" {
		query.Set("category", params.Category)
	}

	req, err := s.client.newRequestWithQuery(ctx, http.MethodGet, "files", query, nil)
	if err != nil {
		return nil, err
	}

	var files FileList
	if err := s.client.doJSON(req, &files); err != nil {
		return nil, err
	}
	return &files, nil
}

func (s *FilesService) Get(ctx context.Context, fileID string) (*File, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "files/"+fileID, nil)
	if err != nil {
		return nil, err
	}

	var file File
	if err := s.client.doJSON(req, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

func (s *FilesService) Create(ctx context.Context, request CreateFileRequest) (*File, error) {
	category := request.Category
	if category == "" {
		category = "Evidence"
	}

	var file File
	if err := s.client.doMultipartFile(ctx, "files", request.Account, request.FilePath, category, &file); err != nil {
		return nil, err
	}
	return &file, nil
}
