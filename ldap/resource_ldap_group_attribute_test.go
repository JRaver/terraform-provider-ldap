package ldap

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
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

// ── orderLike ────────────────────────────────────────────────────────────────

func TestOrderLike(t *testing.T) {
	tests := []struct {
		name   string
		actual []string
		prior  []string
		want   []string
	}{
		{
			name:   "same set, different order returns prior",
			actual: []string{"ccc", "aaa", "bbb"},
			prior:  []string{"aaa", "bbb", "ccc"},
			want:   []string{"aaa", "bbb", "ccc"},
		},
		{
			name:   "case-only difference returns prior",
			actual: []string{"mixeduser", "CN=Bob,DC=example,DC=com"},
			prior:  []string{"cn=bob,dc=example,dc=com", "MixedUser"},
			want:   []string{"cn=bob,dc=example,dc=com", "MixedUser"},
		},
		{
			name:   "mixed-case value in AD order returns prior",
			actual: []string{"user03@example.com", "user02@example.com", "MixedUser@example.com", "user01@example.com"},
			prior:  []string{"MixedUser@example.com", "user01@example.com", "user02@example.com", "user03@example.com"},
			want:   []string{"MixedUser@example.com", "user01@example.com", "user02@example.com", "user03@example.com"},
		},
		{
			name:   "AD returns different casing returns prior with config casing",
			actual: []string{"user01@example.com", "mixeduser@EXAMPLE.COM"},
			prior:  []string{"MixedUser@example.com", "user01@example.com"},
			want:   []string{"MixedUser@example.com", "user01@example.com"},
		},
		{
			name:   "all upper case in AD returns prior",
			actual: []string{"USER01@EXAMPLE.COM", "MIXEDUSER@EXAMPLE.COM"},
			prior:  []string{"MixedUser@example.com", "user01@example.com"},
			want:   []string{"MixedUser@example.com", "user01@example.com"},
		},
		{
			name:   "case-insensitive duplicates are matched one-to-one",
			actual: []string{"mixeduser@example.com", "MIXEDUSER@example.com"},
			prior:  []string{"MixedUser@example.com", "MixedUser@example.com"},
			want:   []string{"MixedUser@example.com", "MixedUser@example.com"},
		},
		{
			name:   "mixed-case value replaced returns actual",
			actual: []string{"user01@example.com", "MixedUsex@example.com"},
			prior:  []string{"MixedUser@example.com", "user01@example.com"},
			want:   []string{"user01@example.com", "MixedUsex@example.com"},
		},
		{
			name:   "mixed-case value removed returns actual",
			actual: []string{"user01@example.com", "user02@example.com"},
			prior:  []string{"MixedUser@example.com", "user01@example.com"},
			want:   []string{"user01@example.com", "user02@example.com"},
		},
		{
			name:   "value replaced returns actual",
			actual: []string{"aaa", "ddd", "ccc"},
			prior:  []string{"aaa", "bbb", "ccc"},
			want:   []string{"aaa", "ddd", "ccc"},
		},
		{
			name:   "value added returns actual",
			actual: []string{"aaa", "bbb", "ccc"},
			prior:  []string{"bbb", "aaa"},
			want:   []string{"aaa", "bbb", "ccc"},
		},
		{
			name:   "value removed returns actual",
			actual: []string{"bbb"},
			prior:  []string{"aaa", "bbb"},
			want:   []string{"bbb"},
		},
		{
			name:   "duplicates are matched one-to-one",
			actual: []string{"aaa", "bbb"},
			prior:  []string{"aaa", "aaa"},
			want:   []string{"aaa", "bbb"},
		},
		{
			name:   "empty prior (import) returns actual",
			actual: []string{"bbb", "aaa"},
			prior:  []string{},
			want:   []string{"bbb", "aaa"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := orderLike(tt.actual, tt.prior)
			if len(got) != len(tt.want) {
				t.Fatalf("length mismatch: got %v, want %v", got, tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("index %d: got %q, want %q (full: %v)", i, got[i], tt.want[i], got)
				}
			}
		})
	}
}

const largeOwnersCount = 59

// largeOwnersFixture returns a sorted config-style list with one mixed-case value first,
// and the same set shuffled deterministically to mimic AD's internal order.
func largeOwnersFixture() (configOrder, adOrder []string) {
	configOrder = []string{"MixedUser@example.com"}
	for i := 1; i < largeOwnersCount; i++ {
		configOrder = append(configOrder, fmt.Sprintf("user%02d@example.com", i))
	}
	// 7 is coprime with 59, so i -> (7i+3) mod 59 is a permutation.
	adOrder = make([]string, largeOwnersCount)
	for i := range adOrder {
		adOrder[i] = configOrder[(7*i+3)%largeOwnersCount]
	}
	return configOrder, adOrder
}

func TestOrderLike_LargeOwnersList(t *testing.T) {
	configOrder, adOrder := largeOwnersFixture()

	mixedIdx := -1
	for i, v := range adOrder {
		if v == "MixedUser@example.com" {
			mixedIdx = i
		}
	}
	if mixedIdx <= 0 {
		t.Fatalf("fixture: mixed-case value should be shuffled away from index 0, got %d", mixedIdx)
	}

	assertEqual := func(t *testing.T, got, want []string) {
		t.Helper()
		if len(got) != len(want) {
			t.Fatalf("length mismatch: got %d, want %d", len(got), len(want))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("index %d: got %q, want %q", i, got[i], want[i])
			}
		}
	}

	t.Run("same set in AD order returns config order", func(t *testing.T) {
		assertEqual(t, orderLike(adOrder, configOrder), configOrder)
	})

	t.Run("AD lower-cases mixed-case value returns config order", func(t *testing.T) {
		lowered := make([]string, len(adOrder))
		for i, v := range adOrder {
			lowered[i] = strings.ToLower(v)
		}
		assertEqual(t, orderLike(lowered, configOrder), configOrder)
	})

	t.Run("AD upper-cases everything returns config order", func(t *testing.T) {
		upper := make([]string, len(adOrder))
		for i, v := range adOrder {
			upper[i] = strings.ToUpper(v)
		}
		assertEqual(t, orderLike(upper, configOrder), configOrder)
	})

	t.Run("one owner replaced in AD returns AD order", func(t *testing.T) {
		changed := append([]string(nil), adOrder...)
		changed[mixedIdx] = "MixedUsex@example.com"
		assertEqual(t, orderLike(changed, configOrder), changed)
	})

	t.Run("one owner removed in AD returns AD order", func(t *testing.T) {
		removed := append(append([]string(nil), adOrder[:mixedIdx]...), adOrder[mixedIdx+1:]...)
		assertEqual(t, orderLike(removed, configOrder), removed)
	})
}

func TestOrderLike_RandomShuffles(t *testing.T) {
	const iterations = 500
	rng := rand.New(rand.NewPCG(42, 1024))
	configOrder, _ := largeOwnersFixture()

	randomCase := func(s string) string {
		switch rng.IntN(3) {
		case 0:
			return strings.ToLower(s)
		case 1:
			return strings.ToUpper(s)
		default:
			return s
		}
	}

	for it := 0; it < iterations; it++ {
		actual := make([]string, len(configOrder))
		for i, v := range configOrder {
			actual[i] = randomCase(v)
		}
		rng.Shuffle(len(actual), func(i, j int) { actual[i], actual[j] = actual[j], actual[i] })

		if got := orderLike(actual, configOrder); !slices.Equal(got, configOrder) {
			t.Fatalf("iteration %d: same set should return prior\nactual: %v\ngot:    %v", it, actual, got)
		}

		mutated := slices.Clone(actual)
		mutated[rng.IntN(len(mutated))] = fmt.Sprintf("Intruder%03d@example.com", it)
		if got := orderLike(mutated, configOrder); !slices.Equal(got, mutated) {
			t.Fatalf("iteration %d: replaced value should return actual\nactual: %v\ngot:    %v", it, mutated, got)
		}

		idx := rng.IntN(len(actual))
		shrunk := slices.Delete(slices.Clone(actual), idx, idx+1)
		if got := orderLike(shrunk, configOrder); !slices.Equal(got, shrunk) {
			t.Fatalf("iteration %d: removed value should return actual\nactual: %v\ngot:    %v", it, shrunk, got)
		}
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
