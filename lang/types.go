package lang

type Types int

const (
	TypeInt32 Types = iota
	TypeString
)

var TypeKeywords = map[Keyword]Types{
	KeywordInt32:  TypeInt32,
	KeywordString: TypeString,
}
