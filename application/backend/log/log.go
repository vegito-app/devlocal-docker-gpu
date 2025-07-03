package log

import (
	"github.com/rs/zerolog"
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

}

// Output: {"time":1516134303,"level":"debug","message":"hello world"}
