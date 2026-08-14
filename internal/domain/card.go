package domain

type ArcanaType string

const (
	Major ArcanaType = "Major"
	Minor ArcanaType = "Minor"
)

type Meanings struct {
	Upright  string `json:"upright"`
	Reversed string `json:"reversed"`
}

type Card struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	ArcanaType  ArcanaType `json:"arcana_type"`
	Suit        *string    `json:"suit"` // Pointer allows null for Major Arcana
	Value       int        `json:"value"`
	Meanings    Meanings   `json:"meanings"`
	Description string     `json:"description"`
}

type DrawResult struct {
	Card       Card
	IsReversed bool
	Meaning    string
	Position   *string
}
