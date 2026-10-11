package assets

import _ "embed"

// CardsJSON is the full deck, compiled into the binary so it runs from any working directory.
//
//go:embed data/cards.json
var CardsJSON []byte
