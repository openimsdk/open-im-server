package checker

type Checker interface {
	Check() error
}

func Validate(args any) error {
	if checker, ok := args.(Checker); ok {
		if err := checker.Check(); err != nil {
			return err
		}
	}
	return nil
}
