package nerinyan

type BeatmapSet struct {
	Artist             string             `json:"artist"`
	ArtistUnicode      string             `json:"artist_unicode"`
	Creator            string             `json:"creator"`
	FavouriteCount     int                `json:"favourite_count"`
	Hype               Hype               `json:"hype"`
	ID                 int                `json:"id"`
	NSFW               bool               `json:"nsfw"`
	PlayCount          int                `json:"play_count"`
	PreviewURL         string             `json:"preview_url"`
	Source             string             `json:"source"`
	Status             string             `json:"status"`
	Title              string             `json:"title"`
	TitleUnicode       string             `json:"title_unicode"`
	UserID             int                `json:"user_id"`
	Video              bool               `json:"video"`
	Availability       Availability       `json:"availability"`
	BPM                string             `json:"bpm"`
	CanBeHyped         bool               `json:"can_be_hyped"`
	DiscussionEnabled  bool               `json:"discussion_enabled"`
	DiscussionLocked   bool               `json:"discussion_locked"`
	IsScoreable        bool               `json:"is_scoreable"`
	LastUpdated        string             `json:"last_updated"`
	LegacyThreadURL    string             `json:"legacy_thread_url"`
	NominationsSummary NominationsSummary `json:"nominations_summary"`
	Ranked             int                `json:"ranked"`
	RankedDate         string             `json:"ranked_date"`
	Storyboard         bool               `json:"storyboard"`
	SubmittedDate      string             `json:"submitted_date"`
	Tags               string             `json:"tags"`
	HasFavourited      bool               `json:"has_favourited"`
	Beatmaps           []Beatmap          `json:"beatmaps"`
	Description        Description        `json:"description"`
	Genre              Genre              `json:"genre"`
	Language           Language           `json:"language"`
	RatingsString      string             `json:"ratings_string"`
}
type Beatmap struct {
	DifficultyRating float64 `json:"difficulty_rating"`
	ID               int     `json:"id"`
	Mode             string  `json:"mode"`
	Status           string  `json:"status"`
	TotalLength      int     `json:"total_length"`
	UserID           int     `json:"user_id"`
	Version          string  `json:"version"`
	Accuracy         int     `json:"accuracy"`
	AR               float64 `json:"ar"`
	BeatmapSetID     int     `json:"beatmapset_id"`
	BPM              string  `json:"bpm"`
	Convert          bool    `json:"convert"`
	CountCircles     int     `json:"count_circles"`
	CountSliders     int     `json:"count_sliders"`
	CountSpinners    int     `json:"count_spinners"`
	CS               float64 `json:"cs"`
	DeletedAt        string  `json:"deleted_at"`
	Drain            int     `json:"drain"`
	HitLength        int     `json:"hit_length"`
	IsScoreable      bool    `json:"is_scoreable"`
	LastUpdated      string  `json:"last_updated"`
	ModeInt          int     `json:"mode_int"`
	Passcount        int     `json:"passcount"`
	Playcount        int     `json:"playcount"`
	Ranked           int     `json:"ranked"`
	URL              string  `json:"url"`
	Checksum         string  `json:"checksum"`
	MaxCombo         int     `json:"max_combo"`
}

type Hype struct {
	Current  interface{} `json:"current"`
	Required interface{} `json:"required"`
}

type Availability struct {
	DownloadDisabled bool   `json:"download_disabled"`
	MoreInformation  string `json:"more_information"`
}

type NominationsSummary struct {
	Current  int `json:"current"`
	Required int `json:"required"`
}

type Description struct {
	Description string `json:"description"`
}

type Genre struct {
	ID   interface{} `json:"id"`
	Name interface{} `json:"name"`
}

type Language struct {
	ID   interface{} `json:"id"`
	Name interface{} `json:"name"`
}
