package rocrail

import "fmt"

// Direction marks a forward or reverse locomotive direction.
type Direction bool

const (
	// Forward is the standard direction.
	Forward Direction = true
	// Reverse is for locomotives driving back-first or pushing a train in front of it.
	Reverse Direction = false
)

// Logo represents a locomotive / drive unit.
type Loco struct {
	id   string
	conn *Connection

	velocity int
	dir      Direction
}

// Loco returns a locomotive instance for the given ID (i.e. "0003").
func (c *Connection) Loco(id string) *Loco {
	// TODO read current state
	return &Loco{id: id, conn: c, dir: Forward}
}

// SetDirection allows you to set a new direction for this locomotive.
// It will attempt to regain the same speed in the opposite direction.
func (l *Loco) SetDirection(d Direction) error {
	l.dir = d
	return l.sendVelocity()
}

// SetVelocity specifies the desired speed for this locomotive - between 0 and 100.
func (l *Loco) SetVelocity(val int) error {
	l.velocity = val
	return l.sendVelocity()
}

func (l *Loco) sendVelocity() error {
	xml := fmt.Sprintf("<lc id=\"%s\" V=\"%d\" dir=\"%t\" cmd=\"velocity\" />", l.id, l.velocity, l.dir)
	return l.conn.SendXML("lc", xml)
}
