package events

import "time"

// This is the events module, the most important mechanism that keeps track of what we have wrote in the textarea, segmenting it into useful pieces.
// The content should be kept synchronized with the main content.

// i think we only need to track some sorts of events. The simplest is just two.With more modular design, there might be more types.
// i dont even know how to handle delete. maybe there is no such need and can be handled elsewhere,where the edis are updated somehow. yeah.

type EventType = int

const (
	Write EventType = iota
	Delete
)

type Event struct {
	ID        int
	CreatedAt time.Time // when the event started
	UpdatedAt time.Time // when teh event ended
	Type      EventType
	Content   []rune // here we use a list of runes
	Pos       int    // right now just the line number, we will use it to track and to search. maybe better scheme.
}

// data struct used for single event storing
type EventMgr struct {
	Eventlist    []*Event
	Eventmap     map[int]*Event
	PendingRunes []rune
}

// there is no point removing events, I think. If memory is concerned, we should remove that on a higher level.

func (em *EventMgr) Init() {
	// build list and map
	em.Eventlist = []*Event{}
	em.Eventmap = make(map[int]*Event)
}

func (em *EventMgr) AppendEvent(e *Event) {
	if e != nil {
		// add event to list
		em.Eventlist = append(em.Eventlist, e)
		// add to map
		em.Eventmap[e.ID] = e
	}
}

func (em *EventMgr) GetEvent(id int) *Event {
	if value, ok := em.Eventmap[id]; ok {
		return value
	}
	return nil
}
