package model

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDateTimeToPlainUsesUTC(t *testing.T) {
	raw, err := bson.Marshal(bson.D{{Key: "t", Value: primitive.NewDateTimeFromTime(time.Date(2026, 5, 29, 21, 52, 36, 0, time.UTC))}})
	if err != nil {
		t.Fatal(err)
	}
	var doc bson.D
	if err := bson.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	tm := ToPlain(doc.Map()["t"]).(time.Time)
	if tm.Location() != time.UTC {
		t.Fatalf("location=%v", tm.Location())
	}
}
