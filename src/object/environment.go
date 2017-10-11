package object

type Environment struct {
	vars  map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{vars: make(map[string]Object)}
}

func (env *Environment) NewLinkedEnvironment() *Environment {
	res := NewEnvironment()
	res.outer = env
	return res
}

func (env *Environment) Get(name string) (result Object, ok bool) {
	result, ok = env.vars[name]
	if !ok && env.outer != nil {
		result, ok = env.outer.Get(name)
	}
	return
}

func (env *Environment) Set(name string, value Object) {
	env.vars[name] = value
}
