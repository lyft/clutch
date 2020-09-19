package experimentstore

type mockCounter struct {
	count int64
}

func (counter *mockCounter) Inc(count int64) {
	counter.count += count
}
