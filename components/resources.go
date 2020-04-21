package components

import "fmt"


type Resourcer interface{
	Prepare(map[string]interface{}) string
	Run() string
}

type MapOfResource map[string]func() (Resourcer, error)

func (mor MapOfResource) getResource(name string) (Resourcer, error){
	r, ok := mor[name]
	if !ok {
		fmt.Println("i am in get")
		return nil, fmt.Errorf(" InvalidResourceTYpe: Resource type %s not found", name)
	}

	return r()
}
