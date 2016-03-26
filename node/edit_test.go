package node

import (
	"log"
	"github.com/blitter/meta/yang"
	"strings"
	"testing"
	"github.com/blitter/meta"
)

const EDIT_TEST_MODULE = `
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

func TestEditListItem(t *testing.T) {
	m := YangFromString(EDIT_TEST_MODULE)
	root := testDataRoot()
	bd := MapNode(root)
	json := NewJsonReader(strings.NewReader(`{"origin":{"country":"Canada"}}`)).Node()

	// UPDATE
	// Here we're testing editing a specific list item. With FindTarget walk controller
	// needs to leave walkstate in a position for WalkTarget controller to make the edit
	// on the right item.
	log.Println("Testing edit\n")
	c := NewContext()
	if err := c.Select(m, bd).Find("fruits=apple").UpdateFrom(json).LastErr; err != nil {
		t.Fatal(err)
	}
	actual := MapValue(root, "fruits.1.origin.country")
	if actual != "Canada" {
		t.Error("Edit failed", actual)
	}

	// INSERT
	log.Println("Testing insert\n")
	json = NewJsonReader(strings.NewReader(`{"fruits":[{"name":"pear","origin":{"country":"Columbia"}}]}`)).Node()
	if err := c.Select(m, bd).Find("fruits").InsertFrom(json).LastErr; err != nil {
		t.Fatal(err)
	}
	actual, found := root["fruits"]
	if !found {
		t.Error("fruits not found")
	} else {
		fruits := actual.([]map[string]interface{})
		if len(fruits) != 3 {
			t.Error("Expected 3 fruits but got ", len(fruits))
		}
	}
}

func testDataRoot() map[string]interface{} {
	return map[string]interface{}{
		"fruits" : []map[string]interface{}{
			map[string]interface{}{
				"name" : "banana",
				"origin" : map[string]interface{}{
					"country" : "Brazil",
				},
			},
			map[string]interface{}{
				"name" : "apple",
				"origin" : map[string]interface{}{
					"country": "US",
				},
			},
		},
	}
}


func YangFromString(s string) (*meta.Module) {
	m, err := yang.LoadModuleCustomImport(s, nil)
	if err != nil {
		panic(err)
	}
	return m
}
