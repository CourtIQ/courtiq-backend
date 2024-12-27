// internal/utils/pipeline_builder.go

package utils

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BuildUserSearchPipeline returns a mongo.Pipeline for searching the "users" collection
// with the given text query, excluding `excludeUserID`, applying fuzzy/partial matches, etc.
func BuildUserSearchPipeline(
	query string,
	excludeUserID primitive.ObjectID,
	limit, offset int,
) mongo.Pipeline {
	searchStage := bson.D{
		{
			"$search", bson.D{
				{"index", "users_search"},
				{"compound", bson.D{
					{"should", bson.A{
						// Exact username match (highest score)
						bson.D{
							{"text", bson.D{
								{"query", query},
								{"path", "username"},
								{"score", bson.D{{"boost", bson.D{{"value", 5}}}}},
							}},
						},
						// Autocomplete on all name fields with fuzzy matching
						bson.D{
							{"compound", bson.D{
								{"should", bson.A{
									bson.D{
										{"autocomplete", bson.D{
											{"query", query},
											{"path", "firstName"},
											{"fuzzy", bson.D{
												{"maxEdits", 1},
												{"prefixLength", 1},
											}},
											{"score", bson.D{{"boost", bson.D{{"value", 3}}}}},
										}},
									},
									bson.D{
										{"autocomplete", bson.D{
											{"query", query},
											{"path", "lastName"},
											{"fuzzy", bson.D{
												{"maxEdits", 1},
												{"prefixLength", 1},
											}},
											{"score", bson.D{{"boost", bson.D{{"value", 3}}}}},
										}},
									},
									bson.D{
										{"autocomplete", bson.D{
											{"query", query},
											{"path", "displayName"},
											{"fuzzy", bson.D{
												{"maxEdits", 1},
												{"prefixLength", 1},
											}},
											{"score", bson.D{{"boost", bson.D{{"value", 2}}}}},
										}},
									},
								}},
							}},
						},
					}},
				}},
			},
		},
	}

	// Rest of your pipeline stages remain the same
	matchStage := bson.D{
		{
			"$match", bson.D{
				{"_id", bson.D{{"$ne", excludeUserID}}},
			},
		},
	}

	skipStage := bson.D{{"$skip", offset}}
	limitStage := bson.D{{"$limit", limit}}

	projectStage := bson.D{
		{
			"$project", bson.D{
				{"id", "$_id"},
				{"username", 1},
				{"displayName", 1},
				{"firstName", 1},
				{"lastName", 1},
				{"profilePicture", 1},
				{"location", 1},
				{"score", bson.D{{"$meta", "searchScore"}}},
			},
		},
	}

	return mongo.Pipeline{
		searchStage,
		matchStage,
		skipStage,
		limitStage,
		projectStage,
	}
}

func BuildTennisCourtLocationNamePipeline(
    lat, lng float64,
    maxDistance float64,
    query string,
    limit int,
	offset int,
) mongo.Pipeline {

    pipeline := mongo.Pipeline{}

    // 1) $geoNear stage
    geoNearStage := bson.D{{
        "$geoNear", bson.D{
            {"near", bson.D{
                {"type", "Point"},
                {"coordinates", bson.A{lng, lat}}, // [lng, lat]
            }},
            {"distanceField", "distance"}, 
            {"spherical", true},
            // optionally: {"maxDistance", maxDistance},
        },
    }}

    pipeline = append(pipeline, geoNearStage)

    // 2) $match stage for name ~ query
    if query != "" {
        matchStage := bson.D{{
            "$match", bson.M{
                "name": bson.M{"$regex": query, "$options": "i"},
            },
        }}
        pipeline = append(pipeline, matchStage)
    }

    if offset > 0 {
        pipeline = append(pipeline, bson.D{{"$skip", offset}})
    }

    // 3) $limit stage
    if limit > 0 {
        pipeline = append(pipeline, bson.D{{"$limit", limit}})
    }

    return pipeline
}

