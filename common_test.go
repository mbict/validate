package validate

import (
	"mime/multipart"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

// These types are mostly contrived examples, but they're used
// across many test cases. The idea is to cover all the scenarios
// that this binding package might encounter in actual use.
type (
	// For basic test cases with a required field
	Post struct {
		Title   string `validate:"Required"`
		Content string
	}

	// To be used as a nested struct (with a required field)
	Person struct {
		Name  string `validate:"Required"`
		Email string
	}

	// For advanced test cases: multiple values, embedded
	// and nested structs, an ignored field, and single
	// and multiple file uploads
	BlogPost struct {
		Post
		Id           int `validate:"Required"`
		Ignored      string
		Ratings      []int
		Author       Person
		Coauthor     *Person
		Readers      []Person
		Contributors []*Person
		HeaderImage  *multipart.FileHeader
		Pictures     []*multipart.FileHeader
		unexported   string
	}

	EmbedPerson struct {
		*Person
	}

	SadForm struct {
		AlphaDash    string   `validate:"AlphaDash"`
		AlphaDashDot string   `validate:"AlphaDashDot"`
		MinSize      string   `validate:"Min(5)"`
		MinSizeSlice []string `validate:"Min(5)"`
		MaxSize      string   `validate:"Max(1)"`
		MaxSizeSlice []string `validate:"Max(1)"`
		Email        string   `validate:"Email"`
		Url          string   `validate:"Url"`
		UrlEmpty     string   `validate:"Url"`
		Range        int      `validate:"Range(1,2)"`
		RangeInvalid int      `validate:"Range(1)"`
		In           string   `validate:"In(1,2,3)"`
		InInvalid    string   `validate:"In(1,2,3)"`
		NotIn        string   `validate:"NotIn(1,2,3)"`
		Include      string   `validate:"Include(a)"`
		Exclude      string   `validate:"Exclude(a)"`
	}
)

func (p Post) Validate(errors Errors) Errors {
	if len(p.Title) < 10 {
		errors.Add([]string{"Title"}, "LengthError", "Life is too short")
	}
	return errors
}
