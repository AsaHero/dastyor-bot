package inerr

// error not found
type ErrNotFound struct {
	name string
}

func (e *ErrNotFound) Error() string {
	return e.name + " not found"
}

func IsErrNotFound(err error) bool {
	_, ok := err.(*ErrNotFound)
	return ok
}

func NewErrNotFound(text string) *ErrNotFound {
	return &ErrNotFound{text}
}

// error conflict
type ErrConflict struct {
	name string
}

func (e *ErrConflict) Error() string {
	return e.name + " already exist"
}

func IsErrConflict(err error) bool {
	_, ok := err.(*ErrConflict)
	return ok
}

func NewErrConflict(text string) *ErrConflict {
	return &ErrConflict{text}
}

// error no changes
type ErrNoChanges struct {
	name string
}

func (e *ErrNoChanges) Error() string {
	return e.name + " is not changed"
}

func IsErrNoChanges(err error) bool {
	_, ok := err.(*ErrNoChanges)
	return ok
}

func NewErrNoChanges(text string) *ErrNoChanges {
	return &ErrNoChanges{text}
}
