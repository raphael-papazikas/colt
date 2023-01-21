package colt

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Hooks[T Document] struct {
	deleteOne  []func(filter *bson.M) error
	updateOne  []func(filter *bson.M, model *T) error
	updateMany []func(filter *bson.M, update *bson.M) error
	findOne    []func(filter *bson.M) error
	find       []func(filter *bson.M) error
	count      []func(filter *bson.M) error
}

func (hooks *Hooks[T]) OnDeleteOne(callback func(filter *bson.M) error) {
	hooks.deleteOne = append(hooks.deleteOne, callback)
}

func (hooks *Hooks[T]) OnUpdateOne(callback func(filter *bson.M, model *T) error) {
	hooks.updateOne = append(hooks.updateOne, callback)
}

func (hooks *Hooks[T]) OnUpdateMany(callback func(filter *bson.M, update *bson.M) error) {
	hooks.updateMany = append(hooks.updateMany, callback)
}

func (hooks *Hooks[T]) OnFindOne(callback func(filter *bson.M) error) {
	hooks.findOne = append(hooks.findOne, callback)
}

func (hooks *Hooks[T]) OnFind(callback func(filter *bson.M) error) {
	hooks.find = append(hooks.find, callback)
}

func (hooks *Hooks[T]) OnCount(callback func(filter *bson.M) error) {
	hooks.count = append(hooks.count, callback)
}

type Collection[T Document] struct {
	collection  *mongo.Collection
	BeforeHooks Hooks[T]
}

func (repo *Collection[T]) Insert(model T) (T, error) {
	if model.GetID() == "" {
		model.SetID(repo.NewId().Hex())
	}

	if hook, ok := any(model).(BeforeInsertHook); ok {
		if err := hook.BeforeInsert(); err != nil {
			return model, err
		}
	}

	res, err := repo.collection.InsertOne(DefaultContext(), model)
	model.SetID(res.InsertedID.(string))
	return model, err
}

func (repo *Collection[T]) UpdateById(id string, model T) error {
	return repo.UpdateOne(bson.M{"_id": id}, model)
}

func (repo *Collection[T]) UpdateOne(filter bson.M, model T) error {
	for _, hook := range repo.BeforeHooks.updateOne {
		if err := hook(&filter, &model); err != nil {
			return err
		}
	}

	if hook, ok := any(model).(BeforeUpdateHook); ok {
		if err := hook.BeforeUpdate(); err != nil {
			return err
		}
	}

	_, err := repo.collection.UpdateOne(DefaultContext(), filter, bson.M{"$set": model})
	return err
}

func (repo *Collection[T]) UpdateMany(filter bson.M, doc bson.M) error {
	for _, hook := range repo.BeforeHooks.updateMany {
		if err := hook(&filter, &doc); err != nil {
			return err
		}
	}

	_, err := repo.collection.UpdateMany(DefaultContext(), filter, doc)
	return err
}

func (repo *Collection[T]) DeleteById(id string) error {
	return repo.DeleteOne(bson.M{"_id": id})
}

func (repo *Collection[T]) DeleteOne(filter bson.M) error {
	for _, hook := range repo.BeforeHooks.deleteOne {
		if err := hook(&filter); err != nil {
			return err
		}
	}

	res, err := repo.collection.DeleteOne(DefaultContext(), filter)

	if err != nil {
		return err
	}

	if res.DeletedCount < 1 {
		return errors.New("could not delete")
	}

	return nil
}

func (repo *Collection[T]) FindById(id string) (T, error) {
	return repo.FindOne(bson.M{"_id": id})
}

func (repo *Collection[T]) FindOne(filter bson.M) (T, error) {
	var target T
	for _, hook := range repo.BeforeHooks.findOne {
		if err := hook(&filter); err != nil {
			return target, err
		}
	}

	err := repo.collection.FindOne(DefaultContext(), filter).Decode(&target)

	return target, err
}

func (repo *Collection[T]) Find(filter bson.M, opts ...*options.FindOptions) ([]T, error) {
	var result []T

	for _, hook := range repo.BeforeHooks.find {
		if err := hook(&filter); err != nil {
			return result, err
		}
	}

	csr, err := repo.collection.Find(DefaultContext(), filter, opts...)
	if err = csr.All(DefaultContext(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *Collection[T]) CountDocuments(filter bson.M) (int64, error) {
	for _, hook := range repo.BeforeHooks.count {
		if err := hook(&filter); err != nil {
			return 0, err
		}
	}

	count, err := repo.collection.CountDocuments(DefaultContext(), filter)
	return count, err
}

func (repo *Collection[T]) NewId() primitive.ObjectID {
	return primitive.NewObjectID()
}
