module github.com/riesinger/hugo-golunr

go 1.24.0

require (
	github.com/adrg/frontmatter v0.2.0
	github.com/grokify/html-strip-tags-go v0.1.0
	github.com/spf13/afero v1.15.0
	github.com/writeas/go-strip-markdown v2.0.1+incompatible
)

require (
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/riesinger/hugo-golunr/internal/post => ./internal/post
