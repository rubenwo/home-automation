package qr

const (
	DefaultLogoSize = 0.3
)

type OutputFormat string

const (
	PNG OutputFormat = "png"
	JPG OutputFormat = "jpg"
	SVG OutputFormat = "svg"
)

type config struct {
	logoSize   float64
	blendAlpha bool
}

type options func(c *config)

func WithLogoSize(logoSize float64) func(c *config) {
	return func(c *config) {
		if logoSize < DefaultLogoSize && logoSize > 0 {
			c.logoSize = logoSize
		}
	}
}

func WithAlphaBlending(c *config) {
	c.blendAlpha = true
}
