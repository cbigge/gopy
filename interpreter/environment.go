package interpreter

type Environment struct {
	env map[string]Item
}

func NewEnv() *Environment {
	e := make(map[string]Item)
	return &Environment{env: e}
}

func (e *Environment) Get(k string) (Item, bool) {
	val, ok := e.env[k]
	return val, ok
}

func (e *Environment) Store(k string, i Item) Item {
	e.env[k] = i
	return i
}