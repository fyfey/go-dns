package opcode

type Opcode int

const (
	Query  Opcode = 0
	IQuery Opcode = 1
	Status Opcode = 2
	Notify Opcode = 4
	Update Opcode = 5
)
