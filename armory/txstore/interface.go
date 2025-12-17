package txstore

type TxStore interface {
	Init() error

	InsertTx(tx *TxEvent) error
	InsertBatch(txs []*TxEvent) error

	// 顺序回放（状态机）
	ScanFrom(
		fromBlock int64,
		fn func(*TxEvent) error,
	) error

	// 时间 + 类型筛选（报表 / 分析）
	QueryByTimeAndType(
		start, end int64,
		txType string,
		fn func(*TxEvent) error,
	) error

	Close() error
}
