package log // import "go.microcore.dev/framework/log"

type OutputFormat string

const (
	// FormatText writes plain text logs, human-readable, suitable for local or file output.
	FormatText OutputFormat = "text"

	// FormatJSON writes structured JSON logs, recommended for production and log aggregation.
	FormatJSON OutputFormat = "json"

	// FormatPretty writes colorized, developer-friendly logs, automatically disables
	// color if output is not a terminal.
	FormatPretty OutputFormat = "pretty"
)
