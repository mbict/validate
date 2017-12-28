package examples

//go:generate validate -type Example1

type Example1 struct {
	ExampleEmbedded
	Id              int       `validate:"required"`
	Name            string    `validate:"required"`
	SliceSet        []string  `validate:"min(10)"`
	SlicePtrSet     []*string `validate:"min(10)"`
	SingleStruct    ExampleStructValidate
	SinglePtrStruct *ExampleStructValidate
	SliceStruct     []ExampleStructValidate
	SlicePtrStruct  []*ExampleStructValidate
}

func (e *Example1) Test() string {
	return "test"
}

func (e Example1) Validate() error {
	return nil
}

type ExampleEmbedded struct {
	Tags    []string `validate:"alphanum"`
	Keyword string   `validate:"omitEmpty;min(3)"`
}

type ExampleStructValidate struct {
}

func (e ExampleStructValidate) Validate() error {
	return nil
}
