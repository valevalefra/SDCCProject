package utility

type MessageType int
type TimeStamp int

const (
	Request MessageType = iota + 1 // request mutual lock
	Reply                          // reply mutual lock
	Release                        // release mutual lock
)

type Message struct {
	Type   MessageType
	SendID int
	TS     TimeStamp
	Text   string
}
