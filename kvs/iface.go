package kvs

// IKvs is the interface for a key-value store
type IKvs interface {
	CommandExecutor(inputs []string) (string, error)
	Set(k, v string)
	Get(k string) (string, error)
	Delete(k string)
	Count(v string) string
	Begin()
	Commit() error
	Rollback() error
}
