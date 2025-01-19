package recfile

import (
    "fmt"
    "os"
    "slices"
    "strings"
)

type RecordRepository struct {
    Schema      RecordSchema
    Records     map[string]Record
    IDs         []string
    BackingFile string
}

func (r *RecordRepository) AddEnumToSchema(fieldName string, values []string) {
    r.Schema = r.Schema.WithEnum(fieldName, values)
}

func (r *RecordRepository) Merge(records []Record) {
    for _, record := range records {
        id := record.FindValueForKeyIgnoreCase(r.Schema.KeyFieldName)
        if _, exists := r.Records[id]; exists {
            r.ReplaceById(id, record)
        } else {
            r.Append(record)
        }
    }
}

func (r *RecordRepository) ReplaceById(originalId string, record Record) {
    currentRecordId := record.FindValueForKeyIgnoreCase(r.Schema.KeyFieldName)
    if currentRecordId == originalId {
        r.Records[originalId] = record
    } else { // ID has changed
        r.DeleteById(originalId)
        r.Append(record)
    }
}

func (r *RecordRepository) DeleteById(id string) {
    delete(r.Records, id)
    for i, recordID := range r.IDs {
        if recordID == id {
            r.IDs = append(r.IDs[:i], r.IDs[i+1:]...)
            return
        }
    }
}

func (r *RecordRepository) FindById(id string) Record {
    return r.Records[id]
}

func (r *RecordRepository) FindByField(fieldName string, value string) []Record {
    var matchingRecords []Record
    for _, record := range r.Records {
        if record.FindValueForKeyIgnoreCase(fieldName) == value {
            matchingRecords = append(matchingRecords, record)
        }
    }
    return matchingRecords
}

func (r *RecordRepository) Count() int {
    return len(r.Records)
}

func (r *RecordRepository) Append(record Record) {
    id := record.FindValueForKeyIgnoreCase(r.Schema.KeyFieldName)
    if _, exists := r.Records[id]; exists {
        return
    }
    r.Records[id] = record
    r.IDs = append(r.IDs, id)
}

func (r *RecordRepository) SortByID() {
    slices.Sort(r.IDs)
}

func (r *RecordRepository) Write() error {
    var recordsAsList []Record
    for _, id := range r.IDs {
        record := r.Records[id]
        recordsAsList = append(recordsAsList, record)
    }
    file, err := os.Create(r.BackingFile)
    if err != nil {
        return err
    }
    defer file.Close()
    return WriteWithSchema(file, recordsAsList, r.Schema)
}

func NewRecordRepository(filename string) *RecordRepository {
    file, _ := os.Open(filename)
    records, schema := ReadAndClose(file)

    repo := &RecordRepository{
        Schema:      schema,
        BackingFile: filename,
        Records:     make(map[string]Record, len(records)),
        IDs:         make([]string, len(records)),
    }
    for recordIndex, record := range records {
        id := record.FindValueForKeyIgnoreCase(schema.KeyFieldName)
        repo.Records[id] = record
        repo.IDs[recordIndex] = id
    }
    return repo
}

type RepoMan struct {
    repos map[string]*RecordRepository
}

func NewRepoMan() *RepoMan {
    return &RepoMan{
        repos: make(map[string]*RecordRepository),
    }
}

func (r *RepoMan) GetRepo(schemaName string) *RecordRepository {
    return r.repos[schemaName]
}

func (r *RepoMan) AddRepo(fileName string) string {
    repo := NewRecordRepository(fileName)
    schemaName := repo.Schema.RecordType
    r.repos[schemaName] = repo
    return schemaName
}

func (r *RepoMan) Status() string {
    if len(r.repos) == 0 {
        return "No records in repo"
    }

    var lines []string
    for schemaName, repo := range r.repos {
        reportLine := fmt.Sprintf("%s: %d records", schemaName, repo.Count())
        lines = append(lines, reportLine)
    }

    return fmt.Sprintf("RepoMan with %d schemas:\n%s", len(r.repos), strings.Join(lines, "\n"))
}
