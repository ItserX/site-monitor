package storage

import "context"

type Site struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Active bool   `json:"active"`
}

type Storage interface {
	AddSite(ctx context.Context, site Site) (string, error)
	GetSites(ctx context.Context) ([]Site, error)
	GetSiteByID(ctx context.Context, id string) (*Site, error)
	UpdateSite(ctx context.Context, site Site) error
	DeleteSite(ctx context.Context, id string) error
}
