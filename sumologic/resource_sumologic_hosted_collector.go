package sumologic

import (
	"github.com/brandonstevens/sumologic-sdk-go"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
)

func resourceHostedCollector() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostedCollectorCreate,
		Read:   resourceHostedCollectorRead,
		Update: resourceHostedCollectorUpdate,
		Delete: resourceHostedCollectorDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"category": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"timezone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UTC",
			},
		},
	}
}

func resourceHostedCollectorCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sumologic.Client)

	collector, err := client.CreateHostedCollector(sumologic.Collector{
		Name:          d.Get("name").(string),
		CollectorType: "Hosted",
		Description:   d.Get("description").(string),
		Category:      d.Get("category").(string),
		TimeZone:      d.Get("timezone").(string),
	})
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(collector.ID))

	return resourceHostedCollectorUpdate(d, m)
}

func resourceHostedCollectorRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sumologic.Client)

	id, _ := strconv.Atoi(d.Id())
	collector, _, err := client.GetHostedCollector(id)
	if err == sumologic.ErrCollectorNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	d.Set("name", collector.Name)
	d.Set("description", collector.Description)
	d.Set("category", collector.Category)
	d.Set("timezone", collector.TimeZone)

	return nil
}

func resourceHostedCollectorUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sumologic.Client)

	id, _ := strconv.Atoi(d.Id())
	collector := sumologic.Collector{
		ID:            id,
		Name:          d.Get("name").(string),
		CollectorType: "Hosted",
		Description:   d.Get("description").(string),
		Category:      d.Get("category").(string),
		TimeZone:      d.Get("timezone").(string),
	}

	_, etag, _ := client.GetHostedCollector(id)

	updatedCollector, err := client.UpdateHostedCollector(collector, etag)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(updatedCollector.ID))

	return resourceHostedCollectorRead(d, m)
}

func resourceHostedCollectorDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*sumologic.Client)

	id, _ := strconv.Atoi(d.Id())
	err := client.DeleteHostedCollector(id)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
