package plugins

import (
	"context"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var p *sync.Pool

func init() {
	p = &sync.Pool{
		New: func() interface{} {
			return &generator{r: ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)}
		},
	}
}

type generator struct {
	r io.Reader
}

func (g *generator) New() ulid.ULID {
	return ulid.MustNew(ulid.Timestamp(time.Now()), g.r)
}

func Generate() string {
	g := p.Get().(*generator)
	id := g.New().String()
	p.Put(g)
	return id
}

type ID string

func (id ID) String() string {
	return string(id)
}

func (id ID) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	value := id.String()
	return clause.Expr{
		SQL:  "?",
		Vars: []interface{}{value},
	}
}

func (id *ID) Generate() {
	if id.String() == "" {
		*id = ID(Generate())
	}
}
