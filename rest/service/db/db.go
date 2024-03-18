package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DBClient struct {
	client *dynamodb.Client
	name   string
}

// AttributeValueMemberB - Binary, For example: "B":"dGhpcyB0ZXh0IGlzIGJhc2U2NC1lbmNvZGVk"
// AttributeValueMemberBOOL - Bool, For example: "BOOL": true
// AttributeValueMemberBS - Binary Set, For example: "BS": ["U3Vubnk=", "UmFpbnk=", "U25vd3k="]
// AttributeValueMemberL - List, For example: "L": [ {"S": "Cookies"} , {"S": "Coffee"}, {"N": "3.14159"}]
// AttributeValueMemberM - Map, For example: "M": {"Name": {"S": "Joe"}, "Age": {"N": "35"}}
// AttributeValueMemberN - Number, For example: "N": "123.45"
// AttributeValueMemberNS - Number Set, For example: "NS": ["42.2", "-19", "7.5", "3.14"]
// AttributeValueMemberNULL - NULL, For example: "NULL": true
// AttributeValueMemberS - String, For example: "S": "Hello"
// AttributeValueMemberSS - String Set, For example: "SS": ["Giraffe", "Hippo", "Zebra"]

type Keys string

const (
	PK Keys = "PK"
	SK Keys = "SK"
)

type Items map[Keys]interface{}

func (items Items) GetKey() map[string]types.AttributeValue {
	pk, err := attributevalue.Marshal(items["PK"])
	if err != nil {
		panic(err)
	}
	sk, err := attributevalue.Marshal(items["SK"])
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"PK": pk, "SK": sk}
}

// String returns the every properties of Test.
func (item Items) String() string {
	var strs string = ""
	for key, value := range item {
		strs += fmt.Sprintf("\n\t%s: %v", key, value)
	}
	strs += "\n"
	return strs
}

func init() {
	os.Setenv("MODE", "test")
	os.Setenv("AWS_ACCESS_KEY_ID", "id")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "token")
}

// setup client
func CreateClient(ctx context.Context, tableName string) *DBClient {
	mode := os.Getenv("MODE")
	var client *dynamodb.Client
	if mode == "test" {
		// Create a static credentials provider
		staticCreds := credentials.NewStaticCredentialsProvider("id", "token", "")
		// Configure the AWS client
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"), config.WithCredentialsProvider(staticCreds))
		if err != nil {
			// Handle error
			log.Fatal(err)
		}
		// Set the DynamoDB Local endpoint
		cfg.EndpointResolver = aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: "http://localhost:8000",
				}, nil
			},
		)
		// Create the DynamoDB client
		client = dynamodb.NewFromConfig(cfg)
	} else {
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
		if err != nil {
			panic(err)
		}
		// Create the DynamoDB client
		client = dynamodb.NewFromConfig(cfg)
	}

	return &(DBClient{client, tableName})
}

// create table
func (db *DBClient) CreateTable(ctx context.Context) (*types.TableDescription, error) {
	var tableDesc *types.TableDescription
	table, err := db.client.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName:   aws.String(db.name),
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		log.Printf("Couldn't create table %v. Here's why: %v\n", db.name, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(db.client)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(db.name)}, 30*time.Second)
		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
		tableDesc = table.TableDescription
	}
	return tableDesc, err
}

// create secondary index (global/local)
func (db *DBClient) CreateIndex(ctx context.Context) {
	db.client.UpdateTable(ctx, &dynamodb.UpdateTableInput{
		TableName: aws.String(db.name),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("GSI1PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("GSI1SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		GlobalSecondaryIndexUpdates: []types.GlobalSecondaryIndexUpdate{
			{
				Create: &types.CreateGlobalSecondaryIndexAction{
					IndexName: aws.String("GSI1_COUNT_PK"),
					KeySchema: []types.KeySchemaElement{
						{
							AttributeName: aws.String("GSI1PK"),
							KeyType:       types.KeyTypeHash,
						},
						{
							AttributeName: aws.String("GSI1SK"),
							KeyType:       types.KeyTypeRange,
						},
					},
				},
			},
		},
	})
}

// delete table
func (db *DBClient) DeleteTable(ctx context.Context) error {
	_, err := db.client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(db.name),
	})
	return err
}

// list table
func (db *DBClient) ListTable(ctx context.Context) {
	p := dynamodb.NewListTablesPaginator(db.client, nil, func(o *dynamodb.ListTablesPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	for p.HasMorePages() {
		out, err := p.NextPage(ctx)
		if err != nil {
			panic(err)
		}

		for _, tn := range out.TableNames {
			fmt.Println(tn)
		}
	}
}

// insert item into table
func (db *DBClient) InsertItem(ctx context.Context, item Items) error {
	data, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}
	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{TableName: aws.String(db.name), Item: data})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}

// get item from table
func (db *DBClient) GetItem(ctx context.Context, key map[string]types.AttributeValue) Items {
	var out Items
	response, err := db.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: key, TableName: aws.String(db.name),
	})
	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", key, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &out)
		if err != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}
	return out
}

func (db *DBClient) QueryRead(ctx context.Context, pk string) []Items {
	var out []Items
	keyEx := expression.Key("PK").Equal(expression.Value(pk))
	filtEx := expression.Name("Count").GreaterThan(expression.Value(5))
	projEx := expression.NamesList(expression.Name("Count"))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).WithFilter(filtEx).WithProjection(projEx).Build()
	if err != nil {
		log.Panicln("failed to build Expressions")
	}
	response, err := db.client.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(db.name),
		// expressions for read
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		log.Println(err)
		log.Panicln("failed to query")
	}
	err = attributevalue.UnmarshalListOfMaps(response.Items, &out)
	if err != nil {
		log.Panicln("failed to extract data")
	}
	return out
}

func (db *DBClient) QueryWrite(ctx context.Context, item Items) error {
	data, err := attributevalue.MarshalMap(item)
	fmt.Println(item)
	if err != nil {
		panic(err)
	}
	condExpr := expression.AttributeNotExists(expression.Name("PK")).And(expression.BeginsWith(expression.Name("SK"), "TEST#"))
	expr, err := expression.NewBuilder().WithCondition(condExpr).Build()
	if err != nil {
		log.Panicln("failed to build Expressions")
	}
	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:                 aws.String(db.name),
		Item:                      data,
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}
