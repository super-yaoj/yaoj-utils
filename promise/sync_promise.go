package promise

//-------Async Promise
type SyncPromise[Resolve any, Reject any] struct {
	res      Resolve           // last resolve
	rej      Reject            // last rejected
	IsReject func(Reject) bool // judge whether T is rejected data
}

func NewSyncPromise[Res any, Rej any](isrej func(Rej) bool, callback func() (Res, Rej)) *SyncPromise[Res, Rej] {
	res, rej := callback()
	return &SyncPromise[Res, Rej]{res, rej, isrej}
}
//give a callback function that accepts data and returns next data
//data should be able to judge whether it's rejected
func (p *SyncPromise[Res, Rej]) Then(callback func(Res) (Res, Rej)) *SyncPromise[Res, Rej] {
	if p.IsReject(p.rej) {
		return p
	}
	p.res, p.rej = callback(p.res)
	return p
}

func (p *SyncPromise[Res, Rej]) Catch(callback func(Rej)) {
	if p.IsReject(p.rej) {
		callback(p.rej)
	}
}