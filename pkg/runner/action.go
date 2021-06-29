package runner

type MessageType int

const (
	WarningMessage MessageType = iota + 1
	ErrorMessage
	NoticeMessage
)

type Message struct {
	Type MessageType
	Body string
}

type ActionResult struct {
	Err       error
	Skipped   bool
	ActionRef string
	Messages  []Message
}

func (result *ActionResult) Ok() bool {
	return result.Err == nil
}

type ActionResultIter struct {
	len     int
	channel chan *ActionResult
}

func (iter *ActionResultIter) Next() (*ActionResult, bool) {
	iter.len--

	return <-iter.channel, iter.len <= 0
}

type Action interface {
	Run(ctx *Context) ([]Message, error)
}

type ActionFunc func(*Context) ([]Message, error)

func (fn ActionFunc) Run(ctx *Context) ([]Message, error) {
	return fn(ctx)
}
