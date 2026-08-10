package ldap

import (
	"testing"
)

// ── parseResourceLDAPGroupAttributeID ───────────────────────────────────────

func TestParseResourceLDAPGroupAttributeID_Valid(t *testing.T) {
	tests := []struct {
		id        string
		wantGroup string
		wantAttr  string
	}{
		{
			id:        "CN=MyGroup,OU=Groups,DC=example,DC=com/department",
			wantGroup: "CN=MyGroup,OU=Groups,DC=example,DC=com",
			wantAttr:  "department",
		},
		{
			id:        "CN=MyGroup,OU=Groups,DC=example,DC=com/telephoneNumber",
			wantGroup: "CN=MyGroup,OU=Groups,DC=example,DC=com",
			wantAttr:  "telephoneNumber",
		},
		{
			id:        "CN=My Group,OU=Groups,DC=example,DC=com/extensionAttribute1",
			wantGroup: "CN=My Group,OU=Groups,DC=example,DC=com",
			wantAttr:  "extensionAttribute1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			gotGroup, gotAttr, err := parseResourceLDAPGroupAttributeID(tt.id)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotGroup != tt.wantGroup {
				t.Errorf("group_dn: got %q, want %q", gotGroup, tt.wantGroup)
			}
			if gotAttr != tt.wantAttr {
				t.Errorf("attribute_name: got %q, want %q", gotAttr, tt.wantAttr)
			}
		})
	}
}

func TestParseResourceLDAPGroupAttributeID_Invalid(t *testing.T) {
	invalid := []string{
		"",
		"no-slash-at-all",
		"/missingGroupDN",
		"CN=Group,OU=x,DC=y/",
	}

	for _, id := range invalid {
		t.Run(id, func(t *testing.T) {
			_, _, err := parseResourceLDAPGroupAttributeID(id)
			if err == nil {
				t.Errorf("expected error for id %q, got nil", id)
			}
		})
	}
}

// ── resourceLDAPGroupAttributeID ────────────────────────────────────────────

func TestResourceLDAPGroupAttributeID_RoundTrip(t *testing.T) {
	groupDN := "CN=MyGroup,OU=Groups,DC=example,DC=com"
	attrName := "department"

	id := resourceLDAPGroupAttributeID(groupDN, attrName)

	gotGroup, gotAttr, err := parseResourceLDAPGroupAttributeID(id)
	if err != nil {
		t.Fatalf("round-trip parse failed: %v", err)
	}
	if gotGroup != groupDN {
		t.Errorf("group_dn: got %q, want %q", gotGroup, groupDN)
	}
	if gotAttr != attrName {
		t.Errorf("attribute_name: got %q, want %q", gotAttr, attrName)
	}
}

// ── expandStringList ─────────────────────────────────────────────────────────

func TestExpandStringList(t *testing.T) {
	tests := []struct {
		name  string
		input []interface{}
		want  []string
	}{
		{
			name:  "single value",
			input: []interface{}{"value1"},
			want:  []string{"value1"},
		},
		{
			name:  "multiple values",
			input: []interface{}{"aaa", "bbb", "ccc"},
			want:  []string{"aaa", "bbb", "ccc"},
		},
		{
			name:  "empty list",
			input: []interface{}{},
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandStringList(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("length mismatch: got %d, want %d", len(got), len(tt.want))
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("index %d: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// ── Schema validation ────────────────────────────────────────────────────────

func TestResourceLDAPGroupAttributeSchema(t *testing.T) {
	res := resourceLDAPGroupAttribute()

	requiredForceNew := []string{"group_dn", "attribute_name"}
	for _, attr := range requiredForceNew {
		s, ok := res.Schema[attr]
		if !ok {
			t.Errorf("schema missing field %q", attr)
			continue
		}
		if !s.Required {
			t.Errorf("field %q should be Required", attr)
		}
		if !s.ForceNew {
			t.Errorf("field %q should be ForceNew", attr)
		}
	}

	valuesSchema, ok := res.Schema["attribute_values"]
	if !ok {
		t.Fatal("schema missing field \"attribute_values\"")
	}
	if !valuesSchema.Required {
		t.Error("field \"attribute_values\" should be Required")
	}
	if valuesSchema.ForceNew {
		t.Error("field \"attribute_values\" should NOT be ForceNew (must be updatable in-place)")
	}
	if valuesSchema.MinItems != 1 {
		t.Errorf("field \"attribute_values\" MinItems: got %d, want 1", valuesSchema.MinItems)
	}
}
