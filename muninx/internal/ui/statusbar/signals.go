package statusbar

import "time"

// QuitPrompt returns a ready-made Signal for confirming a quit action, so
// every view asks the same way.
func QuitPrompt() Signal {
	return Signal{
		Content: "Quit muninx?  Ctrl+C = confirm   n / esc = cancel",
		Fg:      "51", // bright blue
		Bold:    true,
	}
}

// flashDuration is how long a transient confirmation message stays visible
// before it auto-clears.
const flashDuration = 750 * time.Millisecond

// Synced returns a brief confirmation that pending changes were pushed to
// the database.
func Synced() Signal {
	return Signal{
		Content:  "Changes pushed to database!",
		Fg:       "51", // bright green
		Bold:     true,
		Duration: flashDuration,
	}
}

// Saved returns a brief confirmation that the note's modified content was
// saved.
func Saved() Signal {
	return Signal{
		Content:  "Modified content saved.",
		Fg:       "51", // bright green
		Bold:     true,
		Duration: flashDuration,
	}
}

// Committed returns a brief confirmation that a commit was pushed.
func Committed() Signal {
	return Signal{
		Content:  "Commit pushed.",
		Fg:       "51", // bright green
		Bold:     true,
		Duration: flashDuration,
	}
}

// CommitOmitted returns a brief confirmation that the last commit was undone.
func CommitOmitted() Signal {
	return Signal{
		Content:  "Last commit omitted.",
		Fg:       "51", // bright green
		Bold:     true,
		Duration: flashDuration,
	}
}
