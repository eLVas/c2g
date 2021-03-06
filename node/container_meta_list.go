package node

import (
	"github.com/c2stack/c2g/meta"
)

type ContainerMetaList struct {
	next       meta.Meta
	main       meta.MetaIterator
	choiceCase *ContainerMetaList
	s          Selection
}

func NewContainerMetaList(s Selection) *ContainerMetaList {
	i := &ContainerMetaList{
		main: meta.NewMetaListIterator(s.Path.meta, true),
		s:    s,
	}
	i.lookAhead()
	return i
}

func newChoiceCaseIterator(s Selection, m *meta.ChoiceCase) *ContainerMetaList {
	i := &ContainerMetaList{
		main: meta.NewMetaListIterator(m, true),
		s:    s,
	}
	i.lookAhead()
	return i
}

func (self *ContainerMetaList) Next() meta.Meta {
	var next = self.next
	self.lookAhead()
	return next
}

func (self *ContainerMetaList) lookAhead() error {
	self.next = nil
	var m meta.Meta
	for  {
		if self.choiceCase != nil {
			m = self.choiceCase.Next()
			if m == nil {
				self.choiceCase = nil
				continue
			}
		} else if self.main != nil {
			m = self.main.NextMeta()
			if m == nil {
				break
			}
		}
		if choice, isChoice := m.(*meta.Choice); isChoice {
			if chosen, err := self.s.Node.Choose(self.s, choice); err != nil {
				return err
			} else if chosen != nil {
				self.choiceCase = newChoiceCaseIterator(self.s, chosen)
				continue
			}
		} else {
			self.next = m
			break
		}
	}
	return nil
}
