package daemon

type ProcFunc func(chan int, chan int)

type Daemon struct {
	port int
	proc ProcFunc
	ch   chan int
}
