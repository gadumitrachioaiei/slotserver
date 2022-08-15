package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/gadumitrachioaiei/slotserver/slot"
)

type Users struct {
	c            *dynamodb.Client
	defaultChips int
}

// NewUsers returns an initialized users service.
func NewUsers(c *dynamodb.Client, defaultChips int) (*Users, error) {
	users := &Users{c: c, defaultChips: defaultChips}
	return users, nil
}

func (u *Users) Update(ctx context.Context, id string, amount int) (slot.User, error) {
	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					Item: map[string]types.AttributeValue{
						"ID":    &types.AttributeValueMemberS{Value: id},
						"Chips": &types.AttributeValueMemberN{Value: strconv.Itoa(u.defaultChips)},
					},
					TableName:           aws.String("Users"),
					ConditionExpression: aws.String("attribute_not_exists(ID)"),
				},
			},
		},
	}
	_, err := u.c.TransactWriteItems(ctx, input)
	if err != nil && !isTransactionConditionalCheckFailed(err) {
		return slot.User{}, fmt.Errorf("cannot create user: %v", err)
	}
	input = &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Update: &types.Update{
					Key: map[string]types.AttributeValue{
						"ID": &types.AttributeValueMemberS{Value: id},
					},
					TableName:           aws.String("Users"),
					UpdateExpression:    aws.String("SET Chips = Chips + :amount"),
					ConditionExpression: aws.String("Chips >= :amount"),
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":amount": &types.AttributeValueMemberN{Value: strconv.Itoa(amount)},
					},
				},
			},
		},
	}
	_, err = u.c.TransactWriteItems(ctx, input)
	if err != nil {
		if isTransactionConditionalCheckFailed(err) {
			return slot.User{}, errors.New("not enough chips")
		}
		return slot.User{}, fmt.Errorf("cannot update user: %v", err)
	}
	//TODO: how do I read the updated users object ?
	return slot.User{ID: id, Chips: 0}, nil
}

// createTable creates the user table.
//
// I don't know why but it doesn't work, it hangs.
func (u *Users) createTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := u.c.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Chips"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("ID"),
			KeyType:       types.KeyTypeHash,
		}},
		TableName: aws.String("Users"),
	})
	if err != nil {
		return fmt.Errorf("cannot create users table: %v", err)
	}
	return err
}

func isTransactionConditionalCheckFailed(err error) bool {
	var tce *types.TransactionCanceledException
	if !errors.As(err, &tce) {
		return false
	}
	if len(tce.CancellationReasons) != 1 || tce.CancellationReasons[0].Code == nil ||
		*tce.CancellationReasons[0].Code != "ConditionalCheckFailed" {
		return false
	}
	return true
}
