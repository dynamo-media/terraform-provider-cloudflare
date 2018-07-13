package cloudflare

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/cloudflare/cloudflare-go"
	"fmt"
	"strings"
	"log"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceCloudFlareZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudFlareZoneCreate,
		Delete: resourceCloudFlareZoneDelete,
		Update: resourceCloudFlareZoneUpdate,
		Read:   resourceCloudFlareZoneRead,
		Importer: &schema.ResourceImporter{
			State: resourceCloudFlareZoneImport,
		},

		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"jump_start": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
				ForceNew: true,
			},

			"organization_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"name_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"plan": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"free", "pro", "business", "enterprise"}, false),
			},
		},
	}
}

func resourceCloudFlareZoneCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	org := cloudflare.Organization{}
	orgId, orgIdOk := d.GetOk("organization_id")
	if orgIdOk {
		org.ID = orgId.(string)
	}
	client.ListOrganizations()
	z, err := client.CreateZone(d.Get("domain").(string), d.Get("jump_start").(bool), org)
	if err != nil {
		return fmt.Errorf("failed to create zone: %s", err)
	}

	if d.Get("plan").(string) != "free" {
		updateZoneSubscription(client, z.ID, d.Get("plan").(string))
	}

	d.SetId(z.ID)

	return resourceCloudFlareZoneRead(d, meta)
}

func resourceCloudFlareZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	zoneId := d.Id()
	z, err := client.ZoneDetails(zoneId)
	if err != nil {
		if strings.Contains(err.Error(), "Invalid zone identifier") {
			log.Printf("[INFO] zone %s not found", zoneId)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read zone: %s", err)
	}

	if z.Plan.Name == "Free Website" {
		d.Set("plan", "free")
	} else {
		s, err := client.ZoneSubscriptionDetails(z.ID)
		if err != nil {
			return fmt.Errorf("failed to read zone subscription details: %s", err)
		}
		d.Set("plan", s.RatePlan.ID)
	}

	d.Set("domain", z.Name)
	d.Set("organization_id", z.Owner.ID)
	d.Set("name_servers", z.NameServers)

	return nil
}

func resourceCloudFlareZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	zoneId := d.Id()
	_, err := client.DeleteZone(zoneId)
	if err != nil {
		return fmt.Errorf("error deleting zone: %s", err)
	}

	return nil
}

func resourceCloudFlareZoneImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceCloudFlareZoneRead(d, meta)
	if err != nil {
		return nil, err
	}

	d.Set("jump_start", d.Get("jump_start").(bool))

	return []*schema.ResourceData{d}, nil
}

func updateZoneSubscription(client *cloudflare.API, zoneId string, plan string) error {
	s := cloudflare.ZoneSubscriptionDetails{}
	s.RatePlan.ID = plan
	_, err := client.UpdateZoneSubscription(zoneId, s)
	if err != nil {
		return fmt.Errorf("failed to update zone: %s", err)
	}

	return nil
}

func resourceCloudFlareZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	plan, ok := d.GetOk("plan")
	if !ok {
		return nil
	}

	err := updateZoneSubscription(client, d.Id(), plan.(string))
	if err != nil {
		return err
	}

	return resourceCloudFlareZoneRead(d, meta)
}