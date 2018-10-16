package sumologic

import (
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nextgenhealthcare/sumologic-sdk-go"
	"strconv"
	"time"
)

func resourceAWSCloudTrailSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAWSCloudTrailSourceCreate,
		Read:   resourceAWSCloudTrailSourceRead,
		Update: resourceAWSCloudTrailSourceUpdate,
		Delete: resourceAWSCloudTrailSourceDelete,

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
			"scan_interval": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"content_type": &schema.Schema{
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
			"paused": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"cutoff_relative_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"third_party_ref": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resources": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_type": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"path": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Required: true,
												},
												"bucket_name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"path_expression": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"authentication": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Required: true,
												},
												"role_arn": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAWSCloudTrailSourceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sumologic.Client)

	s := sumologic.AWSCloudTrailSource{
		Name:               d.Get("name").(string),
		SourceType:         d.Get("source_type").(string),
		ScanInterval:       d.Get("scan_interval").(int),
		ContentType:        d.Get("content_type").(string),
		Description:        d.Get("description").(string),
		Category:           d.Get("category").(string),
		TimeZone:           d.Get("timezone").(string),
		Paused:             d.Get("paused").(bool),
		CutoffRelativeTime: d.Get("cutoff_relative_time").(string),
		ThirdPartyRef: sumologic.AWSBucketThirdPartyRef{
			Resources: make([]sumologic.AWSBucketResource, 0),
		},
	}
	a := sumologic.AWSBucketResource{
		ServiceType: d.Get("third_party_ref.0.resources.0.service_type").(string),
		Path: sumologic.AWSBucketPath{
			Type:           d.Get("third_party_ref.0.resources.0.path.0.type").(string),
			BucketName:     d.Get("third_party_ref.0.resources.0.path.0.bucket_name").(string),
			PathExpression: d.Get("third_party_ref.0.resources.0.path.0.path_expression").(string),
		},
		Authentication: sumologic.AWSBucketAuthentication{
			Type:    d.Get("third_party_ref.0.resources.0.authentication.0.type").(string),
			RoleARN: d.Get("third_party_ref.0.resources.0.authentication.0.role_arn").(string),
		},
	}
	s.ThirdPartyRef.Resources = append(s.ThirdPartyRef.Resources, a)

	// Retry due to IAM eventual consistency
	var out *sumologic.AWSCloudTrailSource
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		var err error
		out, err = client.CreateAWSCloudTrailSource(d.Get("collector_id").(int), s)

		if err == sumologic.ErrAwsAuthenticationError {
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})
	if err != nil {
		return err
	}

	source := *out

	d.SetId(strconv.Itoa(source.ID))

	return resourceAWSCloudTrailSourceUpdate(d, m)
}

func resourceAWSCloudTrailSourceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*sumologic.Client)

	id, _ := strconv.Atoi(d.Id())
	source, _, err := client.GetAWSCloudTrailSource(d.Get("collector_id").(int), id)
	if err == sumologic.ErrSourceNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	d.Set("name", source.Name)
	d.Set("source_type", source.SourceType)
	d.Set("scan_interval", source.ScanInterval)
	d.Set("content_type", source.ContentType)
	d.Set("description", source.Description)
	d.Set("category", source.Category)
	d.Set("timezone", source.TimeZone)
	d.Set("paused", source.Paused)
	d.Set("third_party_ref.0.resources.0.service_type", source.ThirdPartyRef.Resources[0].ServiceType)
	d.Set("third_party_ref.0.resources.0.path.0.type", source.ThirdPartyRef.Resources[0].Path.Type)
	d.Set("third_party_ref.0.resources.0.path.0.bucket_name", source.ThirdPartyRef.Resources[0].Path.BucketName)
	d.Set("third_party_ref.0.resources.0.path.0.path_expression", source.ThirdPartyRef.Resources[0].Path.PathExpression)
	d.Set("third_party_ref.0.resources.0.authentication.0.type", source.ThirdPartyRef.Resources[0].Authentication.Type)
	d.Set("third_party_ref.0.resources.0.authentication.0.role_arn", source.ThirdPartyRef.Resources[0].Authentication.RoleARN)

	return nil
}

func resourceAWSCloudTrailSourceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*sumologic.Client)

	id, _ := strconv.Atoi(d.Id())
	source := sumologic.AWSCloudTrailSource{
		ID:                 id,
		Name:               d.Get("name").(string),
		SourceType:         d.Get("source_type").(string),
		ScanInterval:       d.Get("scan_interval").(int),
		ContentType:        d.Get("content_type").(string),
		Description:        d.Get("description").(string),
		Category:           d.Get("category").(string),
		TimeZone:           d.Get("timezone").(string),
		Paused:             d.Get("paused").(bool),
		CutoffRelativeTime: d.Get("cutoff_relative_time").(string),
		ThirdPartyRef: sumologic.AWSBucketThirdPartyRef{
			Resources: make([]sumologic.AWSBucketResource, 0),
		},
	}
	a := sumologic.AWSBucketResource{
		ServiceType: d.Get("third_party_ref.0.resources.0.service_type").(string),
		Path: sumologic.AWSBucketPath{
			Type:           d.Get("third_party_ref.0.resources.0.path.0.type").(string),
			BucketName:     d.Get("third_party_ref.0.resources.0.path.0.bucket_name").(string),
			PathExpression: d.Get("third_party_ref.0.resources.0.path.0.path_expression").(string),
		},
		Authentication: sumologic.AWSBucketAuthentication{
			Type:    d.Get("third_party_ref.0.resources.0.authentication.0.type").(string),
			RoleARN: d.Get("third_party_ref.0.resources.0.authentication.0.role_arn").(string),
		},
	}
	source.ThirdPartyRef.Resources = append(source.ThirdPartyRef.Resources, a)

	_, etag, _ := client.GetAWSCloudTrailSource(d.Get("collector_id").(int), id)

	// Retry due to IAM eventual consistency
	var out *sumologic.AWSCloudTrailSource
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		var err error
		out, err = client.UpdateAWSCloudTrailSource(d.Get("collector_id").(int), source, etag)

		if err == sumologic.ErrAwsAuthenticationError {
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})
	if err != nil {
		return err
	}

	updatedSource := *out

	d.SetId(strconv.Itoa(updatedSource.ID))

	return resourceAWSCloudTrailSourceRead(d, m)
}

func resourceAWSCloudTrailSourceDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*sumologic.Client)

	id, _ := strconv.Atoi(d.Id())
	err := client.DeleteAWSCloudTrailSource(d.Get("collector_id").(int), id)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
