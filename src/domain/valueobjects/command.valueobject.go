package valueobjects

type CommandValue string

const (
	MCP            CommandValue = "https.local:8080"
	RAG            CommandValue = "https.local:8081"
	ResponseModule CommandValue = "https.local:8082"
)

var ValidCommands = []CommandValue{
	MCP,
	RAG,
	ResponseModule,
}


