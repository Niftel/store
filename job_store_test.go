package store

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadableByScopesQueryFailsClosedWithoutScopes(t *testing.T) {
	query, args, ok, err := readableByScopesQuery(nil, nil, 50)
	if err != nil {
		t.Fatal(err)
	}
	if ok || query != "" || args != nil {
		t.Fatalf("empty scopes produced query=%q args=%v ok=%v", query, args, ok)
	}
}

func TestReadableByScopesQueryUsesAuthorizedTemplateAndInventoryIDs(t *testing.T) {
	tests := []struct {
		name          string
		templateIDs   []int64
		inventoryIDs  []int64
		wantFragments []string
		absent        []string
		wantArgs      []interface{}
	}{
		{
			name: "template only", templateIDs: []int64{11, 12},
			wantFragments: []string{"FROM job_templates jt", "jt.id IN (?, ?)"},
			absent:        []string{"inventory_sync_history"},
			wantArgs:      []interface{}{int64(11), int64(12), 50},
		},
		{
			name: "inventory only", inventoryIDs: []int64{21},
			wantFragments: []string{"FROM inventory_sync_history ish", "ish.inventory_id IN (?)"},
			absent:        []string{"job_templates"},
			wantArgs:      []interface{}{int64(21), 50},
		},
		{
			name: "both", templateIDs: []int64{31}, inventoryIDs: []int64{41, 42},
			wantFragments: []string{"FROM job_templates jt", "FROM inventory_sync_history ish", " OR ", "ORDER BY uj.created_at DESC, uj.id DESC"},
			wantArgs:      []interface{}{int64(31), int64(41), int64(42), 50},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, args, ok, err := readableByScopesQuery(tt.templateIDs, tt.inventoryIDs, 50)
			if err != nil {
				t.Fatal(err)
			}
			if !ok {
				t.Fatal("authorized scopes did not produce a query")
			}
			for _, fragment := range tt.wantFragments {
				if !strings.Contains(query, fragment) {
					t.Fatalf("query missing %q: %s", fragment, query)
				}
			}
			for _, fragment := range tt.absent {
				if strings.Contains(query, fragment) {
					t.Fatalf("query unexpectedly contains %q: %s", fragment, query)
				}
			}
			if !reflect.DeepEqual(args, tt.wantArgs) {
				t.Fatalf("args=%#v want %#v", args, tt.wantArgs)
			}
		})
	}
}
