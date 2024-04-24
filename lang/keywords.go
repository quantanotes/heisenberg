package lang

type Keyword string

const (
	KeywordFrom   Keyword = "from"
	KeywordCreate Keyword = "create"
	KeywordTable  Keyword = "table"
	KeywordSelect Keyword = "select"
	KeywordInsert Keyword = "insert"
	KeywordUpdate Keyword = "update"
	KeywordDelete Keyword = "delete"
	KeywordFilter Keyword = "filter"
	KeywordInt32  Keyword = "int32"
	KeywordString Keyword = "string"
)

var Keywords = map[Keyword]bool{
	KeywordFrom:   true,
	KeywordCreate: true,
	KeywordTable:  true,
	KeywordSelect: true,
	KeywordUpdate: true,
	KeywordInsert: true,
	KeywordDelete: true,
	KeywordFilter: true,
	KeywordInt32:  true,
	KeywordString: true,
}
