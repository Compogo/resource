package mountinfo

import "github.com/Compogo/resource/infrastructure/resource/cgroup/path"

type Option func(p *Parser)

func WithPath(path path.Path) Option {
	return func(p *Parser) {
		p.path = path
	}
}
