package flow

type item func() (err error)

type Mixin struct {
	Error error
}

func (this *Mixin) Exec(items ...item) error {
	for _, item := range items {
		this.Error = item()

		if this.Error != nil {
			break
		}
	}

	return this.Error
}

func Go(items ...item) (err error) {
	return new(Mixin).Exec(items...)
}
