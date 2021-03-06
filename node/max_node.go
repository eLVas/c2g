package node

import "github.com/c2stack/c2g/c2"

type MaxNode struct {
	Count int
	Max   int
}

func (self MaxNode) CheckContainerPreConstraints(r *ChildRequest) (bool, error) {
	if r.IsNavigation() {
		return true, nil
	}
	self.Count++
	if self.Count > self.Max {
		r.ConstraintsHandler.IncompleteResponse(r.Selection.Path)
		// FATAL
		return false, c2.NewErrC("Too many nodes", 413)
	}
	return true, nil
}
