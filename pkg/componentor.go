package pkg

type Componentor interface {
	Use(driver string) Componentor
	Config(configuration map[string]interface{}) error
	Component() interface{}
}
