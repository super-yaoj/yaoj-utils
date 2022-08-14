package promise

//error promise
type ErrorPromise struct {
	*SyncPromise[struct{}, error]
}

func NewErrorPromise(callback func() error) *ErrorPromise {
	return &ErrorPromise{
		NewSyncPromise(
			func(err error) bool { return err != nil },
			func() (struct{}, error) {
				return struct{}{}, callback()
			},
		),
	}
}

func (p *ErrorPromise) Then(callback func() error) *ErrorPromise {
	p.SyncPromise.Then(func(struct{}) (struct{}, error) {
		return struct{}{}, callback()
	})
	return p
}

func (p *ErrorPromise) Catch(callback func(error)) {
	p.SyncPromise.Catch(callback)
}