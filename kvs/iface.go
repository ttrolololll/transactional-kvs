package kvs

// IKvs is the interface for a key-value store
type IKvs interface {
	Set(k, v string)
	Get(k string) string
	Delete(k string)
	Count(v string) uint64
	//Begin()
	//Commit()
	//Rollback()
}
