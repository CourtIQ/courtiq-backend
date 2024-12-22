package utils

import "go.mongodb.org/mongo-driver/mongo/options"

func BuildFindOptions(limit *int, offset *int) *options.FindOptions {
	// Default values
	defaultLimit := 10
	defaultOffset := 0

	findOptions := &options.FindOptions{}

	// If limit is nil or < 1, use the default limit
	if limit == nil || *limit < 1 {
		findOptions.SetLimit(int64(defaultLimit))
	} else {
		findOptions.SetLimit(int64(*limit))
	}

	// If offset is nil or < 0, use the default offset
	if offset == nil || *offset < 0 {
		findOptions.SetSkip(int64(defaultOffset))
	} else {
		findOptions.SetSkip(int64(*offset))
	}

	return findOptions
}
