package provider

import "context"

type Show struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Year   string   `json:"year"`
	Kind   string   `json:"kind"`
	Genres []string `json:"genres"`
}

type Episode struct {
	Season   int    `json:"season"`
	Number   int    `json:"number"`
	Title    string `json:"title"`
	Special  bool   `json:"special"`
	AirDate  string `json:"airdate"`
	Absolute int    `json:"absolute"`
}

type Provider interface {
	Name() string
	Configured() bool
	Search(context.Context, string) ([]Show, error)
	Episodes(context.Context, string) ([]Episode, error)
}

type Info struct {
	Name       string `json:"name"`
	Configured bool   `json:"configured"`
}

type Registry struct{ providers map[string]Provider }

func NewRegistry(providers ...Provider) *Registry {
	r := &Registry{providers: make(map[string]Provider)}
	for _, p := range providers {
		r.providers[p.Name()] = p
	}
	return r
}
func (r *Registry) Get(name string) (Provider, bool) { p, ok := r.providers[name]; return p, ok }
func (r *Registry) List() []Info {
	order := []string{"TVmaze", "TMDB", "AniDB"}
	out := make([]Info, 0, len(r.providers))
	for _, name := range order {
		if p, ok := r.providers[name]; ok {
			out = append(out, Info{Name: p.Name(), Configured: p.Configured()})
		}
	}
	return out
}
