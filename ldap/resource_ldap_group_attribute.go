package ldap

import (
	"context"
	"fmt"
	"strings"

	"github.com/Ouest-France/goldap"
	ldapv3 "github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLDAPGroupAttribute() *schema.Resource {
	return &schema.Resource{
		Description: "`ldap_group_attribute` manages a single LDAP attribute on a group. " +
			"It only controls its own attribute and does not affect any other attributes — " +
			"including those managed by the `ldap_group` resource or external systems.",
		CreateContext: resourceLDAPGroupAttributeCreate,
		ReadContext:   resourceLDAPGroupAttributeRead,
		UpdateContext: resourceLDAPGroupAttributeUpdate,
		DeleteContext: resourceLDAPGroupAttributeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLDAPGroupAttributeImport,
		},

		Schema: map[string]*schema.Schema{
			"group_dn": {
				Description: "The full DN of the LDAP group to manage the attribute on. Changes force recreation.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"attribute_name": {
				Description: "The LDAP attribute name to manage (e.g. `telephoneNumber`, `department`, `extensionAttribute1`). Changes force recreation.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"attribute_values": {
				Description: "List of values for the attribute. At least one value is required.",
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// resourceID builds the composite resource ID: <group_dn>/<attribute_name>
func resourceLDAPGroupAttributeID(groupDN, attrName string) string {
	return fmt.Sprintf("%s/%s", groupDN, attrName)
}

// parseResourceLDAPGroupAttributeID splits the composite ID back into groupDN and attrName.
func parseResourceLDAPGroupAttributeID(id string) (groupDN, attrName string, err error) {
	// SplitN with n=2 splits on the first "/" only, so the DN (which has no "/") is fully preserved.
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid ldap_group_attribute id %q: expected <group_dn>/<attribute_name>", id)
	}
	return parts[0], parts[1], nil
}

func resourceLDAPGroupAttributeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)

	groupDN := d.Get("group_dn").(string)
	attrName := d.Get("attribute_name").(string)
	values := expandStringList(d.Get("attribute_values").([]interface{}))

	req := ldapv3.NewModifyRequest(groupDN, nil)
	req.Replace(attrName, values)

	if err := client.Conn.Modify(req); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set attribute %q on group %q: %w", attrName, groupDN, err))
	}

	d.SetId(resourceLDAPGroupAttributeID(groupDN, attrName))

	return resourceLDAPGroupAttributeRead(ctx, d, m)
}

func resourceLDAPGroupAttributeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)

	groupDN, attrName, err := parseResourceLDAPGroupAttributeID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	req := ldapv3.NewSearchRequest(
		groupDN,
		ldapv3.ScopeBaseObject,
		ldapv3.NeverDerefAliases,
		0, 0, false,
		"(objectclass=group)",
		[]string{attrName},
		nil,
	)

	sr, err := client.Conn.Search(req)
	if err != nil {
		ldapErr, ok := err.(*ldapv3.Error)
		if ok && ldapErr.ResultCode == ldapv3.LDAPResultNoSuchObject {
			// Group was deleted outside of Terraform
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed to read attribute %q on group %q: %w", attrName, groupDN, err))
	}

	if len(sr.Entries) == 0 {
		d.SetId("")
		return nil
	}

	// Find the attribute in the response
	var foundValues []string
	for _, attr := range sr.Entries[0].Attributes {
		if strings.EqualFold(attr.Name, attrName) {
			foundValues = attr.Values
			break
		}
	}

	if len(foundValues) == 0 {
		// Attribute no longer present on the group — mark as deleted
		d.SetId("")
		return nil
	}

	if err := d.Set("group_dn", groupDN); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("attribute_name", attrName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("attribute_values", foundValues); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceLDAPGroupAttributeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)

	groupDN := d.Get("group_dn").(string)
	attrName := d.Get("attribute_name").(string)

	if d.HasChange("attribute_values") {
		values := expandStringList(d.Get("attribute_values").([]interface{}))

		req := ldapv3.NewModifyRequest(groupDN, nil)
		req.Replace(attrName, values)

		if err := client.Conn.Modify(req); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update attribute %q on group %q: %w", attrName, groupDN, err))
		}
	}

	return resourceLDAPGroupAttributeRead(ctx, d, m)
}

func resourceLDAPGroupAttributeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)

	groupDN := d.Get("group_dn").(string)
	attrName := d.Get("attribute_name").(string)

	req := ldapv3.NewModifyRequest(groupDN, nil)
	req.Delete(attrName, []string{})

	if err := client.Conn.Modify(req); err != nil {
		// If the attribute is already gone, treat as success
		ldapErr, ok := err.(*ldapv3.Error)
		if ok && ldapErr.ResultCode == ldapv3.LDAPResultNoSuchAttribute {
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed to delete attribute %q on group %q: %w", attrName, groupDN, err))
	}

	return nil
}

func resourceLDAPGroupAttributeImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	groupDN, attrName, err := parseResourceLDAPGroupAttributeID(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("group_dn", groupDN)
	d.Set("attribute_name", attrName)

	return []*schema.ResourceData{d}, nil
}

// expandStringList converts []interface{} to []string.
func expandStringList(raw []interface{}) []string {
	result := make([]string, 0, len(raw))
	for _, v := range raw {
		result = append(result, v.(string))
	}
	return result
}
