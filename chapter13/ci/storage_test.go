package main

type testStorage struct {
	items map[string]string
}

func (t *testStorage) Read(name string) string {
	return t.items[name]
}

func (t *testStorage) Write(name, content string) {
	t.items[name] = content
}

func newTestStorage() Storage {
	store := &testStorage{}
	store.items = make(map[string]string)
	return store
}
