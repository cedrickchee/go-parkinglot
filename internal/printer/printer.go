package printer

import (
	"fmt"
	"io"
	"reflect"
)

// Fprintf is a convenience wrapper for fmt.Fprintf
// that pretty prints any array, slice, or string.
func Fprintf(w io.Writer, in interface{}) error {
	v := reflect.ValueOf(in)
	if (v.Kind() != reflect.Slice) &&
		(v.Kind() != reflect.Array) &&
		(v.Kind() != reflect.String) {
		return fmt.Errorf("Incompatible input type: %v", v.Kind())
	}
	if v.Len() == 0 {
		return nil
	}
	for i := 0; i < v.Len()-1; i++ {
		fmt.Fprintf(w, "%v, ", v.Index(i))
	}
	fmt.Fprintf(w, "%v\n", v.Index(v.Len()-1))
	return nil
}
