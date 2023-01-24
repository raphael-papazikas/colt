package colt

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"testing"
)

func TestCollection_OnFind(t *testing.T) {
	mockDb := MockSetup()
	collection := GetCollection[*testdoc](mockDb, "testdocs")

	var firstHookCalled = false
	var secondHookCalled = false

	doc := testdoc{Title: fmt.Sprint(rand.Int())}
	doc.SetID("secondId")
	collection.Insert(&doc)

	collection.BeforeHooks.OnFind(func(filter *bson.M) error {
		assert.Equal(t, (*filter)["_id"], "firstId")
		(*filter)["_id"] = "secondId"
		firstHookCalled = true
		return nil
	})

	collection.BeforeHooks.OnFind(func(filter *bson.M) error {
		assert.Equal(t, (*filter)["_id"], "secondId")
		secondHookCalled = true
		return nil
	})

	result, err := collection.Find(bson.M{"_id": "firstId"})
	assert.Equal(t, firstHookCalled, true)
	assert.Equal(t, secondHookCalled, true)
	assert.Nil(t, err)
	assert.Equal(t, result[0].ID, doc.ID)
	assert.Equal(t, result[0].Title, doc.Title)
}

func TestCollection_OnFind_Error(t *testing.T) {
	mockDb := MockSetup()
	collection := GetCollection[*testdoc](mockDb, "testdocs")

	doc := testdoc{Title: fmt.Sprint(rand.Int())}
	collection.Insert(&doc)

	collection.BeforeHooks.OnFind(func(filter *bson.M) error {
		return errors.New("Hook failing")
	})

	result, err := collection.Find(bson.M{})
	assert.NotNil(t, err)
	assert.Equal(t, len(result), 0)
}

func TestCollection_OnFindOne(t *testing.T) {
	mockDb := MockSetup()
	collection := GetCollection[*testdoc](mockDb, "testdocs")

	doc := testdoc{Title: fmt.Sprint(rand.Int())}
	doc.SetID("secondId")
	collection.Insert(&doc)

	collection.BeforeHooks.OnFindOne(func(filter *bson.M) error {
		assert.Equal(t, (*filter)["_id"], "firstId")
		(*filter)["_id"] = "secondId"
		return nil
	})

	result, err := collection.FindOne(bson.M{"_id": "firstId"})
	assert.Nil(t, err)
	assert.Equal(t, result.ID, doc.ID)
	assert.Equal(t, result.Title, doc.Title)
}

func TestCollection_OnFindOne_Error(t *testing.T) {
	mockDb := MockSetup()
	collection := GetCollection[*testdoc](mockDb, "testdocs")

	doc := testdoc{Title: fmt.Sprint(rand.Int())}
	collection.Insert(&doc)

	collection.BeforeHooks.OnFindOne(func(filter *bson.M) error {
		return errors.New("Hook failing")
	})

	result, err := collection.FindOne(bson.M{"_id": doc.ID})
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestCollection_OnCount(t *testing.T) {
	mockDb := MockSetup()
	collection := GetCollection[*testdoc](mockDb, "testdocs")

	doc := testdoc{Title: fmt.Sprint(rand.Int())}
	doc.SetID("secondId")
	collection.Insert(&doc)

	collection.BeforeHooks.OnCount(func(filter *bson.M) error {
		assert.Equal(t, (*filter)["_id"], "firstId")
		(*filter)["_id"] = "secondId"
		return nil
	})

	result, err := collection.CountDocuments(bson.M{"_id": "firstId"})
	assert.Nil(t, err)
	assert.Equal(t, result, int64(1))
}

func TestCollection_OnCount_Error(t *testing.T) {
	mockDb := MockSetup()
	collection := GetCollection[*testdoc](mockDb, "testdocs")

	doc := testdoc{Title: fmt.Sprint(rand.Int())}
	collection.Insert(&doc)

	collection.BeforeHooks.OnCount(func(filter *bson.M) error {
		return errors.New("Hook failing")
	})

	result, err := collection.CountDocuments(bson.M{"_id": doc.ID})
	assert.NotNil(t, err)
	assert.Equal(t, result, int64(0))
}

func TestCollection_OnUpdateOne(t *testing.T)       {}
func TestCollection_OnUpdateOne_Error(t *testing.T) {}

func TestCollection_OnUpdateMany(t *testing.T)       {}
func TestCollection_OnUpdateMany_Error(t *testing.T) {}

func TestCollection_OnDeleteOne(t *testing.T)       {}
func TestCollection_OnDeleteOne_Error(t *testing.T) {}
