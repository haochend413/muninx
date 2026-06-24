package app

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	editstack "github.com/haochend413/muninx/internal/app/editStack"
	"github.com/haochend413/muninx/internal/models"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/haochend413/muninx/sys"
)

// commit.go implements the note commit system: snapshotting a note's
// content as a patch (CommitNote), undoing the most recent snapshot
// (OmitCommit), and flagging one as archived (ArchiveCommit, which doesn't
// change behavior yet — it's just a marker for now).

// CommitNote snapshots the current note: it diffs Content against
// CheckedContent, records the result as a new commit, and advances
// CheckedContent to Content. No-op (returns false) if nothing has changed
// since the last commit.
func (a *App) CommitNote() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil {
		return false
	}
	if note.CheckedContent == note.Content {
		return false
	}

	dmp := diffmatchpatch.New()
	patch := dmp.PatchToText(dmp.PatchMake(note.CheckedContent, note.Content))
	note.Commits = append(note.Commits, &models.NoteCommit{
		Patch:      patch,
		CommitTime: time.Now(),
	})
	note.Diff = dmp.DiffPrettyText(dmp.DiffMain(note.CheckedContent, note.Content, false))
	note.CheckedContent = note.Content

	a.Synced = false
	edit := &editstack.Edit{ID: note.ID, EditType: editstack.UpdateNote}
	if err := a.editMgr.AddEdit(edit); err != nil {
		sys.LogError(err)
	}
	return true
}

// OmitCommit undoes the most recent commit: it applies the reverse of its
// patch to CheckedContent, then removes the commit. No-op (returns false)
// if there are no commits.
func (a *App) OmitCommit() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil || len(note.Commits) == 0 {
		return false
	}

	last := note.Commits[len(note.Commits)-1]
	restored, err := applyReversePatch(last.Patch, note.CheckedContent)
	if err != nil {
		sys.LogError(fmt.Errorf("omit commit: %w", err))
		return false
	}

	dmp := diffmatchpatch.New()
	note.CheckedContent = restored
	note.Diff = dmp.DiffPrettyText(dmp.DiffMain(note.CheckedContent, note.Content, false))
	note.Commits = note.Commits[:len(note.Commits)-1]

	a.Synced = false
	edit := &editstack.Edit{ID: note.ID, EditType: editstack.UpdateNote}
	if err := a.editMgr.AddEdit(edit); err != nil {
		sys.LogError(err)
	}
	return true
}

// ArchiveCommit marks the most recent commit as archived. Archiving doesn't
// change any behavior yet — it's just a flag for now. No-op if there are no
// commits.
func (a *App) ArchiveCommit() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil || len(note.Commits) == 0 {
		return
	}
	note.Commits[len(note.Commits)-1].Archived = true
}

// GetCommitCount returns the number of commits recorded for the active
// note. Returns 0 if there's no active note.
func (a *App) GetCommitCount() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil {
		return 0
	}
	return len(note.Commits)
}

// GetCommitDiffAt returns a red/green ANSI-colored rendering of the commit
// at the given index in the note's commit stack (0 = oldest), for display.
// It reconstructs that commit's before/after text by applying later
// commits' patches in reverse, starting from CheckedContent. Returns a
// placeholder if the index is out of range.
func (a *App) GetCommitDiffAt(index int) string {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil || index < 0 || index >= len(note.Commits) {
		return "No commit at this position."
	}

	after := note.CheckedContent
	for i := len(note.Commits) - 1; i > index; i-- {
		restored, err := applyReversePatch(note.Commits[i].Patch, after)
		if err != nil {
			sys.LogError(fmt.Errorf("render commit diff: %w", err))
			return "Could not render this commit's diff."
		}
		after = restored
	}

	before, err := applyReversePatch(note.Commits[index].Patch, after)
	if err != nil {
		sys.LogError(fmt.Errorf("render commit diff: %w", err))
		return "Could not render this commit's diff."
	}

	dmp := diffmatchpatch.New()
	return dmp.DiffPrettyText(dmp.DiffMain(before, after, false))
}

// GetNextCommitPreview returns a red/green ANSI-colored preview of what
// CommitNote would record right now — the diff between CheckedContent and
// the current Content — without creating a commit or changing any state.
func (a *App) GetNextCommitPreview() string {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	note := a.getCurrentNote()
	if note == nil {
		return "No note selected."
	}
	if note.CheckedContent == note.Content {
		return "No changes to commit."
	}

	dmp := diffmatchpatch.New()
	return dmp.DiffPrettyText(dmp.DiffMain(note.CheckedContent, note.Content, false))
}

// applyReversePatch reverses patchText (a diffmatchpatch unified-diff-style
// patch) and applies it to text, returning what text looked like before the
// patch was originally applied.
func applyReversePatch(patchText, text string) (string, error) {
	dmp := diffmatchpatch.New()
	reversed, err := dmp.PatchFromText(reversePatchText(patchText))
	if err != nil {
		return "", fmt.Errorf("failed to parse reversed patch: %w", err)
	}
	restored, _ := dmp.PatchApply(reversed, text)
	return restored, nil
}

// patchHeaderRe matches a diffmatchpatch patch hunk header, e.g. "@@ -1,4 +1,6 @@".
var patchHeaderRe = regexp.MustCompile(`^@@ -([0-9,]+) \+([0-9,]+) @@$`)

// reversePatchText flips a diffmatchpatch patch so applying it undoes what
// the original did: each hunk's "-old +new" header is swapped, and each
// line's '+'/'-' prefix is swapped (context " " lines are left alone).
func reversePatchText(patch string) string {
	if patch == "" {
		return patch
	}
	lines := strings.Split(patch, "\n")
	for i, line := range lines {
		if m := patchHeaderRe.FindStringSubmatch(line); m != nil {
			lines[i] = "@@ -" + m[2] + " +" + m[1] + " @@"
			continue
		}
		switch {
		case strings.HasPrefix(line, "+"):
			lines[i] = "-" + line[1:]
		case strings.HasPrefix(line, "-"):
			lines[i] = "+" + line[1:]
		}
	}
	return strings.Join(lines, "\n")
}
