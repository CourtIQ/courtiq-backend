package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/CourtIQ/courtiq-backend/search-service/graph/model"
	"github.com/CourtIQ/courtiq-backend/search-service/internal/utils"
	"github.com/CourtIQ/courtiq-backend/shared/pkg/db"
)

// --------------------------
// INTERFACE + CONSTRUCTOR
// --------------------------
type SearchRepository interface {
	SearchUsers(ctx context.Context, query string, excludeUserID primitive.ObjectID, limit, offset int) ([]*model.UserSearchResult, error)
	SearchTennisCourts(ctx context.Context, query string, lat, lng float64, radius float64, limit, offset int) ([]*model.TennisCourtSearchResult, error)
}

type searchRepository struct {
	usersCollection        *mongo.Collection
	tennisCourtsCollection *mongo.Collection
}

func NewSearchRepository(mdb *db.MongoDB) SearchRepository {
	return &searchRepository{
		usersCollection:        mdb.GetCollection(db.UsersCollection),
		tennisCourtsCollection: mdb.GetCollection(db.TennisCourtsCollection),
	}
}

// --------------------------
// SEARCH USERS
// --------------------------
func (r *searchRepository) SearchUsers(
	ctx context.Context,
	query string,
	excludeUserID primitive.ObjectID,
	limit, offset int,
) ([]*model.UserSearchResult, error) {

	log.Printf("[SearchUsers] query=%q excludeUserID=%s limit=%d offset=%d", query, excludeUserID.Hex(), limit, offset)
	pipeline := utils.BuildUserSearchPipeline(query, excludeUserID, limit, offset)

	cursor, err := r.usersCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate error: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []model.UserSearchResult
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	results := make([]*model.UserSearchResult, len(docs))
	for i := range docs {
		results[i] = &docs[i]
	}
	log.Printf("[SearchUsers] found %d user(s)", len(results))
	return results, nil
}

// --------------------------
// SEARCH TENNIS COURTS
// --------------------------
// We'll first do an aggregate pipeline search in Mongo. If we find fewer than 5 results,
// we then fetch more from Google, insert them, and re-query.
func (r *searchRepository) SearchTennisCourts(
	ctx context.Context,
	query string,
	lat, lng float64,
	radius float64,
	limit, offset int,
) ([]*model.TennisCourtSearchResult, error) {

	log.Printf("[SearchTennisCourts] query=%q lat=%.6f lng=%.6f radius=%.1f limit=%d offset=%d",
		query, lat, lng, radius, limit, offset)

	// We'll build a pipeline with a $geoNear stage plus name regex match, skip & limit
	// e.g. maxDistance = radius * 1000 if you're converting km to meters
	maxDistance := radius // if radius is already in meters, else do radius * 1000
	pipeline := utils.BuildTennisCourtLocationNamePipeline(
		lat, lng,
		maxDistance,
		query,
		limit,
		offset,
	)

	// Step 1) Query DB
	cursor, err := r.tennisCourtsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate error: %w", err)
	}
	defer cursor.Close(ctx)

	var existingCourts []model.TennisCourt
	if err := cursor.All(ctx, &existingCourts); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	existingCount := len(existingCourts)
	log.Printf("[SearchTennisCourts] Found %d courts from pipeline for query=%q", existingCount, query)

	// If we found 5 or more in DB, just return them
	if existingCount >= 5 {
		log.Printf("[SearchTennisCourts] >= 5 in DB, skipping Google fetch.")
		return r.mapTennisCourtsToResults(existingCourts), nil
	}

	// Step 2) If fewer than 5 found, fetch from Google
	numToFetch := 10 - existingCount // we'll aim for up to 10 total
	if numToFetch < 1 {
		// If for some reason existing is 4, we fetch 6, etc. But never <1
		numToFetch = 1
	}
	log.Printf("[SearchTennisCourts] only %d in DB, fetching up to %d from Google", existingCount, numToFetch)

	googlePlaces, err := r.fetchFromGoogleAutocomplete(ctx, query, lat, lng, radius, numToFetch)
	if err != nil {
		log.Printf("[SearchTennisCourts] google fetch error: %v", err)
	}

	// Build a set of already-known placeIDs to avoid double insertion
	existingMap := make(map[string]bool)
	for _, c := range existingCourts {
		existingMap[c.GooglePlaceID] = true
	}

	// Step 3) For each new place from Google, if not in DB => fetch details & insert
	for _, pid := range googlePlaces {
		if existingMap[pid] {
			continue
		}

		details, err := r.fetchPlaceDetails(ctx, pid)
		if err != nil {
			log.Printf("[SearchTennisCourts] fetchPlaceDetails error pid=%s: %v", pid, err)
			continue
		}

		now := time.Now()
		tc := model.TennisCourt{
			ID:                       primitive.NewObjectID(),
			GooglePlaceID:            pid,
			Name:                     details.Name,
			Coordinates:              [2]float64{details.Geometry.Location.Lng, details.Geometry.Location.Lat},
			FormattedAddress:         &details.FormattedAddress,
			City:                     &details.City,
			State:                    &details.State,
			Country:                  &details.Country,
			PostalCode:               &details.PostalCode,
			Rating:                   &details.Rating,
			UserRatingsTotal:         &details.UserRatingsTotal,
			BusinessStatus:           &details.BusinessStatus,
			PhoneNumber:              &details.FormattedPhoneNumber,
			InternationalPhoneNumber: &details.InternationalPhoneNumber,
			Website:                  &details.Website,
			OpeningHours:             details.toModelOpeningHours(),
			OpenNow:                  &details.OpenNow,
			LastUpdated:              &now,
		}

		if _, insertErr := r.tennisCourtsCollection.InsertOne(ctx, tc); insertErr != nil {
			log.Printf("[SearchTennisCourts] Insert error pid=%s: %v", pid, insertErr)
			continue
		}
		existingMap[pid] = true
	}

	// Step 4) Re-run the pipeline or a simpler find to get the final data (limit=10).
	// Let's re-run the pipeline but with limit=10 for the final results,
	// ignoring offset for this final check, or you can do offset if you want.
	finalPipeline := utils.BuildTennisCourtLocationNamePipeline(
		lat, lng,
		maxDistance,
		query,
		10,
		0, // maybe ignore offset here, or reuse offset if you prefer
	)
	finalCursor, err := r.tennisCourtsCollection.Aggregate(ctx, finalPipeline)
	if err != nil {
		return nil, fmt.Errorf("final aggregate error: %w", err)
	}
	defer finalCursor.Close(ctx)

	var finalCourts []model.TennisCourt
	if err := finalCursor.All(ctx, &finalCourts); err != nil {
		return nil, fmt.Errorf("decode final courts error: %w", err)
	}

	results := r.mapTennisCourtsToResults(finalCourts)
	log.Printf("[SearchTennisCourts] returning %d tennis courts total after merge", len(results))
	return results, nil
}

// -----------------------------------
// fetchFromGoogleAutocomplete
// -----------------------------------
func (r *searchRepository) fetchFromGoogleAutocomplete(
	ctx context.Context,
	query string,
	lat, lng float64,
	radius float64,
	numToFetch int,
) ([]string, error) {

	googleAPIKey := os.Getenv("GOOGLE_PLACES_API_KEY")
	if googleAPIKey == "" {
		return nil, fmt.Errorf("missing GOOGLE_PLACES_API_KEY env variable")
	}

	autoURL, err := url.Parse("https://maps.googleapis.com/maps/api/place/autocomplete/json")
	if err != nil {
		return nil, fmt.Errorf("url parse error: %w", err)
	}

	q := autoURL.Query()
	// Force "tennis court" in the query
	q.Set("input", "tennis court "+query)
	q.Set("key", googleAPIKey)
	// Use "radius" in meters, up to 50 km, etc.
	q.Set("location", fmt.Sprintf("%.6f,%.6f", lat, lng))
	q.Set("radius", "50000")
	q.Set("strictbounds", "true")

	autoURL.RawQuery = q.Encode()

	resp, err := http.Get(autoURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to call autocomplete: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("autocomplete returned %d: %s", resp.StatusCode, string(body))
	}

	var autoResp struct {
		Predictions []struct {
			PlaceID string `json:"place_id"`
		} `json:"predictions"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&autoResp); err != nil {
		return nil, fmt.Errorf("decode autocomplete response: %w", err)
	}
	if autoResp.Status != "OK" && autoResp.Status != "ZERO_RESULTS" {
		return nil, fmt.Errorf("autocomplete status %q", autoResp.Status)
	}
	if len(autoResp.Predictions) == 0 {
		log.Println("[SearchTennisCourts] No predictions from Autocomplete.")
		return []string{}, nil
	}

	// Collect place IDs, optionally truncating to `numToFetch`
	var placeIDs []string
	for i, p := range autoResp.Predictions {
		if i >= numToFetch {
			break
		}
		placeIDs = append(placeIDs, p.PlaceID)
	}
	return placeIDs, nil
}

// -----------------------------------
// mapTennisCourtsToResults
// -----------------------------------
func (r *searchRepository) mapTennisCourtsToResults(courts []model.TennisCourt) []*model.TennisCourtSearchResult {
	results := make([]*model.TennisCourtSearchResult, 0, len(courts))
	for _, c := range courts {
		var addr string
		if c.FormattedAddress != nil {
			addr = *c.FormattedAddress
		}
		results = append(results, &model.TennisCourtSearchResult{
			ID:            c.ID,
			GooglePlaceID: c.GooglePlaceID,
			DisplayName:   c.Name,
			Address:       addr,
			City:          c.City,
			Country:       c.Country,
			Rating:        c.Rating,
			OpenNow:       c.OpenNow,
			Coordinates:   c.Coordinates,
		})
	}
	return results
}

// -------------------------------------------------------------------
// fetchPlaceDetails
// -------------------------------------------------------------------
func (r *searchRepository) fetchPlaceDetails(ctx context.Context, placeID string) (*GooglePlaceDetails, error) {
	googleAPIKey := os.Getenv("GOOGLE_PLACES_API_KEY")
	if googleAPIKey == "" {
		return nil, fmt.Errorf("missing GOOGLE_PLACES_API_KEY env variable")
	}

	detailsURL, err := url.Parse("https://maps.googleapis.com/maps/api/place/details/json")
	if err != nil {
		return nil, err
	}
	q := detailsURL.Query()
	q.Set("place_id", placeID)
	q.Set("key", googleAPIKey)
	// Optionally: q.Set("fields", "place_id,name,geometry,formatted_address,...")
	detailsURL.RawQuery = q.Encode()

	resp, err := http.Get(detailsURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("place details returned %d: %s", resp.StatusCode, string(body))
	}

	var raw struct {
		Result GooglePlaceDetails `json:"result"`
		Status string             `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode place details: %w", err)
	}
	if raw.Status != "OK" {
		return nil, fmt.Errorf("place details status %q", raw.Status)
	}
	res := &raw.Result
	res.parseAddressComponents()
	return res, nil
}

// -------------------------------------------------------------------
// GooglePlaceDetails + Helpers
// -------------------------------------------------------------------
type GooglePlaceDetails struct {
	PlaceID                  string              `json:"place_id"`
	Name                     string              `json:"name"`
	Geometry                 geometryType        `json:"geometry"`
	FormattedAddress         string              `json:"formatted_address"`
	FormattedPhoneNumber     string              `json:"formatted_phone_number"`
	InternationalPhoneNumber string              `json:"international_phone_number"`
	Website                  string              `json:"website"`
	BusinessStatus           string              `json:"business_status"`
	Rating                   float64             `json:"rating"`
	UserRatingsTotal         int                 `json:"user_ratings_total"`
	OpeningHours             *googleOpeningHours `json:"opening_hours"`
	AddressComponents        []AddressComponent  `json:"address_components"`
	Types                    []string            `json:"types"`

	City        string
	State       string
	Country     string
	PostalCode  string
	OpenNow     bool
	LastUpdated time.Time
}

type geometryType struct {
	Location struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"location"`
}

type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type googleOpeningHours struct {
	OpenNow     bool           `json:"open_now"`
	WeekdayText []string       `json:"weekday_text"`
	Periods     []googlePeriod `json:"periods"`
}

type googlePeriod struct {
	Open  googlePeriodTime `json:"open"`
	Close googlePeriodTime `json:"close"`
}

type googlePeriodTime struct {
	Day  int    `json:"day"`
	Time string `json:"time"`
}

// parseAddressComponents sets City, State, Country, PostalCode from address components.
func (g *GooglePlaceDetails) parseAddressComponents() {
	for _, comp := range g.AddressComponents {
		for _, t := range comp.Types {
			switch t {
			case "locality":
				g.City = comp.LongName
			case "administrative_area_level_1":
				g.State = comp.ShortName
			case "country":
				g.Country = comp.LongName
			case "postal_code":
				g.PostalCode = comp.LongName
			}
		}
	}
	g.LastUpdated = time.Now()
	if g.OpeningHours != nil {
		g.OpenNow = g.OpeningHours.OpenNow
	}
}

// Converts the googleOpeningHours into our GraphQL model.OpeningHours
func (g *GooglePlaceDetails) toModelOpeningHours() *model.OpeningHours {
	if g.OpeningHours == nil {
		return nil
	}
	wd := g.OpeningHours.WeekdayText
	ps := make([]*model.OpeningPeriod, 0, len(g.OpeningHours.Periods))
	for _, p := range g.OpeningHours.Periods {
		openTime := p.Open.Time
		closeTime := p.Close.Time
		day := p.Open.Day
		ps = append(ps, &model.OpeningPeriod{
			Day:       &day,
			OpenTime:  &openTime,
			CloseTime: &closeTime,
		})
	}
	return &model.OpeningHours{
		WeekdayText: wd,
		Periods:     ps,
	}
}
