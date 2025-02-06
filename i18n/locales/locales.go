package locales

import "embed"

//go:embed en
//go:embed fr
//go:embed it
//go:embed es
//go:embed pt
var Content embed.FS
