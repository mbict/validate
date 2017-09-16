package validate

import "regexp"

const (
	alphaRegexString        = "^[a-zA-Z]+$"
	alphaNumericRegexString = "^[a-zA-Z0-9]+$"
	alphaDashRegexString    = "[^\\d\\w-_]"
	alphaDashDotRegexString = "[^\\d\\w-_\\.]"
	numericRegexString      = "^[-+]?[0-9]+(?:\\.[0-9]+)?$"
	numberRegexString       = "^[0-9]+$"
	emailRegexString        = "^[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?$"
	base64RegexString       = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
	uUID3RegexString        = "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
	uUID4RegexString        = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	uUID5RegexString        = "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	uUIDRegexString         = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
	aSCIIRegexString        = "^[\x00-\x7F]*$"
	urlRegexString          = `^(http|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?$`
)

var (
	alphaRegex        = regexp.MustCompile(alphaRegexString)
	alphaNumericRegex = regexp.MustCompile(alphaNumericRegexString)
	alphaDashRegex    = regexp.MustCompile(alphaDashRegexString)
	alphaDashDotRegex = regexp.MustCompile(alphaDashDotRegexString)
	numericRegex      = regexp.MustCompile(numericRegexString)
	numberRegex       = regexp.MustCompile(numberRegexString)
	emailRegex        = regexp.MustCompile(emailRegexString)
	base64Regex       = regexp.MustCompile(base64RegexString)
	uUID3Regex        = regexp.MustCompile(uUID3RegexString)
	uUID4Regex        = regexp.MustCompile(uUID4RegexString)
	uUID5Regex        = regexp.MustCompile(uUID5RegexString)
	uUIDRegex         = regexp.MustCompile(uUIDRegexString)
	aSCIIRegex        = regexp.MustCompile(aSCIIRegexString)
	urlRegex          = regexp.MustCompile(urlRegexString)
)
