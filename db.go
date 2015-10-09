package cellnet

type KVDatabase interface {
	Start(rawcfg interface{}) error

	Stop()

	Insert(evq EventQueue, collName string, doc interface{}, callback func(error))

	FindOne(evq EventQueue, collName string, query interface{}, callback interface{})

	Update(evq EventQueue, collName string, selector interface{}, doc interface{}, callback func(error))
}
