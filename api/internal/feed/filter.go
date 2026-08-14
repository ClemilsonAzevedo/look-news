package feed

// todo: criar client para fazer chamadas ao groq ai e fazer fultro

// type req struct {
// 	Query    string    `json:"query"`
// 	Articles []Article `json:"articles"`
// }

// type response struct {
// 	Relevant []string `json:"relevant"`
// }

type Filter struct {
	parser *Parser
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) ApplyFilter() ([]Article, error) {
	return []Article{}, nil
}
