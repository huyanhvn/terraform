package aws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
)

func TestValidateEcrRepositoryName(t *testing.T) {
	validNames := []string{
		"nginx-web-app",
		"project-a/nginx-web-app",
		"domain.ltd/nginx-web-app",
		"3chosome-thing.com/01different-pattern",
		"0123456789/999999999",
		"double/forward/slash",
		"000000000000000",
	}
	for _, v := range validNames {
		_, errors := validateEcrRepositoryName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid ECR repository name: %q", v, errors)
		}
	}

	invalidNames := []string{
		// length > 256
		"3cho_some-thing.com/01different.-_pattern01different.-_pattern01diff" +
			"erent.-_pattern01different.-_pattern01different.-_pattern01different" +
			".-_pattern01different.-_pattern01different.-_pattern01different.-_pa" +
			"ttern01different.-_pattern01different.-_pattern234567",
		// length < 2
		"i",
		"special@character",
		"different+special=character",
		"double//slash",
		"double..dot",
		"/slash-at-the-beginning",
		"slash-at-the-end/",
	}
	for _, v := range invalidNames {
		_, errors := validateEcrRepositoryName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid ECR repository name", v)
		}
	}
}

func TestValidateCloudWatchEventRuleName(t *testing.T) {
	validNames := []string{
		"HelloWorl_d",
		"hello-world",
		"hello.World0125",
	}
	for _, v := range validNames {
		_, errors := validateCloudWatchEventRuleName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid CW event rule name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"special@character",
		"slash/in-the-middle",
		// Length > 64
		"TooLooooooooooooooooooooooooooooooooooooooooooooooooooooooongName",
	}
	for _, v := range invalidNames {
		_, errors := validateCloudWatchEventRuleName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid CW event rule name", v)
		}
	}
}

func TestValidateLambdaFunctionName(t *testing.T) {
	validNames := []string{
		"arn:aws:lambda:us-west-2:123456789012:function:ThumbNail",
		"FunctionName",
		"function-name",
	}
	for _, v := range validNames {
		_, errors := validateLambdaFunctionName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Lambda function name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"/FunctionNameWithSlash",
		"function.name.with.dots",
		// length > 140
		"arn:aws:lambda:us-west-2:123456789012:function:TooLoooooo" +
			"ooooooooooooooooooooooooooooooooooooooooooooooooooooooo" +
			"ooooooooooooooooongFunctionName",
	}
	for _, v := range invalidNames {
		_, errors := validateLambdaFunctionName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Lambda function name", v)
		}
	}
}

func TestValidateLambdaQualifier(t *testing.T) {
	validNames := []string{
		"123",
		"prod",
		"PROD",
		"MyTestEnv",
		"$LATEST",
	}
	for _, v := range validNames {
		_, errors := validateLambdaQualifier(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Lambda function qualifier: %q", v, errors)
		}
	}

	invalidNames := []string{
		// No ARNs allowed
		"arn:aws:lambda:us-west-2:123456789012:function:prod",
		// length > 128
		"TooLooooooooooooooooooooooooooooooooooooooooooooooooooo" +
			"ooooooooooooooooooooooooooooooooooooooooooooooooooo" +
			"oooooooooooongQualifier",
	}
	for _, v := range invalidNames {
		_, errors := validateLambdaQualifier(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Lambda function qualifier", v)
		}
	}
}

func TestValidateLambdaPermissionAction(t *testing.T) {
	validNames := []string{
		"lambda:*",
		"lambda:InvokeFunction",
		"*",
	}
	for _, v := range validNames {
		_, errors := validateLambdaPermissionAction(v, "action")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Lambda permission action: %q", v, errors)
		}
	}

	invalidNames := []string{
		"yada",
		"lambda:123",
		"*:*",
		"lambda:Invoke*",
	}
	for _, v := range invalidNames {
		_, errors := validateLambdaPermissionAction(v, "action")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Lambda permission action", v)
		}
	}
}

func TestValidateAwsAccountId(t *testing.T) {
	validNames := []string{
		"123456789012",
		"999999999999",
	}
	for _, v := range validNames {
		_, errors := validateAwsAccountId(v, "account_id")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid AWS Account ID: %q", v, errors)
		}
	}

	invalidNames := []string{
		"12345678901",   // too short
		"1234567890123", // too long
		"invalid",
		"x123456789012",
	}
	for _, v := range invalidNames {
		_, errors := validateAwsAccountId(v, "account_id")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid AWS Account ID", v)
		}
	}
}

func TestValidateArn(t *testing.T) {
	v := ""
	_, errors := validateArn(v, "arn")
	if len(errors) != 0 {
		t.Fatalf("%q should not be validated as an ARN: %q", v, errors)
	}

	validNames := []string{
		"arn:aws:elasticbeanstalk:us-east-1:123456789012:environment/My App/MyEnvironment", // Beanstalk
		"arn:aws:iam::123456789012:user/David",                                             // IAM User
		"arn:aws:rds:eu-west-1:123456789012:db:mysql-db",                                   // RDS
		"arn:aws:s3:::my_corporate_bucket/exampleobject.png",                               // S3 object
		"arn:aws:events:us-east-1:319201112229:rule/rule_name",                             // CloudWatch Rule
		"arn:aws:lambda:eu-west-1:319201112229:function:myCustomFunction",                  // Lambda function
		"arn:aws:lambda:eu-west-1:319201112229:function:myCustomFunction:Qualifier",        // Lambda func qualifier
	}
	for _, v := range validNames {
		_, errors := validateArn(v, "arn")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid ARN: %q", v, errors)
		}
	}

	invalidNames := []string{
		"arn",
		"123456789012",
		"arn:aws",
		"arn:aws:logs",
		"arn:aws:logs:region:*:*",
	}
	for _, v := range invalidNames {
		_, errors := validateArn(v, "arn")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid ARN", v)
		}
	}
}

func TestValidatePolicyStatementId(t *testing.T) {
	validNames := []string{
		"YadaHereAndThere",
		"Valid-5tatement_Id",
		"1234",
	}
	for _, v := range validNames {
		_, errors := validatePolicyStatementId(v, "statement_id")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Statement ID: %q", v, errors)
		}
	}

	invalidNames := []string{
		"Invalid/StatementId/with/slashes",
		"InvalidStatementId.with.dots",
		// length > 100
		"TooooLoooooooooooooooooooooooooooooooooooooooooooo" +
			"ooooooooooooooooooooooooooooooooooooooooStatementId",
	}
	for _, v := range invalidNames {
		_, errors := validatePolicyStatementId(v, "statement_id")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Statement ID", v)
		}
	}
}

func TestValidateCIDRNetworkAddress(t *testing.T) {
	cases := []struct {
		CIDR              string
		ExpectedErrSubstr string
	}{
		{"notacidr", `must contain a valid CIDR`},
		{"10.0.1.0/16", `must contain a valid network CIDR`},
		{"10.0.1.0/24", ``},
	}

	for i, tc := range cases {
		_, errs := validateCIDRNetworkAddress(tc.CIDR, "foo")
		if tc.ExpectedErrSubstr == "" {
			if len(errs) != 0 {
				t.Fatalf("%d/%d: Expected no error, got errs: %#v",
					i+1, len(cases), errs)
			}
		} else {
			if len(errs) != 1 {
				t.Fatalf("%d/%d: Expected 1 err containing %q, got %d errs",
					i+1, len(cases), tc.ExpectedErrSubstr, len(errs))
			}
			if !strings.Contains(errs[0].Error(), tc.ExpectedErrSubstr) {
				t.Fatalf("%d/%d: Expected err: %q, to include %q",
					i+1, len(cases), errs[0], tc.ExpectedErrSubstr)
			}
		}
	}
}

func TestValidateHTTPMethod(t *testing.T) {
	type testCases struct {
		Value    string
		ErrCount int
	}

	invalidCases := []testCases{
		{
			Value:    "incorrect",
			ErrCount: 1,
		},
		{
			Value:    "delete",
			ErrCount: 1,
		},
	}

	for _, tc := range invalidCases {
		_, errors := validateHTTPMethod(tc.Value, "http_method")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %q to trigger a validation error.", tc.Value)
		}
	}

	validCases := []testCases{
		{
			Value:    "ANY",
			ErrCount: 0,
		},
		{
			Value:    "DELETE",
			ErrCount: 0,
		},
		{
			Value:    "OPTIONS",
			ErrCount: 0,
		},
	}

	for _, tc := range validCases {
		_, errors := validateHTTPMethod(tc.Value, "http_method")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %q not to trigger a validation error.", tc.Value)
		}
	}
}

func TestValidateLogMetricFilterName(t *testing.T) {
	validNames := []string{
		"YadaHereAndThere",
		"Valid-5Metric_Name",
		"This . is also %% valid@!)+(",
		"1234",
		strings.Repeat("W", 512),
	}
	for _, v := range validNames {
		_, errors := validateLogMetricFilterName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Log Metric Filter Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"Here is a name with: colon",
		"and here is another * invalid name",
		"*",
		// length > 512
		strings.Repeat("W", 513),
	}
	for _, v := range invalidNames {
		_, errors := validateLogMetricFilterName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Log Metric Filter Name", v)
		}
	}
}

func TestValidateLogMetricTransformationName(t *testing.T) {
	validNames := []string{
		"YadaHereAndThere",
		"Valid-5Metric_Name",
		"This . is also %% valid@!)+(",
		"1234",
		"",
		strings.Repeat("W", 255),
	}
	for _, v := range validNames {
		_, errors := validateLogMetricFilterTransformationName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Log Metric Filter Transformation Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"Here is a name with: colon",
		"and here is another * invalid name",
		"also $ invalid",
		"*",
		// length > 255
		strings.Repeat("W", 256),
	}
	for _, v := range invalidNames {
		_, errors := validateLogMetricFilterTransformationName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Log Metric Filter Transformation Name", v)
		}
	}
}

func TestValidateLogGroupName(t *testing.T) {
	validNames := []string{
		"ValidLogGroupName",
		"ValidLogGroup.Name",
		"valid/Log-group",
		"1234",
		"YadaValid#0123",
		"Also_valid-name",
		strings.Repeat("W", 512),
	}
	for _, v := range validNames {
		_, errors := validateLogGroupName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Log Metric Filter Transformation Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"Here is a name with: colon",
		"and here is another * invalid name",
		"also $ invalid",
		"This . is also %% invalid@!)+(",
		"*",
		"",
		// length > 512
		strings.Repeat("W", 513),
	}
	for _, v := range invalidNames {
		_, errors := validateLogGroupName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Log Metric Filter Transformation Name", v)
		}
	}
}

func TestValidateS3BucketLifecycleTimestamp(t *testing.T) {
	validDates := []string{
		"2016-01-01",
		"2006-01-02",
	}

	for _, v := range validDates {
		_, errors := validateS3BucketLifecycleTimestamp(v, "date")
		if len(errors) != 0 {
			t.Fatalf("%q should be valid date: %q", v, errors)
		}
	}

	invalidDates := []string{
		"Jan 01 2016",
		"20160101",
	}

	for _, v := range invalidDates {
		_, errors := validateS3BucketLifecycleTimestamp(v, "date")
		if len(errors) == 0 {
			t.Fatalf("%q should be invalid date", v)
		}
	}
}

func TestValidateS3BucketLifecycleStorageClass(t *testing.T) {
	validStorageClass := []string{
		"STANDARD_IA",
		"GLACIER",
	}

	for _, v := range validStorageClass {
		_, errors := validateS3BucketLifecycleStorageClass(v, "storage_class")
		if len(errors) != 0 {
			t.Fatalf("%q should be valid storage class: %q", v, errors)
		}
	}

	invalidStorageClass := []string{
		"STANDARD",
		"1234",
	}
	for _, v := range invalidStorageClass {
		_, errors := validateS3BucketLifecycleStorageClass(v, "storage_class")
		if len(errors) == 0 {
			t.Fatalf("%q should be invalid storage class", v)
		}
	}
}

func TestValidateS3BucketReplicationRuleId(t *testing.T) {
	validId := []string{
		"YadaHereAndThere",
		"Valid-5Rule_ID",
		"This . is also %% valid@!)+*(:ID",
		"1234",
		strings.Repeat("W", 255),
	}
	for _, v := range validId {
		_, errors := validateS3BucketReplicationRuleId(v, "id")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid lifecycle rule id: %q", v, errors)
		}
	}

	invalidId := []string{
		// length > 255
		strings.Repeat("W", 256),
	}
	for _, v := range invalidId {
		_, errors := validateS3BucketReplicationRuleId(v, "id")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid replication configuration rule id", v)
		}
	}
}

func TestValidateS3BucketReplicationRulePrefix(t *testing.T) {
	validId := []string{
		"YadaHereAndThere",
		"Valid-5Rule_ID",
		"This . is also %% valid@!)+*(:ID",
		"1234",
		strings.Repeat("W", 1024),
	}
	for _, v := range validId {
		_, errors := validateS3BucketReplicationRulePrefix(v, "id")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid lifecycle rule id: %q", v, errors)
		}
	}

	invalidId := []string{
		// length > 1024
		strings.Repeat("W", 1025),
	}
	for _, v := range invalidId {
		_, errors := validateS3BucketReplicationRulePrefix(v, "id")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid replication configuration rule id", v)
		}
	}
}

func TestValidateS3BucketReplicationDestinationStorageClass(t *testing.T) {
	validStorageClass := []string{
		s3.StorageClassStandard,
		s3.StorageClassStandardIa,
		s3.StorageClassReducedRedundancy,
	}

	for _, v := range validStorageClass {
		_, errors := validateS3BucketReplicationDestinationStorageClass(v, "storage_class")
		if len(errors) != 0 {
			t.Fatalf("%q should be valid storage class: %q", v, errors)
		}
	}

	invalidStorageClass := []string{
		"FOO",
		"1234",
	}
	for _, v := range invalidStorageClass {
		_, errors := validateS3BucketReplicationDestinationStorageClass(v, "storage_class")
		if len(errors) == 0 {
			t.Fatalf("%q should be invalid storage class", v)
		}
	}
}

func TestValidateS3BucketReplicationRuleStatus(t *testing.T) {
	validRuleStatuses := []string{
		s3.ReplicationRuleStatusEnabled,
		s3.ReplicationRuleStatusDisabled,
	}

	for _, v := range validRuleStatuses {
		_, errors := validateS3BucketReplicationRuleStatus(v, "status")
		if len(errors) != 0 {
			t.Fatalf("%q should be valid rule status: %q", v, errors)
		}
	}

	invalidRuleStatuses := []string{
		"FOO",
		"1234",
	}
	for _, v := range invalidRuleStatuses {
		_, errors := validateS3BucketReplicationRuleStatus(v, "status")
		if len(errors) == 0 {
			t.Fatalf("%q should be invalid rule status", v)
		}
	}
}

func TestValidateS3BucketLifecycleRuleId(t *testing.T) {
	validId := []string{
		"YadaHereAndThere",
		"Valid-5Rule_ID",
		"This . is also %% valid@!)+*(:ID",
		"1234",
		strings.Repeat("W", 255),
	}
	for _, v := range validId {
		_, errors := validateS3BucketLifecycleRuleId(v, "id")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid lifecycle rule id: %q", v, errors)
		}
	}

	invalidId := []string{
		// length > 255
		strings.Repeat("W", 256),
	}
	for _, v := range invalidId {
		_, errors := validateS3BucketLifecycleRuleId(v, "id")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid lifecycle rule id", v)
		}
	}
}

func TestValidateIntegerInRange(t *testing.T) {
	validIntegers := []int{-259, 0, 1, 5, 999}
	min := -259
	max := 999
	for _, v := range validIntegers {
		_, errors := validateIntegerInRange(min, max)(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be an integer in range (%d, %d): %q", v, min, max, errors)
		}
	}

	invalidIntegers := []int{-260, -99999, 1000, 25678}
	for _, v := range invalidIntegers {
		_, errors := validateIntegerInRange(min, max)(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an integer outside range (%d, %d)", v, min, max)
		}
	}
}

func TestResourceAWSElastiCacheClusterIdValidation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "tEsting",
			ErrCount: 1,
		},
		{
			Value:    "t.sting",
			ErrCount: 1,
		},
		{
			Value:    "t--sting",
			ErrCount: 1,
		},
		{
			Value:    "1testing",
			ErrCount: 1,
		},
		{
			Value:    "testing-",
			ErrCount: 1,
		},
		{
			Value:    randomString(65),
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateElastiCacheClusterId(tc.Value, "aws_elasticache_cluster_cluster_id")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the ElastiCache Cluster cluster_id to trigger a validation error")
		}
	}
}

func TestValidateDbEventSubscriptionName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"valid02-name",
		"Valid-Name1",
	}
	for _, v := range validNames {
		_, errors := validateDbEventSubscriptionName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid RDS Event Subscription Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"Here is a name with: colon",
		"and here is another * invalid name",
		"also $ invalid",
		"This . is also %% invalid@!)+(",
		"*",
		"",
		" ",
		"_",
		// length > 255
		strings.Repeat("W", 256),
	}
	for _, v := range invalidNames {
		_, errors := validateDbEventSubscriptionName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid RDS Event Subscription Name", v)
		}
	}
}

func TestValidateJsonString(t *testing.T) {
	type testCases struct {
		Value    string
		ErrCount int
	}

	invalidCases := []testCases{
		{
			Value:    `{0:"1"}`,
			ErrCount: 1,
		},
		{
			Value:    `{'abc':1}`,
			ErrCount: 1,
		},
		{
			Value:    `{"def":}`,
			ErrCount: 1,
		},
		{
			Value:    `{"xyz":[}}`,
			ErrCount: 1,
		},
	}

	for _, tc := range invalidCases {
		_, errors := validateJsonString(tc.Value, "json")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %q to trigger a validation error.", tc.Value)
		}
	}

	validCases := []testCases{
		{
			Value:    ``,
			ErrCount: 0,
		},
		{
			Value:    `{}`,
			ErrCount: 0,
		},
		{
			Value:    `{"abc":["1","2"]}`,
			ErrCount: 0,
		},
	}

	for _, tc := range validCases {
		_, errors := validateJsonString(tc.Value, "json")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %q not to trigger a validation error.", tc.Value)
		}
	}
}

func TestValidateApiGatewayIntegrationType(t *testing.T) {
	type testCases struct {
		Value    string
		ErrCount int
	}

	invalidCases := []testCases{
		{
			Value:    "incorrect",
			ErrCount: 1,
		},
		{
			Value:    "aws_proxy",
			ErrCount: 1,
		},
	}

	for _, tc := range invalidCases {
		_, errors := validateApiGatewayIntegrationType(tc.Value, "types")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %q to trigger a validation error.", tc.Value)
		}
	}

	validCases := []testCases{
		{
			Value:    "MOCK",
			ErrCount: 0,
		},
		{
			Value:    "AWS_PROXY",
			ErrCount: 0,
		},
	}

	for _, tc := range validCases {
		_, errors := validateApiGatewayIntegrationType(tc.Value, "types")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected %q not to trigger a validation error.", tc.Value)
		}
	}
}

func TestValidateSQSQueueName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"valid02-name",
		"Valid-Name1",
		"_",
		"-",
		strings.Repeat("W", 80),
	}
	for _, v := range validNames {
		if errors := validateSQSQueueName(v, "name"); len(errors) > 0 {
			t.Fatalf("%q should be a valid SQS queue Name", v)
		}
	}

	invalidNames := []string{
		"Here is a name with: colon",
		"another * invalid name",
		"also $ invalid",
		"This . is also %% invalid@!)+(",
		"*",
		"",
		" ",
		".",
		strings.Repeat("W", 81), // length > 80
	}
	for _, v := range invalidNames {
		if errors := validateSQSQueueName(v, "name"); len(errors) == 0 {
			t.Fatalf("%q should be an invalid SQS queue Name", v)
		}
	}
}

func TestValidateSQSFifoQueueName(t *testing.T) {
	validNames := []string{
		"valid-name.fifo",
		"valid02-name.fifo",
		"Valid-Name1.fifo",
		"_.fifo",
		"a.fifo",
		"A.fifo",
		"9.fifo",
		"-.fifo",
		fmt.Sprintf("%s.fifo", strings.Repeat("W", 75)),
	}
	for _, v := range validNames {
		if errors := validateSQSFifoQueueName(v, "name"); len(errors) > 0 {
			t.Fatalf("%q should be a valid SQS FIFO queue Name: %v", v, errors)
		}
	}

	invalidNames := []string{
		"Here is a name with: colon",
		"another * invalid name",
		"also $ invalid",
		"This . is also %% invalid@!)+(",
		".fifo",
		"*",
		"",
		" ",
		".",
		strings.Repeat("W", 81), // length > 80
	}
	for _, v := range invalidNames {
		if errors := validateSQSFifoQueueName(v, "name"); len(errors) == 0 {
			t.Fatalf("%q should be an invalid SQS FIFO queue Name: %v", v, errors)
		}
	}
}

func TestValidateSNSSubscriptionProtocol(t *testing.T) {
	validProtocols := []string{
		"lambda",
		"sqs",
		"sqs",
		"application",
		"http",
		"https",
	}
	for _, v := range validProtocols {
		if _, errors := validateSNSSubscriptionProtocol(v, "protocol"); len(errors) > 0 {
			t.Fatalf("%q should be a valid SNS Subscription protocol: %v", v, errors)
		}
	}

	invalidProtocols := []string{
		"Email",
		"email",
		"Email-JSON",
		"email-json",
		"SMS",
		"sms",
	}
	for _, v := range invalidProtocols {
		if _, errors := validateSNSSubscriptionProtocol(v, "protocol"); len(errors) == 0 {
			t.Fatalf("%q should be an invalid SNS Subscription protocol: %v", v, errors)
		}
	}
}

func TestValidateSecurityRuleType(t *testing.T) {
	validTypes := []string{
		"ingress",
		"egress",
	}
	for _, v := range validTypes {
		if _, errors := validateSecurityRuleType(v, "type"); len(errors) > 0 {
			t.Fatalf("%q should be a valid Security Group Rule type: %v", v, errors)
		}
	}

	invalidTypes := []string{
		"foo",
		"ingresss",
	}
	for _, v := range invalidTypes {
		if _, errors := validateSecurityRuleType(v, "type"); len(errors) == 0 {
			t.Fatalf("%q should be an invalid Security Group Rule type: %v", v, errors)
		}
	}
}
