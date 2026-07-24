package media

type Preset struct {
	Name     string `json:"name"`
	Template string `json:"template"`
}

var Presets = []Preset{
	{Name: "Plex", Template: "{n} ({y})/Season {s}/{n} - {s00e00} - {t}"},
	{Name: "Kodi", Template: "{n} ({y})/Season {s}/{n} S{s00}E{e00} {t}"},
	{Name: "Jellyfin", Template: "{n} ({y})/Season {s}/{n} S{s00}E{e00} {t}"},
	{Name: "Simple", Template: "{n} - {s00e00} - {t}"},
}

func PresetByName(name string) (Preset, bool) {
	for _, p := range Presets {
		if p.Name == name {
			return p, true
		}
	}
	return Preset{}, false
}
