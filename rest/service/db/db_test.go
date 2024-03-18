package db_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"restapi/service/db"
	"restapi/utils/logger"
	"testing"
	"time"

	re "restapi/utils/error"
)

func TestCreateTable(t *testing.T) {
	cases := []struct {
		tableName string
	}{
		{
			tableName: "test-table",
		},
	}
	context := context.TODO()

	for _, c := range cases {
		clt := db.CreateClient(context, c.tableName)
		clt.DeleteTable(context)
		des, err := clt.CreateTable(context)

		if !(err == nil && *(des.TableName) == c.tableName) {
			t.Errorf("CreateTable(context, %v) = %v, expected: %v", c.tableName, *(des.TableName), c.tableName)
		}
	}
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

type Test struct {
	TestID   string                 `json:"TestID"`
	Time     int                    `json:"Time"`
	Count    int                    `json:"Count"`
	IsActive bool                   `json:"IsActive"`
	List     []string               `json:"List"`
	Map      map[string]interface{} `json:"Map"`
}

func (t Test) GetItems() db.Items {
	return db.Items{
		"PK":       "TEST#" + t.TestID,
		"SK":       "TEST#" + fmt.Sprintf("%v", t.Count),
		"TestID":   t.TestID,
		"Time":     t.Time,
		"Count":    t.Count,
		"IsActive": t.IsActive,
		"List":     t.List,
		"Map":      t.Map,
	}
}
func ExtractTest(data db.Items) (Test, error) {
	var result Test
	delete(data, "PK")
	delete(data, "SK")
	jsonData, err := json.Marshal(data)

	if err != nil {
		return result, re.RequestError{
			Code:  500,
			Msg:   "from ExtractTest func when marshaling",
			Cause: err,
		}
	}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return result, re.RequestError{
			Code:  500,
			Msg:   "from ExtractTest func when unmarshaling",
			Cause: err,
		}
	}
	return result, nil
}

func (t Test) GetItemsWithGSI() db.Items {
	return db.Items{
		"PK":       "TEST#" + t.TestID,
		"SK":       "TEST#" + fmt.Sprintf("%v", t.Count),
		"TestID":   t.TestID,
		"Count":    t.Count,
		"IsActive": t.IsActive,
		"List":     t.List,
		"Map":      t.Map,
	}
}

func TestInsertItem(t *testing.T) {
	context := context.TODO()
	cases := []struct {
		inputItem Test
	}{
		{
			inputItem: Test{
				TestID:   "123",
				Time:     int(time.Now().UnixNano()),
				Count:    rand.Intn(10),
				IsActive: rand.Intn(2) == 0,
				List:     []string{"Hello", "World"},
				Map: map[string]interface{}{
					"hello": 1, "world": 2,
				},
			},
		},
		{
			inputItem: Test{
				TestID:   "123",
				Time:     int(time.Now().UnixNano()),
				Count:    rand.Intn(10),
				IsActive: rand.Intn(2) == 0,
				List:     []string{"Hello", "World"},
				Map: map[string]interface{}{
					"hello": 1, "world": 2,
				},
			},
		},
	}
	for _, c := range cases {
		items := c.inputItem.GetItems()
		clt := db.CreateClient(context, "test-table")
		err := clt.InsertItem(context, items)
		if err != nil {
			log.Fatalln("failed to insert")
		}
		output := clt.GetItem(context, items.GetKey())
		test, _ := ExtractTest(output)
		logger.Log.Info(test)
	}
}

func TestQuery(t *testing.T) {
	ctx := context.TODO()
	cases := []struct {
		inputPK string
	}{
		{
			inputPK: "TEST#123",
		},
	}
	for _, c := range cases {
		clt := db.CreateClient(ctx, "test-table")
		output := clt.QueryRead(ctx, c.inputPK)
		for _, out := range output {
			logger.Log.Info(out)
		}
	}
}
