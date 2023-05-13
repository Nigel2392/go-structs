package structs

type ValidatorMap map[string][]func(interface{}) error

func (v ValidatorMap) Add(field string, validator func(interface{}) error) {
	v[field] = append(v[field], validator)
}

func (v ValidatorMap) Set(field string, validator func(interface{}) error) {
	v[field] = []func(interface{}) error{validator}
}

func (v ValidatorMap) Remove(field string) {
	delete(v, field)
}

func (v ValidatorMap) Validate(field string, value interface{}) error {
	for _, validator := range v[field] {
		if err := validator(value); err != nil {
			return err
		}
	}
	return nil
}
