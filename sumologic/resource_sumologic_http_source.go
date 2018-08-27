package sumologic

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nextgenhealthcare/sumologic-sdk-go"
	"strconv"
)

func resourceHTTPSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceHTTPSourceCreate,
		Read:   resourceHTTPSourceRead,
		Update: resourceHTTPSourceUpdate,
		Delete: resourceHTTPSourceDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"collector_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"source_type": &schema.Schema{
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
				Computed: true,
			},
			"message_per_request": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"multiline_processing_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceHTTPSourceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sumologic.Client)

	source, err := client.CreateHTTPSource(d.Get("collector_id").(int), sumologic.HTTPSource{
		Name:                       d.Get("name").(string),
		SourceType:                 d.Get("source_type").(string),
		Description:                d.Get("description").(string),
		Category:                   d.Get("category").(string),
		TimeZone:                   d.Get("timezone").(string),
		MessagePerRequest:          d.Get("message_per_request").(bool),
		MultilineProcessingEnabled: d.Get("multiline_processing_enabled").(bool),
	})
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(source.ID))

	return resourceHTTPSourceUpdate(d, m)
}

func resourceHTTPSourceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sumologic.Client)

	id, _ := strconv.Atoi(d.Id())
	source, _, err := client.GetHTTPSource(d.Get("collector_id").(int), id)
	if err == sumologic.ErrSourceNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	d.Set("name", source.Name)
	d.Set("source_type", source.SourceType)
	d.Set("description", source.Description)
	d.Set("category", source.Category)
	d.Set("timezone", source.TimeZone)
	d.Set("message_per_request", source.MessagePerRequest)
	d.Set("multiline_processing_enabled", source.MultilineProcessingEnabled)
	d.Set("url", source.Url)

	return nil
}

func resourceHTTPSourceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sumologic.Client)

	id, _ := strconv.Atoi(d.Id())
	source := sumologic.HTTPSource{
		ID:                         id,
		Name:                       d.Get("name").(string),
		SourceType:                 d.Get("source_type").(string),
		Description:                d.Get("description").(string),
		Category:                   d.Get("category").(string),
		TimeZone:                   d.Get("timezone").(string),
		MessagePerRequest:          d.Get("message_per_request").(bool),
		MultilineProcessingEnabled: d.Get("multiline_processing_enabled").(bool),
	}

	_, etag, _ := client.GetHTTPSource(d.Get("collector_id").(int), id)

	updatedSource, err := client.UpdateHTTPSource(d.Get("collector_id").(int), source, etag)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(updatedSource.ID))

	return resourceHTTPSourceRead(d, m)
}

func resourceHTTPSourceDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*sumologic.Client)

	id, _ := strconv.Atoi(d.Id())
	err := client.DeleteHTTPSource(d.Get("collector_id").(int), id)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
