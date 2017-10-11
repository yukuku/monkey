package object

type Environment struct {
	vars map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{vars: make(map[string]Object)}
}

func (env *Environment) Get(name string) (result Object, ok bool) {
	result, ok = env.vars[name]
	return
}

func (env *Environment) Set(name string, value Object) {
	env.vars[name] = value
}
