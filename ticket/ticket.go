package ticket

import lib "github.com/google/uuid"

type Ticket struct {
	uuid lib.UUID
}

func New() *Ticket {
	return &Ticket{
		uuid: lib.New(),
	}
}

func FromRaw(b []byte) (*Ticket, error) {
	uuid, err := lib.FromBytes(b)
	if err != nil {
		return nil, err
	}

	return &Ticket{
		uuid: uuid,
	}, nil
}

func Parse(s string) (*Ticket, error) {
	uuid, err := lib.Parse(s)
	if err != nil {
		return nil, err
	}

	return &Ticket{
		uuid: uuid,
	}, nil
}

func (t *Ticket) String() string {
	return t.uuid.String()
}

func (t *Ticket) Raw() []byte {
	return t.uuid[:]
}
