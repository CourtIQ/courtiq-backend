package utils

import "go.mongodb.org/mongo-driver/mongo/options"

func BuildFindOptions(limit *int, offset *int) *options.FindOptions {
	findOptions := &options.FindOptions{}
	if limit != nil {
		findOptions.SetLimit(int64(*limit))
	}
	if offset != nil {
		findOptions.SetSkip(int64(*offset))
	}
	return findOptions
}
