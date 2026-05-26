# UI scheme

## General Behavior

There are 4 Views: MenuView, WriteView, QuitView, and FindNoteView (overlay).

When the user launches, MenuView appears, showing the ASCII logo header, an interactive notes table (most recent note ID first), and a text input bar (currently display-only).

## Views and Keybindings

### MenuView

**Layout:** header (ASCII logo) | notes table | text input | help bar

| Key       | Action                                            |
| --------- | ------------------------------------------------- |
| `N`       | Create new note (thread[0]/branch[0]) → WriteView |
| `Enter`   | Open selected note → WriteView                    |
| `j` / `↓` | Move cursor down                                  |
| `k` / `↑` | Move cursor up                                    |
| `Ctrl+F`  | Open FindNote overlay                             |
| `Ctrl+Q`  | Sync to database                                  |
| `Ctrl+C`  | Open quit confirmation                            |

### WriteView

**Layout:** Related Notes list (40% width) | Editor (textarea_vim, 60% width) | help bar

Two focus modes: **TextArea** (default on open) and **List**.

| Key      | Condition        | Action                                               |
| -------- | ---------------- | ---------------------------------------------------- |
| `Tab`    | any              | Toggle focus between editor and list                 |
| `Ctrl+S` | any              | Save current note                                    |
| `Ctrl+X` | any              | Save and return to MenuView                          |
| `Esc`    | list focused     | Save and return to MenuView                          |
| `Esc`    | textarea focused | Handled internally by textarea_vim (vim normal mode) |
| `Enter`  | list focused     | Switch current note to selected related note         |
| `Ctrl+F` | any              | Save current note and open FindNote overlay          |
| `Ctrl+Q` | any              | Save, sync to database, refresh menu table           |
| `Ctrl+C` | any              | Open quit confirmation                               |

### QuitView

Simple confirmation prompt.

| Key         | Action                                          |
| ----------- | ----------------------------------------------- |
| `y`         | Save current note, sync DB, persist state, quit |
| `n` / `Esc` | Cancel and return to previous view              |

### FindNoteView (overlay)

**Layout:** Threads table | Branches table | Notes table | Content viewport | help bar

Three focus columns: **Threads** (default on open) → **Branches** → **Notes**.

The dependent tables update automatically as the cursor moves:

- Moving the cursor in Threads refreshes the Branches table (for the selected thread) and clears Notes.
- Moving the cursor in Branches refreshes the Notes table (for the selected branch).
- Moving the cursor in Notes updates the viewport content.

The viewport shows different content depending on which column is focused:

- **Threads**: thread ID, name, summary, branch count, frequency
- **Branches**: branch ID, name, summary, note count, frequency
- **Notes**: note ID, last-edit time, full content

| Key       | Condition     | Action                                                |
| --------- | ------------- | ----------------------------------------------------- |
| `Esc`     | any           | Close overlay, return to previous view                |
| `h` / `←` | any           | Move focus one column left                            |
| `l` / `→` | any           | Move focus one column right                           |
| `j` / `↓` | any           | Navigate cursor down in focused column                |
| `k` / `↑` | any           | Navigate cursor up in focused column                  |
| `Enter`   | Notes focused | Save current note and open selected note in WriteView |

> Note: a Commits column is planned once NoteCommit persistence is implemented.

## Open / Close FindNoteView

- **Open**: `Ctrl+F` from MenuView or WriteView (saves current note first if in WriteView).
- **Close**: `Esc` from FindNoteView — returns to whichever view was active before opening.

## View Arrangement Schemes

```
MenuView:
  ┌─────────────────────────────┐
  │  ASCII logo (openmenu)      │
  ├─────────────────────────────┤
  │  Notes table (ID ↓ desc)    │
  ├─────────────────────────────┤
  │  Text input                 │
  ├─────────────────────────────┤
  │  Help bar                   │
  └─────────────────────────────┘

WriteView:
  ┌──────────────────┬──────────┐
  │  Editor          │ Related  │
  │  (textarea_vim)  │ Notes    │
  │  60% width       │ 40%      │
  └──────────────────┴──────────┘
  │  Help bar                   │
  └─────────────────────────────┘

FindNoteView:
  ┌────────┬─────────┬───────┬─────────────────┐
  │Threads │Branches │ Notes │   Content        │
  │        │         │       │   viewport       │
  └────────┴─────────┴───────┴─────────────────┘
  │  Help bar                                   │
  └─────────────────────────────────────────────┘
```
