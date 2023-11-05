package field

import (
	"reflect"
	"testing"
)

type fieldValuesForTest struct {
	f1 int64
	f2 string
}

func TestField_ValueOf(t *testing.T) {
	type args struct {
		r any
	}

	var (
		f1 Field = "f1"
		f2 Field = "f2"
	)
	f1.InitValueOf(func(a any) any { return a.(*fieldValuesForTest).f1 })
	f2.InitValueOf(func(a any) any { return a.(*fieldValuesForTest).f2 })

	fvs := &fieldValuesForTest{
		f1: 100,
		f2: "abc",
	}

	tests := []struct {
		name string
		f    Field
		args args
		want any
	}{
		{
			"f1",
			f1,
			args{r: fvs},
			fvs.f1,
		},
		{
			"f2",
			f2,
			args{r: fvs},
			fvs.f2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.f.ValueOf(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Field.ValueOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
