package components

import "fmt"


type Resourcer interface{
	Prepare(map[string]interface{}) error
	Run() error
	Validate(map[string]interface{}) error
}

type MapOfResource map[string]func() (Resourcer, error)

func (mor MapOfResource) getResource(name string) (Resourcer, error){
	r, ok := mor[name]
	if !ok {
		return nil, fmt.Errorf("InvalidResourceTYpe: Resource type %s not found", name)
	}

	return r()
}
