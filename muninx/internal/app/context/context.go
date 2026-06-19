package context

// Context is a wrapper that helps ordering and filtering the raw lists of notes, and store lists from the map.
// In the future, maybe we can add a config menu that allows for multiple orders.
// Context also defines current list, which in combination with cursor, defines current item pointers.

type ContextPtr int

const (
	None    ContextPtr = -1
	Default ContextPtr = 0
	Recent  ContextPtr = 1
	Search  ContextPtr = 2
)

// This is replicative, but might be useful in the future.
type ContextOrder int

const (
	CreateAt ContextOrder = 0 // default , time order
	UpdateAt ContextOrder = 1 // recent, most recently updated
)
