package node

import (
	"testing"
	"github.com/c2stack/c2g/meta/yang"
	"bytes"
	"github.com/c2stack/c2g/meta"
)

func TestDump(t *testing.T) {
	mstr := `
module food {
	prefix "x";
	namespace "y";
	revision 0000-00-00 {
		description "";
	}
	list fruits  {
		key "name";
		leaf name {
			type string;
		}
		container origin {
			leaf country {
				type string;
			}
		}
	}
}
`
	m, err := yang.LoadModuleCustomImport(mstr, nil)
	if err != nil {
		t.Fatal(err)
	}
	var actual bytes.Buffer
	var dump bytes.Buffer
	out := Dump(NewJsonWriter(&actual).Node(), &dump)
	if err = SelectModule(meta.PathStreamSource("../yang"), m, true).Root().InsertInto(out).LastErr; err != nil {
		t.Fatal(err)
	}
	t.Log(dump.String())
}
