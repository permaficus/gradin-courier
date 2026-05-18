package courier

import (
	"context"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(database *mongo.Database) *Repository {
	return &Repository{collection: database.Collection(CollectionName)}
}

func (repository *Repository) EnsureIndexes(ctx context.Context) error {
	models := []mongo.IndexModel{
		{Keys: bson.D{{Key: "name", Value: 1}}, Options: options.Index().SetName("couriers_name_idx")},
		{Keys: bson.D{{Key: "level", Value: 1}}, Options: options.Index().SetName("couriers_level_idx")},
		{Keys: bson.D{{Key: "registered_at", Value: 1}}, Options: options.Index().SetName("couriers_registered_at_idx")},
		{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetName("couriers_email_unique_sparse_idx").SetUnique(true).SetSparse(true)},
	}
	_, err := repository.collection.Indexes().CreateMany(ctx, models)
	return err
}

func (repository *Repository) Create(ctx context.Context, courier Courier) (Courier, error) {
	courier.ID = bson.NewObjectID()
	_, err := repository.collection.InsertOne(ctx, courier)
	return courier, err
}

func (repository *Repository) List(ctx context.Context, params QueryParams) (ListResult, error) {
	filter := repository.buildFilter(params)
	total, err := repository.collection.CountDocuments(ctx, filter)
	if err != nil {
		return ListResult{}, err
	}

	sort := bson.D{{Key: "name", Value: 1}}
	switch params.Sort {
	case "registered_at":
		sort = bson.D{{Key: "registered_at", Value: 1}}
	case "-registered_at":
		sort = bson.D{{Key: "registered_at", Value: -1}}
	case "created_at":
		sort = bson.D{{Key: "created_at", Value: 1}}
	case "-created_at":
		sort = bson.D{{Key: "created_at", Value: -1}}
	case "-name":
		sort = bson.D{{Key: "name", Value: -1}}
	}

	cursor, err := repository.collection.Find(ctx, filter, options.Find().
		SetSort(sort).
		SetSkip(int64((params.Page-1)*params.PerPage)).
		SetLimit(int64(params.PerPage)),
	)
	if err != nil {
		return ListResult{}, err
	}
	defer cursor.Close(ctx)

	couriers := make([]Courier, 0)
	if err := cursor.All(ctx, &couriers); err != nil {
		return ListResult{}, err
	}
	totalPages := total / int64(params.PerPage)
	if total%int64(params.PerPage) != 0 {
		totalPages++
	}

	return ListResult{Couriers: couriers, Page: params.Page, PerPage: params.PerPage, Total: total, TotalPages: totalPages}, nil
}

func (repository *Repository) FindByID(ctx context.Context, id bson.ObjectID) (Courier, error) {
	var courier Courier
	err := repository.collection.FindOne(ctx, bson.M{"_id": id, "deleted_at": nil}).Decode(&courier)
	return courier, err
}

func (repository *Repository) Update(ctx context.Context, id bson.ObjectID, courier Courier) (Courier, error) {
	set := bson.M{
		"name":          courier.Name,
		"level":         courier.Level,
		"status":        courier.Status,
		"registered_at": courier.RegisteredAt,
		"updated_at":    courier.UpdatedAt,
	}
	unset := bson.M{}
	applyOptionalString(set, unset, "email", courier.Email)
	applyOptionalString(set, unset, "phone", courier.Phone)
	applyOptionalString(set, unset, "vehicle_type", courier.VehicleType)
	applyOptionalString(set, unset, "license_plate", courier.LicensePlate)

	update := bson.M{"$set": set}
	if len(unset) > 0 {
		update["$unset"] = unset
	}
	_, err := repository.collection.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": nil}, update)
	if err != nil {
		return Courier{}, err
	}
	return repository.FindByID(ctx, id)
}

func applyOptionalString(set bson.M, unset bson.M, field string, value *string) {
	if value == nil {
		unset[field] = ""
		return
	}
	set[field] = value
}

func (repository *Repository) SoftDelete(ctx context.Context, id bson.ObjectID, deletedAt time.Time) error {
	result, err := repository.collection.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": nil}, bson.M{"$set": bson.M{"deleted_at": deletedAt, "updated_at": deletedAt}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (repository *Repository) buildFilter(params QueryParams) bson.M {
	filter := bson.M{"deleted_at": nil}
	if len(params.Levels) > 0 {
		filter["level"] = bson.M{"$in": params.Levels}
	}
	if strings.TrimSpace(params.Search) != "" {
		words := strings.Fields(params.Search)
		and := make([]bson.M, 0, len(words)+1)
		for _, word := range words {
			and = append(and, bson.M{"name": bson.M{"$regex": regexp.QuoteMeta(word), "$options": "i"}})
		}
		and = append(and, filter)
		return bson.M{"$and": and}
	}
	return filter
}
