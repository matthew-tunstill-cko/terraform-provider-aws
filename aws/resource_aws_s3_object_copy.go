package aws

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func resourceAwsS3ObjectCopy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsS3ObjectCopyCreate,
		Read:   resourceAwsS3BucketObjectRead,
		Update: resourceAwsS3ObjectCopyUpdate,
		Delete: resourceAwsS3BucketObjectDelete,

		Schema: map[string]*schema.Schema{
			"acl": {
				Type:          schema.TypeString,
				Default:       s3.ObjectCannedACLPrivate,
				Optional:      true,
				ValidateFunc:  validation.StringInSlice(s3.ObjectCannedACL_Values(), false),
				ConflictsWith: []string{"grant"},
			},

			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"cache_control": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"content_disposition": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"content_encoding": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"content_language": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"copy_if_match": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"copy_if_modified_since": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"copy_if_none_match": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"copy_if_unmodified_since": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"customer_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"customer_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"customer_key_md5": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expected_bucket_owner": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"expected_source_bucket_owner": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expires": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"grant": {
				Type:          schema.TypeSet,
				Optional:      true,
				Set:           grantHash,
				ConflictsWith: []string{"acl"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(s3.Type_Values(), false),
						},

						"uri": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"permissions": {
							Type:     schema.TypeSet,
							Required: true,
							Set:      schema.HashString,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									//write permission not valid here
									s3.PermissionFullControl,
									s3.PermissionRead,
									s3.PermissionReadAcp,
									s3.PermissionWriteAcp,
								}, false),
							},
						},
					},
				},
			},

			"key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"kms_encryption_context": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateArn,
				Sensitive:    true,
			},

			"kms_key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateArn,
				Sensitive:    true,
			},

			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"metadata": {
				Type:         schema.TypeMap,
				ValidateFunc: validateMetadataIsLowerCase,
				Optional:     true,
				Computed:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
			},

			"metadata_directive": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(s3.MetadataDirective_Values(), false),
			},

			"object_lock_legal_hold_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(s3.ObjectLockLegalHoldStatus_Values(), false),
			},

			"object_lock_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(s3.ObjectLockMode_Values(), false),
			},

			"object_lock_retain_until_date": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"request_charged": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"request_payer": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(s3.RequestPayer_Values(), false),
			},

			"server_side_encryption": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(s3.ServerSideEncryption_Values(), false),
			},

			"source": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"source_customer_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"source_customer_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"source_customer_key_md5": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"source_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"storage_class": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(s3.StorageClass_Values(), false),
			},

			"tagging_directive": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(s3.TaggingDirective_Values(), false),
			},

			"tags": tagsSchema(),

			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"website_redirect": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAwsS3ObjectCopyCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceAwsS3ObjectCopyCopy(d, meta)
}

func resourceAwsS3ObjectCopyUpdate(d *schema.ResourceData, meta interface{}) error {
	// if any of these exist, let the API decide whether to copy
	for _, key := range []string{
		"copy_if_match",
		"copy_if_modified_since",
		"copy_if_none_match",
		"copy_if_unmodified_since",
	} {
		if _, ok := d.GetOk(key); ok {
			return resourceAwsS3ObjectCopyCopy(d, meta)
		}
	}

	args := []string{
		"acl",
		"bucket",
		"cache_control",
		"content_disposition",
		"content_encoding",
		"content_language",
		"content_type",
		"copy_if_match",
		"copy_if_modified_since",
		"copy_if_none_match",
		"copy_if_unmodified_since",
		"customer_algorithm",
		"customer_key",
		"customer_key_md5",
		"expected_bucket_owner",
		"expected_source_bucket_owner",
		"expires",
		"grant_full_control",
		"grant_read",
		"grant_read_acp",
		"grant_write_acp",
		"key",
		"kms_encryption_context",
		"kms_key_id",
		"metadata",
		"metadata_directive",
		"object_lock_legal_hold_status",
		"object_lock_mode",
		"object_lock_retain_until_date",
		"request_payer",
		"server_side_encryption",
		"source",
		"source_customer_algorithm",
		"source_customer_key",
		"source_customer_key_md5",
		"storage_class",
		"tagging_directive",
		"tags",
		"website_redirect",
	}
	if d.HasChanges(args...) {
		return resourceAwsS3ObjectCopyCopy(d, meta)
	}

	return nil
}

var resourceAWSS3ObjectCopyCountApiCalls = 0

func resourceAwsS3ObjectCopyCopy(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).s3conn

	input := &s3.CopyObjectInput{
		Bucket:     aws.String(d.Get("bucket").(string)),
		Key:        aws.String(d.Get("key").(string)),
		CopySource: aws.String(url.QueryEscape(d.Get("source").(string))),
	}

	if v, ok := d.GetOk("acl"); ok {
		input.ACL = aws.String(v.(string))
	}

	if v, ok := d.GetOk("cache_control"); ok {
		input.CacheControl = aws.String(v.(string))
	}

	if v, ok := d.GetOk("content_disposition"); ok {
		input.ContentDisposition = aws.String(v.(string))
	}

	if v, ok := d.GetOk("content_encoding"); ok {
		input.ContentEncoding = aws.String(v.(string))
	}

	if v, ok := d.GetOk("content_language"); ok {
		input.ContentLanguage = aws.String(v.(string))
	}

	if v, ok := d.GetOk("content_type"); ok {
		input.ContentType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("copy_if_match"); ok {
		input.CopySourceIfMatch = aws.String(v.(string))
	}

	if v, ok := d.GetOk("copy_if_modified_since"); ok {
		input.CopySourceIfModifiedSince = expandS3ObjectDate(v.(string))
	}

	if v, ok := d.GetOk("copy_if_none_match"); ok {
		input.CopySourceIfNoneMatch = aws.String(v.(string))
	}

	if v, ok := d.GetOk("copy_if_unmodified_since"); ok {
		input.CopySourceIfUnmodifiedSince = expandS3ObjectDate(v.(string))
	}

	if v, ok := d.GetOk("customer_algorithm"); ok {
		input.SSECustomerAlgorithm = aws.String(v.(string))
	}

	if v, ok := d.GetOk("customer_key"); ok {
		input.SSECustomerKey = aws.String(v.(string))
	}

	if v, ok := d.GetOk("customer_key_md5"); ok {
		input.SSECustomerKeyMD5 = aws.String(v.(string))
	}

	if v, ok := d.GetOk("expected_bucket_owner"); ok {
		input.ExpectedBucketOwner = aws.String(v.(string))
	}

	if v, ok := d.GetOk("expected_source_bucket_owner"); ok {
		input.ExpectedSourceBucketOwner = aws.String(v.(string))
	}

	if v, ok := d.GetOk("expires"); ok {
		input.Expires = expandS3ObjectDate(v.(string))
	}

	if err := resourceAwsS3ObjectCopyIncludeGrants(input, d); err != nil {
		return err
	}

	if v, ok := d.GetOk("kms_encryption_context"); ok {
		input.SSEKMSEncryptionContext = aws.String(v.(string))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		input.SSEKMSKeyId = aws.String(v.(string))
		input.ServerSideEncryption = aws.String(s3.ServerSideEncryptionAwsKms)
	}

	if v, ok := d.GetOk("metadata"); ok {
		input.Metadata = stringMapToPointers(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("metadata_directive"); ok {
		input.MetadataDirective = aws.String(v.(string))
	}

	if v, ok := d.GetOk("object_lock_legal_hold_status"); ok {
		input.ObjectLockLegalHoldStatus = aws.String(v.(string))
	}

	if v, ok := d.GetOk("object_lock_mode"); ok {
		input.ObjectLockMode = aws.String(v.(string))
	}

	if v, ok := d.GetOk("object_lock_retain_until_date"); ok {
		input.ObjectLockRetainUntilDate = expandS3ObjectDate(v.(string))
	}

	if v, ok := d.GetOk("request_payer"); ok {
		input.RequestPayer = aws.String(v.(string))
	}

	if v, ok := d.GetOk("server_side_encryption"); ok {
		input.ServerSideEncryption = aws.String(v.(string))
	}

	if v, ok := d.GetOk("source_customer_algorithm"); ok {
		input.CopySourceSSECustomerAlgorithm = aws.String(v.(string))
	}

	if v, ok := d.GetOk("source_customer_key"); ok {
		input.CopySourceSSECustomerKey = aws.String(v.(string))
	}

	if v, ok := d.GetOk("source_customer_key_md5"); ok {
		input.CopySourceSSECustomerKeyMD5 = aws.String(v.(string))
	}

	if v, ok := d.GetOk("storage_class"); ok {
		input.StorageClass = aws.String(v.(string))
	}

	if v, ok := d.GetOk("tagging_directive"); ok {
		input.TaggingDirective = aws.String(v.(string))
	}

	if v := d.Get("tags").(map[string]interface{}); len(v) > 0 {
		// The tag-set must be encoded as URL Query parameters.
		input.Tagging = aws.String(keyvaluetags.New(v).IgnoreAws().UrlEncode())
	}

	if v, ok := d.GetOk("website_redirect"); ok {
		input.WebsiteRedirectLocation = aws.String(v.(string))
	}

	output, err := conn.CopyObject(input)
	if err != nil {
		return fmt.Errorf("Error copying S3 object (bucket: %s; key: %s; source: %s): %s", aws.StringValue(input.Bucket), aws.StringValue(input.Key), aws.StringValue(input.CopySource), err)
	}
	resourceAWSS3ObjectCopyCountApiCalls++
	if resourceAWSS3ObjectCopyCountApiCalls > 1 {
		return fmt.Errorf("Copying S3 object 2 times! (bucket: %s; key: %s; source: %s): %s", aws.StringValue(input.Bucket), aws.StringValue(input.Key), aws.StringValue(input.CopySource), err)
	}

	d.Set("customer_algorithm", output.SSECustomerAlgorithm)
	d.Set("customer_key_md5", output.SSECustomerKeyMD5)

	if output.CopyObjectResult != nil {
		d.Set("etag", strings.Trim(aws.StringValue(output.CopyObjectResult.ETag), `"`))
		d.Set("last_modified", flattenS3ObjectDate(output.CopyObjectResult.LastModified))
	}

	d.Set("expiration", output.Expiration)
	d.Set("kms_encryption_context", output.SSEKMSEncryptionContext)
	d.Set("kms_key_id", output.SSEKMSKeyId)
	d.Set("request_charged", output.RequestCharged)
	d.Set("server_side_encryption", output.ServerSideEncryption)
	d.Set("source_version_id", output.CopySourceVersionId)
	d.Set("version_id", output.VersionId)

	d.SetId(d.Get("key").(string))
	return resourceAwsS3BucketObjectRead(d, meta)
}

func resourceAwsS3ObjectCopyIncludeGrants(input *s3.CopyObjectInput, d *schema.ResourceData) error {
	grants := d.Get("grant").(*schema.Set).List()

	if len(grants) > 0 {
		grantFullControl := make([]*s3.Grantee, 0)
		grantRead := make([]*s3.Grantee, 0)
		grantReadACP := make([]*s3.Grantee, 0)
		grantWriteACP := make([]*s3.Grantee, 0)

		for _, grant := range grants {
			grantMap := grant.(map[string]interface{})

			for _, perm := range grantMap["permissions"].(*schema.Set).List() {
				ge := &s3.Grantee{}

				if i, ok := grantMap["email"].(string); ok && i != "" {
					ge.SetEmailAddress(i)
				}

				if i, ok := grantMap["id"].(string); ok && i != "" {
					ge.SetID(i)
				}

				if t, ok := grantMap["type"].(string); ok && t != "" {
					ge.SetType(t)
				}

				if u, ok := grantMap["uri"].(string); ok && u != "" {
					ge.SetURI(u)
				}

				switch perm.(string) {
				case s3.PermissionFullControl:
					grantFullControl = append(grantFullControl, ge)
				case s3.PermissionRead:
					grantRead = append(grantRead, ge)
				case s3.PermissionReadAcp:
					grantReadACP = append(grantReadACP, ge)
				case s3.PermissionWriteAcp:
					grantWriteACP = append(grantWriteACP, ge)
				}
			}
		}

		if len(grantFullControl) > 0 {
			input.GrantFullControl = aws.String(flattenS3ObjectCopyGrantees(grantFullControl))
		}

		if len(grantRead) > 0 {
			input.GrantRead = aws.String(flattenS3ObjectCopyGrantees(grantRead))
		}

		if len(grantReadACP) > 0 {
			input.GrantReadACP = aws.String(flattenS3ObjectCopyGrantees(grantReadACP))
		}

		if len(grantWriteACP) > 0 {
			input.GrantWriteACP = aws.String(flattenS3ObjectCopyGrantees(grantWriteACP))
		}
	}
	return nil
}

func flattenS3ObjectCopyGrantees(grantees []*s3.Grantee) string {
	//"GrantFullControl": "emailaddress=user1@example.com,emailaddress=user2@example.com",
	//"GrantRead": "uri=http://acs.amazonaws.com/groups/global/AllUsers",
	//"GrantFullControl": "id=examplee7a2f25102679df27bb0ae12b3f85be6f290b936c4393484",
	//"GrantWrite": "uri=http://acs.amazonaws.com/groups/s3/LogDelivery"
	flatGrants := make([]string, 0)
	for _, grantee := range grantees {
		switch *grantee.Type {
		case s3.TypeAmazonCustomerByEmail:
			flatGrants = append(flatGrants, fmt.Sprintf("emailaddress=%s", *grantee.EmailAddress))
		case s3.TypeCanonicalUser:
			flatGrants = append(flatGrants, fmt.Sprintf("id=%s", *grantee.ID))
		case s3.TypeGroup:
			flatGrants = append(flatGrants, fmt.Sprintf("uri=%s", *grantee.URI))
		}
	}

	return strings.Join(flatGrants, ",")
}
