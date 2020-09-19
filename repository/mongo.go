package repository

import (
	"context"

	"github.com/mikestefanello/otcscanner/config"
	"github.com/mikestefanello/otcscanner/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoOrderRepository struct {
	client           *mongo.Client
	config           config.MongoConfig
	filterCompleted  bson.M
	filterIncomplete bson.M
}

// NewMongoOrderRepository creates a new mongo DB repository for orders
func NewMongoOrderRepository(cfg config.MongoConfig) (OrderRepository, error) {
	repo := &mongoOrderRepository{
		config:           cfg,
		filterCompleted:  bson.M{"service": bson.M{"$ne": ""}},
		filterIncomplete: bson.M{"service": ""},
	}
	err := repo.connect()
	return repo, err
}

func (r *mongoOrderRepository) connect() error {
	ctx, cancel := r.contextWithTimeout()
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(r.config.URL))
	if err != nil {
		return err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	r.client = client

	return nil
}

func (r *mongoOrderRepository) contextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), r.config.Timeout)
}

func (r *mongoOrderRepository) getCollection() *mongo.Collection {
	return r.client.Database(r.config.DB).Collection("orders")
}

func (r *mongoOrderRepository) LoadByID(id string) (*models.Order, error) {
	ctx, cancel := r.contextWithTimeout()
	defer cancel()

	o := &models.Order{}
	err := r.getCollection().FindOne(ctx, bson.M{"packageId": id}).Decode(&o)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return o, nil
}

func (r *mongoOrderRepository) LoadAll() (*models.Orders, error) {
	return r.loadWithFilter(bson.M{})
}

func (r *mongoOrderRepository) LoadCompleted() (*models.Orders, error) {
	return r.loadWithFilter(r.filterCompleted)
}

func (r *mongoOrderRepository) LoadIncomplete() (*models.Orders, error) {
	return r.loadWithFilter(r.filterIncomplete)
}

func (r *mongoOrderRepository) DeleteAll() error {
	return r.deleteWithFilter(bson.M{})
}

func (r *mongoOrderRepository) DeleteCompleted() error {
	return r.deleteWithFilter(r.filterCompleted)
}

func (r *mongoOrderRepository) UpdateOne(order *models.Order) error {
	ctx, cancel := r.contextWithTimeout()
	defer cancel()

	filter := bson.M{"packageId": order.PackageID}
	update := bson.M{"$set": order}
	_, err := r.getCollection().UpdateOne(ctx, filter, update)

	return err
}

func (r *mongoOrderRepository) InsertMany(orders *models.Orders) error {
	ctx, cancel := r.contextWithTimeout()
	defer cancel()

	// TODO: Odd that this is needed?
	data := make([]interface{}, 0, len(*orders))
	for _, o := range *orders {
		data = append(data, o)
	}

	_, err := r.getCollection().InsertMany(ctx, data)

	return err
}

func (r *mongoOrderRepository) CountAll() (int64, error) {
	return r.countWithFilter(bson.M{})
}

func (r *mongoOrderRepository) CountCompleted() (int64, error) {
	return r.countWithFilter(r.filterCompleted)
}

func (r *mongoOrderRepository) CountIncomplete() (int64, error) {
	return r.countWithFilter(r.filterIncomplete)
}

func (r *mongoOrderRepository) loadWithFilter(filter bson.M) (*models.Orders, error) {
	ctx, cancel := r.contextWithTimeout()
	defer cancel()

	o := &models.Orders{}
	cursor, err := r.getCollection().Find(ctx, filter)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, o)

	return o, err
}

func (r *mongoOrderRepository) deleteWithFilter(filter bson.M) error {
	ctx, cancel := r.contextWithTimeout()
	defer cancel()

	_, err := r.getCollection().DeleteMany(ctx, filter)

	return err
}

func (r *mongoOrderRepository) countWithFilter(filter bson.M) (int64, error) {
	ctx, cancel := r.contextWithTimeout()
	defer cancel()

	return r.getCollection().CountDocuments(ctx, filter)
}
