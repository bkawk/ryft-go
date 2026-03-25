package ryft

import (
	"context"
	"net/url"
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

func (s *FilesService) List(
	ctx context.Context,
	category string,
	ascending bool,
	limit int,
	startsAfter string,
) (*FileList, error) {
	query := url.Values{}
	if category != "" {
		query.Set("category", category)
	}
	query.Set("ascending", boolString(ascending))
	if limit > 0 {
		query.Set("limit", itoa(limit))
	}
	if startsAfter != "" {
		query.Set("startsAfter", startsAfter)
	}

	req, err := s.client.newRequestWithQuery(ctx, "GET", "files", query, nil)
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
	req, err := s.client.newRequest(ctx, "GET", "files/"+fileID, nil)
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
