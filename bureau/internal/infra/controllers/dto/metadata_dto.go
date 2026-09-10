package dto

// ResponseMetadata identifica a origem e a chamada que produziu a resposta.
// O mesmo transaction_id e registrado no log do servidor, o que permite cruzar
// a transcricao do agente com os registros de execucao de cada ferramenta.
type ResponseMetadata struct {
	Source        string `json:"source"`
	Tool          string `json:"tool"`
	QueriedAt     string `json:"queried_at"`
	TransactionID string `json:"transaction_id"`
}

// Traceable e embutido nas respostas de primeiro nivel das ferramentas MCP.
type Traceable struct {
	Metadata *ResponseMetadata `json:"metadata,omitempty"`
}

func (t *Traceable) SetMetadata(metadata ResponseMetadata) {
	t.Metadata = &metadata
}
