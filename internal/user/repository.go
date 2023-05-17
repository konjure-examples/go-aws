package user

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type userEntity struct {
	id    string
	name  string
	email string
}

type Repository struct {
	tableName *string
	dbClient  *dynamodb.Client
}

func NewRepository(cfg aws.Config, tableName string) *Repository {
	dbClient := dynamodb.NewFromConfig(cfg)

	return &Repository{tableName: aws.String(tableName), dbClient: dbClient}
}

func (r *Repository) create(ctx context.Context, id, userName, email string) error {
	input := &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"PK":     &types.AttributeValueMemberS{Value: id},
			"SK":     &types.AttributeValueMemberS{Value: "USER#"},
			"name":   &types.AttributeValueMemberS{Value: userName},
			"email":  &types.AttributeValueMemberS{Value: email},
			"GSI1SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("EMAIL#%s", email)},
		},
	}

	_, err := r.dbClient.PutItem(ctx, input)

	return err

}

func (r *Repository) getById(ctx context.Context, id string) (*userEntity, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: id},
			"SK": &types.AttributeValueMemberS{Value: "USER#"},
		},
		TableName: r.tableName,
	}

	output, err := r.dbClient.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if output.Item == nil {
		return nil, fmt.Errorf("user with id %s not found", id)
	}

	entity := &userEntity{
		id:    output.Item["PK"].(*types.AttributeValueMemberS).Value,
		name:  output.Item["name"].(*types.AttributeValueMemberS).Value,
		email: output.Item["email"].(*types.AttributeValueMemberS).Value,
	}

	return entity, nil

}

func (r *Repository) getByEmail(ctx context.Context, email string) (*userEntity, error) {
	input := &dynamodb.QueryInput{
		TableName:                r.tableName,
		ConditionalOperator:      "",
		ConsistentRead:           nil,
		ExclusiveStartKey:        nil,
		ExpressionAttributeNames: map[string]string{"pk": "SK", "sk": "GSI1SK"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "USER#"},
			":sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("EMAIL#%s", email)},
		},
		IndexName:              aws.String("GSI1SK"),
		KeyConditionExpression: aws.String("pk = :pk AND sk = :sk"),
		Limit:                  aws.Int32(1),
	}

	output, err := r.dbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(output.Items) == 0 {
		return nil, fmt.Errorf("item not found")
	}

	item := output.Items[0]

	entity := &userEntity{
		id:    item["PK"].(*types.AttributeValueMemberS).Value,
		name:  item["name"].(*types.AttributeValueMemberS).Value,
		email: item["email"].(*types.AttributeValueMemberS).Value,
	}

	return entity, nil
}
