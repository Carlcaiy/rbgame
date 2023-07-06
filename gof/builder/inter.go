package builder

type Item interface {
	Name() string
	Packing() Packing
	Price() int
}

type Packing interface {
	Pack() string
}

type Wrapper struct {
}

func (w *Wrapper) Pack() string {
	return "Wrapper"
}

type Bottle struct {
}

func (b *Bottle) Pack() string {
	return "Bottle"
}

type burger struct {
	Item
}

func (b *burger) Name() string {
	return b.Name()
}

func (b *burger) Price() int {
	return b.Price()
}

func (b *burger) Packing() Packing {
	return new(Wrapper)
}

type coldDrink struct {
	name  string
	price int
}

func (b *coldDrink) Name() string {
	return b.name
}

func (b *coldDrink) Price() int {
	return b.price
}

func (b *coldDrink) Packing() Packing {
	return new(Bottle)
}

type VegBurger struct {
	burger
}

func (v *VegBurger) Price() int {
	return 12
}

func (v *VegBurger) Name() string {
	return "VegBurger"
}

type ChickenBurger struct {
	burger
}

func (c *ChickenBurger) Price() int {
	return 19
}

func (c *ChickenBurger) Name() string {
	return "ChickenBurger"
}
