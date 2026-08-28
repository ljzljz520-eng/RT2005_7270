package model

import "time"

type Record struct {
	ID, Title, Body, Status string
	StartAt, Deadline       time.Time
	ProfileID               string
	Version                 int
}
type Profile struct {
	ID, Name, Email, Timezone string
	Active                    bool
}
type Event struct {
	ID, RecordID, Kind, Message string
	At                          time.Time
}
type Audit struct {
	ID, Actor, Action, Target string
	At                        time.Time
	Metadata                  map[string]string
}

func (r Record) IsExpired(now time.Time) bool  { return !r.Deadline.IsZero() && now.After(r.Deadline) }
func (r Record) CanPublish(now time.Time) bool { return r.Status == "approved" && !r.IsExpired(now) }
func NewRecord(id, title, body string, deadline time.Time) Record {
	return Record{ID: id, Title: title, Body: body, Status: "draft", Deadline: deadline, Version: 1}
}
func NewProfile(id, name, email string) Profile {
	return Profile{ID: id, Name: name, Email: email, Active: true}
}
func NewEvent(id, recordID, kind, message string) Event {
	return Event{ID: id, RecordID: recordID, Kind: kind, Message: message, At: time.Now().UTC()}
}
