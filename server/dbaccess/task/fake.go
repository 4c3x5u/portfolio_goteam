package task

type FakeInserter struct{}

func (f *FakeInserter) Insert(_ Task) error { return nil }
