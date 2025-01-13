package schemas

import "time"

type User struct {
	ID        string    `json:"User.id"`
	Name      string    `json:"User.name"`
	Email     string    `json:"User.email"`
	Password  string    `json:"User.password"`
	CreatedAt time.Time `json:"User.created_at"`
	UpdatedAt time.Time `json:"User.updated_at"`
	DType     []string  `json:"dgraph.type,omitempty"`
}

type LoginUser struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type Research struct {
	User             User        `json:"Research.user"`
	ResearchType     string      `json:"Research.research_type"`
	Title            string      `json:"Research.title"`
	Description      string      `json:"Research.description"`
	PubmedIds        []string    `json:"Research.pubmed_ids"`
	AssociatedChunks []TextChunk `json:"Research.associated_chunks"`
	ResearchResult   string      `json:"Research.research_result"`
	DType            []string    `json:"dgraph.type,omitempty"`
}

type Chat struct {
	User     User          `json:"Chat.user"`
	Research Research      `json:"Chat.research"`
	Messages []ChatMessage `json:"Chat.messages"`
	DType    []string      `json:"dgraph.type,omitempty"`
}

type ChatMessage struct {
	Chat         Chat     `json:"ChatMessage.chat"`
	HumanMessage string   `json:"ChatMessage.human_message"`
	AIMessage    string   `json:"ChatMessage.ai_message"`
	DType        []string `json:"dgraph.type,omitempty"`
}
