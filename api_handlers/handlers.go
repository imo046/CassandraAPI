package api_handlers

import (
	"CassandraAPI/models"
	"CassandraAPI/utils"
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type Handler struct {
	S *gocql.Session
}

func (h Handler) HomeLink(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "My Restful API")
	utils.Panic(err, "Failed to write a message!")
}

func (h Handler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	var entry models.SessionEntry
	requestBody, err := io.ReadAll(r.Body)
	utils.Panic(err, "Failed to read request body!")
	unmarshallErr := json.Unmarshal(requestBody, &entry)
	utils.Panic(unmarshallErr, "Failed to process request body!")
	if insertErr := h.S.Query("INSERT INTO db_entries(entry_id,entry_val) VALUES (?, ?)", entry.EntryId, entry.EntryVal).Exec(); insertErr != nil {
		utils.Panic(insertErr, "Failed to insert data")
	}
	w.WriteHeader(http.StatusCreated)
	Conv, _ := json.MarshalIndent(entry, "", " ")
	fmt.Fprintf(w, "%s", string(Conv))

}

func (h Handler) GetEntries(w http.ResponseWriter, r *http.Request) {
	var entries []models.SessionEntry
	m := make(map[string]interface{})

	iter := h.S.Query("SELECT * FROM db_entries").Iter()
	for iter.MapScan(m) {
		entries = append(entries, models.SessionEntry{
			EntryId:  m["entry_id"].(string),
			EntryVal: m["entry_val"].(int),
		})
		m = map[string]interface{}{}
	}
	Conv, _ := json.MarshalIndent(entries, "", " ")
	fmt.Fprintf(w, "%s", string(Conv))
}

func (h Handler) GetEntry(w http.ResponseWriter, r *http.Request) {
	entryID := mux.Vars(r)["entry_id"]
	var entries []models.SessionEntry
	m := make(map[string]interface{})
	iter := h.S.Query("SELECT * FROM db_entries WHERE entry_id=?", entryID).Iter()
	for iter.MapScan(m) {
		entries = append(entries, models.SessionEntry{
			EntryId:  m["entry_id"].(string),
			EntryVal: m["entry_val"].(int),
		})
		m = map[string]interface{}{}
	}
	Conv, _ := json.MarshalIndent(entries, "", " ")
	fmt.Fprintf(w, "%s", string(Conv))
}

func (h Handler) CountEntries(w http.ResponseWriter, r *http.Request) {
	var Count string
	err := h.S.Query("SELECT count(*) FROM db_entries").Scan(&Count)
	utils.Panic(err, "Failed to count entries!")
	fmt.Fprintf(w, "%s", string(Count))
}

func (h Handler) DeleteOne(w http.ResponseWriter, r *http.Request) {
	entryID := mux.Vars(r)["entry_id"]
	if err := h.S.Query("DELETE FROM db_entries WHERE entry_id=?", entryID).Exec(); err != nil {
		utils.Panic(err, "Failed to delete entry!")
	}
	fmt.Fprintf(w, "Entry %s deleted!", entryID)
}

func (h Handler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if err := h.S.Query("TRUNCATE db_entries").Exec(); err != nil {
		utils.Panic(err, "Failed to delete entries!")
	}
	fmt.Fprintf(w, "All entries deleted!")
}
func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	entryID := mux.Vars(r)["entry_id"]
	var updated models.SessionEntry
	reqBody, err := io.ReadAll(r.Body)
	utils.Panic(err, "Failed to read request body!")
	json.Unmarshal(reqBody, &updated)
	if postErr := h.S.Query("UPDATE db_entries SET entry_val = ? WHERE entry_id = ?", updated.EntryVal, entryID).Exec(); postErr != nil {
		utils.Panic(postErr, "Failed to update entry!")
	}
	fmt.Fprintf(w, "Entry %s updated!", entryID)
}
