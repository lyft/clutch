package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	dynamodbv1 "github.com/lyft/clutch/backend/api/aws/dynamodb/v1"
	awsv1 "github.com/lyft/clutch/backend/api/config/service/aws/v1"
)

// defaults for the dynamodb settings config as set by AWS
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/ServiceQuotas.html#default-limits-throughput-capacity-modes
const (
	AwsMaxRCU       = 40000
	AwsMaxWCU       = 40000
	SafeScaleFactor = 2.0
)

// constant to cover edge case where billing mode is not in table description
// and we cannot infer the billing mode from the provisioned throughput numbers
const TableBillingUnspecified = "UNSPECIFIED"

// get or set defaults for dynamodb scaling
func getScalingLimits(cfg *awsv1.Config) *awsv1.ScalingLimits {
	if cfg.GetDynamodbConfig() == nil && cfg.DynamodbConfig.GetScalingLimits() == nil {
		ds := &awsv1.ScalingLimits{
			MaxReadCapacityUnits:  AwsMaxRCU,
			MaxWriteCapacityUnits: AwsMaxWCU,
			MaxScaleFactor:        float32(SafeScaleFactor),
			EnableOverride:        false,
		}
		return ds
	}
	return cfg.DynamodbConfig.ScalingLimits
}

func (c *client) DescribeTable(ctx context.Context, account, region, tableName string) (*dynamodbv1.Table, error) {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return nil, err
	}

	result, err := getTable(ctx, cl, tableName)

	if err != nil {
		c.log.Warn("unable to find table in region: "+region, zap.Error(err))
		return nil, err
	}

	ret := newProtoForTable(result.Table, account, region)

	return ret, nil
}

func getTable(ctx context.Context, client *regionalClient, tableName string) (*dynamodb.DescribeTableOutput, error) {
	input := &dynamodb.DescribeTableInput{TableName: aws.String(tableName)}
	return client.dynamodb.DescribeTable(ctx, input)
}

// generator for the list of GSI protos
func getGlobalSecondaryIndexes(indexes []types.GlobalSecondaryIndexDescription) []*dynamodbv1.GlobalSecondaryIndex {
	gsis := make([]*dynamodbv1.GlobalSecondaryIndex, len(indexes))
	for idx, i := range indexes {
		gsis[idx] = newProtoForGlobalSecondaryIndex(i)
	}
	return gsis
}

func newProtoForKeySchemas(inputSchema []types.KeySchemaElement) []*dynamodbv1.KeySchema {
	schemaCollection := make([]*dynamodbv1.KeySchema, len(inputSchema))
	for idx, schema := range inputSchema {
		schemaCollection[idx] = &dynamodbv1.KeySchema{
			AttributeName: *schema.AttributeName,
			Type:          newProtoForKeySchemaType(schema.KeyType),
		}
	}
	return schemaCollection
}

func newProtoForKeySchemaType(keyType types.KeyType) dynamodbv1.KeySchema_Type {
	value, ok := dynamodbv1.KeySchema_Type_value[string(keyType)]
	if !ok {
		return dynamodbv1.KeySchema_UNKNOWN
	}
	return dynamodbv1.KeySchema_Type(value)
}
func newProtoForAttributeDefinitions(inputAttributes []types.AttributeDefinition) []*dynamodbv1.AttributeDefinition {
	attributeDefs := make([]*dynamodbv1.AttributeDefinition, len(inputAttributes))
	for idx, attribute := range inputAttributes {
		attributeDefs[idx] = &dynamodbv1.AttributeDefinition{
			AttributeName: *attribute.AttributeName,
			AttributeType: string(attribute.AttributeType),
		}
	}
	return attributeDefs
}

// retrieve one GSI from table description
func getGlobalSecondaryIndex(indexes []types.GlobalSecondaryIndexDescription, targetIndexName string) (*types.GlobalSecondaryIndexDescription, error) {
	for _, i := range indexes {
		if *i.IndexName == targetIndexName {
			return &i, nil
		}
	}
	return nil, status.Error(codes.NotFound, "Global secondary index not found.")
}

func newProtoForTable(t *types.TableDescription, account, region string) *dynamodbv1.Table {
	currentCapacity := &dynamodbv1.Throughput{
		WriteCapacityUnits: aws.ToInt64(t.ProvisionedThroughput.WriteCapacityUnits),
		ReadCapacityUnits:  aws.ToInt64(t.ProvisionedThroughput.ReadCapacityUnits),
	}

	globalSecondaryIndexes := getGlobalSecondaryIndexes(t.GlobalSecondaryIndexes)

	tableStatus := newProtoForTableStatus(t.TableStatus)

	billingMode := newProtoForBillingMode(t)
	keySchemas := newProtoForKeySchemas(t.KeySchema)
	attributeDefinitions := newProtoForAttributeDefinitions(t.AttributeDefinitions)

	ret := &dynamodbv1.Table{
		Name:                   aws.ToString(t.TableName),
		Account:                account,
		Region:                 region,
		GlobalSecondaryIndexes: globalSecondaryIndexes,
		ProvisionedThroughput:  currentCapacity,
		Status:                 tableStatus,
		BillingMode:            billingMode,
		KeySchemas:             keySchemas,
		AttributeDefinitions:   attributeDefinitions,
	}

	return ret
}

func newProtoForContinuousBackups(t *types.ContinuousBackupsDescription, account string, region string) *dynamodbv1.ContinuousBackups {
	backups := &dynamodbv1.ContinuousBackups{
		ContinuousBackupsStatus: newProtoForBackupStatus(string(t.ContinuousBackupsStatus)),
	}

	// if *t.PointInTimeRecoveryDescription is nil
	// then we set PointInTimeRecoveryStatus to UNSPECIFIED
	if t.PointInTimeRecoveryDescription == nil || t.PointInTimeRecoveryDescription.PointInTimeRecoveryStatus == "" {
		backups.PointInTimeRecoveryStatus = dynamodbv1.ContinuousBackups_UNSPECIFIED
	} else {
		backups.PointInTimeRecoveryStatus = newProtoForBackupStatus(string(t.PointInTimeRecoveryDescription.PointInTimeRecoveryStatus))
	}

	// if *t.PointInTimeRecoveryDescription is nil or if EarliestRestorableDateTime is nil
	// then we set EarliestRestorableDateTime to nil
	if t.PointInTimeRecoveryDescription != nil && t.PointInTimeRecoveryDescription.EarliestRestorableDateTime != nil {
		backups.EarliestRestorableDateTime = timestamppb.New(*t.PointInTimeRecoveryDescription.EarliestRestorableDateTime)
	} else {
		backups.EarliestRestorableDateTime = nil
	}

	// if *t.PointInTimeRecoveryDescription is nil or if LatestRestorableDateTime is nil
	// then we set LatestRestorableDateTime to nil
	if t.PointInTimeRecoveryDescription != nil && t.PointInTimeRecoveryDescription.LatestRestorableDateTime != nil {
		backups.LatestRestorableDateTime = timestamppb.New(*t.PointInTimeRecoveryDescription.LatestRestorableDateTime)
	} else {
		backups.LatestRestorableDateTime = nil
	}

	return backups
}

func newProtoForTableStatus(s types.TableStatus) dynamodbv1.Table_Status {
	value, ok := dynamodbv1.Table_Status_value[string(s)]
	if !ok {
		return dynamodbv1.Table_UNKNOWN
	}
	return dynamodbv1.Table_Status(value)
}

func newProtoForIndexStatus(s types.IndexStatus) dynamodbv1.GlobalSecondaryIndex_Status {
	value, ok := dynamodbv1.GlobalSecondaryIndex_Status_value[string(s)]
	if !ok {
		return dynamodbv1.GlobalSecondaryIndex_UNKNOWN
	}
	return dynamodbv1.GlobalSecondaryIndex_Status(value)
}

func newProtoForBackupStatus(s string) dynamodbv1.ContinuousBackups_Status {
	value, ok := dynamodbv1.ContinuousBackups_Status_value[s]
	if !ok {
		return dynamodbv1.ContinuousBackups_UNKNOWN
	}
	return dynamodbv1.ContinuousBackups_Status(value)
}

// manually retrieve the billing mode by inferring it from throughput
// to cover cases where AWS does not return the mode in the table description
// if a table is PROVISIONED, it will have at least 1 RCU/WCU provisioned
// if a table is PAY_PER_REQUEST (on demand), it will have 0 RCU/WCU provisioned
func getBillingMode(t *types.ProvisionedThroughputDescription) types.BillingMode {
	if (*t.ReadCapacityUnits > 0) || (*t.WriteCapacityUnits > 0) {
		return types.BillingModeProvisioned
	}
	if (*t.ReadCapacityUnits == 0) || (*t.WriteCapacityUnits == 0) {
		return types.BillingModePayPerRequest
	}

	return TableBillingUnspecified // unable to infer what billing mode it is
}

func newProtoForBillingMode(t *types.TableDescription) dynamodbv1.Table_BillingMode {
	var billingMode types.BillingMode
	if t.BillingModeSummary == nil {
		billingMode = getBillingMode(t.ProvisionedThroughput)
	} else {
		billingMode = t.BillingModeSummary.BillingMode
	}

	value, ok := dynamodbv1.Table_BillingMode_value[string(billingMode)]
	if !ok {
		return dynamodbv1.Table_BILLING_UNKNOWN
	}

	return dynamodbv1.Table_BillingMode(value)
}

func newProtoForGlobalSecondaryIndex(index types.GlobalSecondaryIndexDescription) *dynamodbv1.GlobalSecondaryIndex {
	currentCapacity := &dynamodbv1.Throughput{
		ReadCapacityUnits:  aws.ToInt64(index.ProvisionedThroughput.ReadCapacityUnits),
		WriteCapacityUnits: aws.ToInt64(index.ProvisionedThroughput.WriteCapacityUnits),
	}

	indexStatus := newProtoForIndexStatus(index.IndexStatus)

	return &dynamodbv1.GlobalSecondaryIndex{
		Name:                  aws.ToString(index.IndexName),
		ProvisionedThroughput: currentCapacity,
		Status:                indexStatus,
	}
}

func isValidIncrease(client *regionalClient, current *types.ProvisionedThroughputDescription, target types.ProvisionedThroughput, ignore_maximums bool) error {
	// check for targets that are lower than current (can't scale down)
	if *current.ReadCapacityUnits > *target.ReadCapacityUnits {
		return status.Errorf(codes.FailedPrecondition, "Target read capacity [%d] is lower than current capacity [%d]", *target.ReadCapacityUnits, *current.ReadCapacityUnits)
	}
	if *current.WriteCapacityUnits > *target.WriteCapacityUnits {
		return status.Errorf(codes.FailedPrecondition, "Target write capacity [%d] is lower than current capacity [%d]", *target.WriteCapacityUnits, *current.WriteCapacityUnits)
	}

	// override not enabled in config or override not set to true in args
	if !client.dynamodbCfg.ScalingLimits.EnableOverride || !ignore_maximums {
		// check for targets that exceed max limits
		if *target.ReadCapacityUnits > client.dynamodbCfg.ScalingLimits.MaxReadCapacityUnits &&
			*current.ReadCapacityUnits < client.dynamodbCfg.ScalingLimits.MaxReadCapacityUnits { // don't apply this check to tables with RCU already above max
			return status.Errorf(codes.FailedPrecondition, "Target read capacity exceeds maximum allowed limits [%d]", client.dynamodbCfg.ScalingLimits.MaxReadCapacityUnits)
		}
		if *target.WriteCapacityUnits > client.dynamodbCfg.ScalingLimits.MaxWriteCapacityUnits &&
			*current.WriteCapacityUnits < client.dynamodbCfg.ScalingLimits.MaxWriteCapacityUnits { // don't apply this check to tables with WCU already above max
			return status.Errorf(codes.FailedPrecondition, "Target write capacity exceeds maximum allowed limits [%d]", client.dynamodbCfg.ScalingLimits.MaxWriteCapacityUnits)
		}

		// check for increases that exceed max increase scale
		if (float32(*target.ReadCapacityUnits) / float32(*current.ReadCapacityUnits)) > client.dynamodbCfg.ScalingLimits.MaxScaleFactor {
			return status.Errorf(codes.FailedPrecondition, "Target read capacity exceeds the scale limit of [%.1f]x current capacity", client.dynamodbCfg.ScalingLimits.MaxScaleFactor)
		}
		if (float32(*target.WriteCapacityUnits) / float32(*current.WriteCapacityUnits)) > client.dynamodbCfg.ScalingLimits.MaxScaleFactor {
			return status.Errorf(codes.FailedPrecondition, "Target write capacity exceeds the scale limit of [%.1f]x current capacity", client.dynamodbCfg.ScalingLimits.MaxScaleFactor)
		}
	}
	return nil
}

// Note: in some cases, table descriptions may not contain the billing mode summary
// even though they're provisioned, so we need to check for throughput settings
// as well to determine if a table capacity is provisioned
func isProvisioned(t *dynamodb.DescribeTableOutput) bool {
	if t.Table.BillingModeSummary == nil {
		billingMode := getBillingMode(t.Table.ProvisionedThroughput)
		return billingMode == types.BillingModeProvisioned
	}
	return t.Table.BillingModeSummary.BillingMode == types.BillingModeProvisioned
}

func (c *client) UpdateCapacity(ctx context.Context, account, region, tableName string, targetTableCapacity *dynamodbv1.Throughput, indexUpdates []*dynamodbv1.IndexUpdateAction, ignore_maximums bool) (*dynamodbv1.Table, error) {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return nil, err
	}

	currentTable, err := getTable(ctx, cl, tableName)
	if err != nil {
		c.log.Error("unable to find table", zap.Error(err))
		return nil, err
	}

	if !(isProvisioned(currentTable)) {
		return nil, status.Error(codes.FailedPrecondition, "Table billing mode is not set to PROVISIONED, cannot scale capacities.")
	}

	// initialize the input we'll pass to AWS
	input := &dynamodb.UpdateTableInput{
		TableName: aws.String(tableName),
	}

	// add target table capacity to input if valid
	if targetTableCapacity != nil { // received new vals for table
		t := types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(targetTableCapacity.ReadCapacityUnits),
			WriteCapacityUnits: aws.Int64(targetTableCapacity.WriteCapacityUnits),
		}

		err = isValidIncrease(cl, currentTable.Table.ProvisionedThroughput, t, ignore_maximums)
		if err != nil {
			return nil, err
		}

		input.ProvisionedThroughput = &t
	}

	// add target index capacities to input if valid
	if len(indexUpdates) > 0 { // received at least one index update
		updates, err := generateIndexUpdates(cl, currentTable, indexUpdates, ignore_maximums)
		if err != nil {
			return nil, err
		}
		input.GlobalSecondaryIndexUpdates = append(input.GlobalSecondaryIndexUpdates, updates...)
	}

	result, err := cl.dynamodb.UpdateTable(ctx, input)
	if err != nil {
		return nil, err
	}

	ret := newProtoForTable(result.TableDescription, account, region)

	return ret, nil
}

// given a list of ddbv1.updates, generate the AWS types for updates
// TO DO: this function currently terminates upon encountering one invalid index update,
// preserving the AWS SDK design of "all-or-nothing" UpdateTable functionality.
// confirm if that is still the desired behavior. Alternatively, we can selectively make updates.
func generateIndexUpdates(cl *regionalClient, t *dynamodb.DescribeTableOutput, indexUpdates []*dynamodbv1.IndexUpdateAction, ignore_maximums bool) ([]types.GlobalSecondaryIndexUpdate, error) {
	currentIndexes := t.Table.GlobalSecondaryIndexes

	var updates []types.GlobalSecondaryIndexUpdate

	for _, i := range indexUpdates {
		currentIndex, err := getGlobalSecondaryIndex(currentIndexes, i.Name)
		if err != nil {
			return nil, err
		}

		targetIndexCapacity := types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(i.IndexThroughput.ReadCapacityUnits),
			WriteCapacityUnits: aws.Int64(i.IndexThroughput.WriteCapacityUnits),
		}

		err = isValidIncrease(cl, currentIndex.ProvisionedThroughput, targetIndexCapacity, ignore_maximums)
		if err != nil {
			return nil, err
		}

		newIndexUpdate := types.GlobalSecondaryIndexUpdate{
			Update: &types.UpdateGlobalSecondaryIndexAction{
				IndexName: aws.String(i.Name),
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(i.IndexThroughput.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(i.IndexThroughput.WriteCapacityUnits),
				},
			},
		}

		updates = append(updates, newIndexUpdate)
	}

	return updates, nil
}

func (c *client) BatchGetItem(ctx context.Context, account string, region string, input *dynamodb.BatchGetItemInput) (*dynamodb.BatchGetItemOutput, error) {
	client, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return nil, err
	}
	return client.dynamodb.BatchGetItem(ctx, input)
}

func (c *client) DescribeContinuousBackups(ctx context.Context, account string, region string, tableName string) (*dynamodbv1.ContinuousBackups, error) {
	cl, err := c.getAccountRegionClient(account, region)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.DescribeContinuousBackupsInput{TableName: aws.String(tableName)}

	res, err := cl.dynamodb.DescribeContinuousBackups(ctx, input)

	if err != nil {
		c.log.Warn("unable to get backups description for table: "+tableName, zap.Error(err))
		return nil, err
	}

	ret := newProtoForContinuousBackups(res.ContinuousBackupsDescription, account, region)

	return ret, nil
}
