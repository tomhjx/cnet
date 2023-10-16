package field

type Field string

var fieldValHandles = map[string]FieldValueHandle{}

func (f Field) ValueOf(r any) any {
	ff := string(f)
	if h, ok := fieldValHandles[ff]; ok {
		return h(r)
	}
	return nil
}

func (f Field) Inject(h FieldValueHandle) {
	fieldValHandles[string(f)] = h
}

func (f Field) String() string {
	return string(f)
}

type FieldValueHandle func(any) any
