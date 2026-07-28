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

func resourceLDAPGroupMember() *schema.Resource {
	return &schema.Resource{
		Description:   "`ldap_group_member` manages a single member of an LDAP group. Unlike `ldap_group.members`, this resource only controls its own member entry and does not affect other members added outside of Terraform.",
		CreateContext: resourceLDAPGroupMemberCreate,
		ReadContext:   resourceLDAPGroupMemberRead,
		DeleteContext: resourceLDAPGroupMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group_dn": {
				Description: "The full DN of the LDAP group to add the member to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"member_dn": {
				Description: "The full DN of the LDAP object (user, group, or computer) to add as a member.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceLDAPGroupMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)

	groupDN := d.Get("group_dn").(string)
	memberDN := d.Get("member_dn").(string)

	req := ldapv3.NewModifyRequest(groupDN, nil)
	req.Add("member", []string{memberDN})

	if err := client.Conn.Modify(req); err != nil {
		return diag.FromErr(fmt.Errorf("failed to add member %q to group %q: %w", memberDN, groupDN, err))
	}

	d.SetId(fmt.Sprintf("%s/%s", groupDN, memberDN))

	return resourceLDAPGroupMemberRead(ctx, d, m)
}

func resourceLDAPGroupMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid resource ID format, expected <group_dn>/<member_dn>, got: %s", d.Id())
	}
	groupDN := parts[0]
	memberDN := parts[1]

	attributes, err := client.ReadGroup(groupDN, 1500)
	if err != nil {
		ldapErr, ok := err.(*ldapv3.Error)
		if ok && ldapErr.ResultCode == 32 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed to read group %q: %w", groupDN, err))
	}

	found := false
	for _, v := range attributes["member"] {
		if strings.EqualFold(v, memberDN) {
			found = true
			break
		}
	}

	if !found {
		d.SetId("")
		return nil
	}

	if err := d.Set("group_dn", groupDN); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("member_dn", memberDN); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceLDAPGroupMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)

	groupDN := d.Get("group_dn").(string)
	memberDN := d.Get("member_dn").(string)

	req := ldapv3.NewModifyRequest(groupDN, nil)
	req.Delete("member", []string{memberDN})

	if err := client.Conn.Modify(req); err != nil {
		return diag.FromErr(fmt.Errorf("failed to remove member %q from group %q: %w", memberDN, groupDN, err))
	}

	return nil
}
