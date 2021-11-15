package data

import "fmt"

func errs2err(errs []error) error {
	var err error
	for _, v := range errs {
		err = fmt.Errorf("%v;%v", err, v)
	}
	return err
}
