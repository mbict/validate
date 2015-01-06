package validate

import . "gopkg.in/check.v1"

type validateSuite struct{}

var _ = Suite(&validateSuite{})

func (s *validateSuite) Test_NoErrors(c *C) {
	v := NewValidator()
	errs := v.Validate(BlogPost{
		Id: 3,
		Post: Post{
			Title:   "Behold The Title!",
			Content: "And some content",
		},
		Author: Person{
			Name: "Matt Holt",
		},
	})
	c.Assert(errs, IsNil)
}

func (s *validateSuite) Test_IdRequired(c *C) {
	v := NewValidator()
	errs := v.Validate(BlogPost{
		Post: Post{
			Title:   "Behold The Title!",
			Content: "And some content",
		},
		Author: Person{
			Name: "Matt Holt",
		},
	})

	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{Error{FieldNames: []string{"Id"}, Classification: RequiredError, Message: "Required"}})
}

func (s *validateSuite) Test_SlicesRequired(c *C) {
	v := NewValidator()
	errs := v.Validate(BlogPost{
		Id:      1,
		Ratings: []int{1, 2, 3, 4, 5, 6},
		Post: Post{
			Title:   "Behold The Title!",
			Content: "And some content",
		},
		Author:   Person{},
		Coauthor: &Person{},
		Readers: []Person{
			Person{
				Name:  "Michael Boke",
				Email: "me@email.com",
			},
			Person{
				Email: "you@email.com",
			},
			Person{},
		},
		Contributors: []*Person{
			&Person{
				Name:  "Michael Boke",
				Email: "me@email.com",
			},
			&Person{
				Email: "you@email.com",
			},
			&Person{},
		},
	})

	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{FieldNames: []string{"Rating"}, Classification: MaxError, Message: "Max"},
		Error{FieldNames: []string{"Author.Name"}, Classification: RequiredError, Message: "Required"},
		Error{FieldNames: []string{"Coauthor.Name"}, Classification: RequiredError, Message: "Required"},
		Error{FieldNames: []string{"Readers.1.Name"}, Classification: RequiredError, Message: "Required"},
		Error{FieldNames: []string{"Readers.2.Name"}, Classification: RequiredError, Message: "Required"},
		Error{FieldNames: []string{"Contributors.1.Name"}, Classification: RequiredError, Message: "Required"},
		Error{FieldNames: []string{"Contributors.2.Name"}, Classification: RequiredError, Message: "Required"},
	})
}

func (s *validateSuite) Test_EmbeddedStructFieldRequired(c *C) {
	v := NewValidator()
	errs := v.Validate(BlogPost{
		Id: 1,
		Post: Post{
			Content: "Content given, but title is required",
		},
		Author: Person{
			Name: "Matt Holt",
		},
	})
	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{
			FieldNames:     []string{"Title"},
			Classification: RequiredError,
			Message:        "Required",
		},
		Error{
			FieldNames:     []string{"Title"},
			Classification: "LengthError",
			Message:        "Life is too short",
		},
	})
}

func (s *validateSuite) Test_NestedStructFieldRequired(c *C) {
	v := NewValidator()
	errs := v.Validate(BlogPost{
		Id: 1,
		Post: Post{
			Title:   "Behold The Title!",
			Content: "And some content",
		},
	})
	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{
			FieldNames:     []string{"Author.Name"},
			Classification: RequiredError,
			Message:        "Required",
		},
	})
}

func (s *validateSuite) Test_RequiredFieldMissingInNestedStructPointer(c *C) {
	v := NewValidator()
	errs := v.Validate(BlogPost{
		Id: 1,
		Post: Post{
			Title:   "Behold The Title!",
			Content: "And some content",
		},
		Author: Person{
			Name: "Matt Holt",
		},
		Coauthor: &Person{},
	})
	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{
			FieldNames:     []string{"Coauthor.Name"},
			Classification: RequiredError,
			Message:        "Required",
		},
	})
}

func (s *validateSuite) Test_AllRequiredFieldsSpecifiedInNestedStructPointer(c *C) {
	v := NewValidator()
	errs := v.Validate(BlogPost{
		Id: 1,
		Post: Post{
			Title:   "Behold The Title!",
			Content: "And some content",
		},
		Author: Person{
			Name: "Matt Holt",
		},
		Coauthor: &Person{
			Name: "Jeremy Saenz",
		},
	})
	c.Assert(errs, IsNil)
}

func (s *validateSuite) Test_CustomStructValidation(c *C) {
	v := NewValidator()
	errs := v.Validate(BlogPost{
		Id: 1,
		Post: Post{
			Title:   "Too short",
			Content: "And some content",
		},
		Author: Person{
			Name: "Matt Holt",
		},
	})
	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{
			FieldNames:     []string{"Title"},
			Classification: "LengthError",
			Message:        "Life is too short",
		},
	})
}

func (s *validateSuite) Test_ListValidation(c *C) {
	v := NewValidator()
	errs := v.Validate([]BlogPost{
		BlogPost{
			Id: 1,
			Post: Post{
				Title:   "First Post",
				Content: "And some content",
			},
			Author: Person{
				Name: "Leeor Aharon",
			},
		},
		BlogPost{
			Id: 2,
			Post: Post{
				Title:   "Second Post",
				Content: "And some content",
			},
			Author: Person{
				Name: "Leeor Aharon",
			},
		},
	})
	c.Assert(errs, IsNil)
}

func (s *validateSuite) Test_ListValidationErrors(c *C) {
	v := NewValidator()
	errs := v.Validate([]BlogPost{
		BlogPost{
			Id: 1,
			Post: Post{
				Title:   "First Post",
				Content: "And some content",
			},
			Author: Person{
				Name: "Leeor Aharon",
			},
		},
		BlogPost{
			Id: 2,
			Post: Post{
				Title:   "Too Short",
				Content: "And some content",
			},
			Author: Person{
				Name: "Leeor Aharon",
			},
		},
	})
	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{
			FieldNames:     []string{"Title"},
			Classification: "LengthError",
			Message:        "Life is too short",
		},
	})
}

func (s *validateSuite) Test_InvalidCustomValidations(c *C) {
	v := NewValidator()
	errs := v.Validate([]SadForm{
		SadForm{
			AlphaDash:    ",",
			AlphaDashDot: ",",
			MinSize:      ",",
			MinSizeSlice: []string{",", ","},
			MaxSize:      ",,",
			MaxSizeSlice: []string{",", ","},
			Email:        ",",
			Url:          ",",
			UrlEmpty:     "",
			Range:        3,
			In:           "2",
			InInvalid:    "4",
			NotIn:        "1",
			Include:      "def",
			Exclude:      "abc",
		},
	})
	c.Assert(errs, NotNil)
	c.Assert(errs, DeepEquals, Errors{
		Error{
			FieldNames:     []string{"0.AlphaDash"},
			Classification: AlphaDashError,
			Message:        "AlphaDash",
		},
		Error{
			FieldNames:     []string{"0.AlphaDashDot"},
			Classification: AlphaDashDotError,
			Message:        "AlphaDashDot",
		},
		Error{
			FieldNames:     []string{"0.MinSize"},
			Classification: MinError,
			Message:        "Min",
		},
		Error{
			FieldNames:     []string{"0.MinSizeSlice"},
			Classification: MinError,
			Message:        "Min",
		},
		Error{
			FieldNames:     []string{"0.MaxSize"},
			Classification: MaxError,
			Message:        "Max",
		},
		Error{
			FieldNames:     []string{"0.MaxSizeSlice"},
			Classification: MaxError,
			Message:        "Max",
		},
		Error{
			FieldNames:     []string{"0.Email"},
			Classification: EmailError,
			Message:        "Email",
		},
		Error{
			FieldNames:     []string{"0.Url"},
			Classification: UrlError,
			Message:        "Url",
		},
		Error{
			FieldNames:     []string{"0.Range"},
			Classification: RangeError,
			Message:        "Range",
		},
		Error{
			FieldNames:     []string{"0.InInvalid"},
			Classification: InError,
			Message:        "In",
		},
		Error{
			FieldNames:     []string{"0.NotIn"},
			Classification: NotInError,
			Message:        "NotIn",
		},
		Error{
			FieldNames:     []string{"0.Include"},
			Classification: IncludeError,
			Message:        "Include",
		},
		Error{
			FieldNames:     []string{"0.Exclude"},
			Classification: ExcludeError,
			Message:        "Exclude",
		},
	})
}

func (s *validateSuite) Test_ListOfValidCustomValidations(c *C) {
	v := NewValidator()
	errs := v.Validate([]SadForm{
		SadForm{
			AlphaDash:    "123-456",
			AlphaDashDot: "123.456",
			MinSize:      "12345",
			MinSizeSlice: []string{"1", "2", "3", "4", "5"},
			MaxSize:      "1",
			MaxSizeSlice: []string{"1"},
			Email:        "123@456.com",
			Url:          "http://123.456",
			Range:        2,
			In:           "1",
			InInvalid:    "1",
			Include:      "abc",
		},
	})
	c.Assert(errs, IsNil)
}
