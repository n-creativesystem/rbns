package entity

import (
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

type Model struct {
	ID        string `gorm:"type:varchar(256);primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Model) Generate() {
	m.ID = generate()
}

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

func generate() string {
	g := p.Get().(*generator)
	id := g.New().String()
	p.Put(g)
	return id
}
