package base

type Option struct {
	value interface{}
}

func NewOption(value interface{}) Option {
	return Option{
		value: value,
	}
}

func (option Option) Value() interface{} {
	return option.value
}

type Options []Option

func (list Options) Filter(handler func(item Option) bool) Options {
	options := make(Options, 0, len(list))
	for _, item := range list {
		if handler(item) {
			options = append(options, item)
		}
	}
	return options
}

func (list Options) Values() []interface{} {
	values := make([]interface{}, 0, len(list))
	for _, option := range list {
		values = append(values, option.Value())
	}
	return values
}
