package recfile

import (
    "fmt"
    "io"
    "os"
    "slices"
    "strings"
)

type RecordRepository struct {
    Schema  RecordSchema
    Records map[string]Record
    IDs     []string
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

func (r *RecordRepository) ReplaceById(id string, record Record) {
    r.Records[id] = record
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

func (r *RecordRepository) Write(w io.Writer) error {
    var recordsAsList []Record
    for _, id := range r.IDs {
        record := r.Records[id]
        recordsAsList = append(recordsAsList, record)
    }

    return WriteWithSchema(w, recordsAsList, r.Schema)
}

func NewRecordRepository(recordsAsList []Record, schema RecordSchema) *RecordRepository {
    repo := &RecordRepository{
        Schema:  schema,
        Records: make(map[string]Record, len(recordsAsList)),
        IDs:     make([]string, len(recordsAsList)),
    }
    for recordIndex, record := range recordsAsList {
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

func (r *RepoMan) AddRepos(fileName string) []string {
    var addedRepos []string
    file, _ := os.Open(fileName)
    recordsByType, schemas := ReadMultiAndClose(file)
    for schemaName, records := range recordsByType {
        if _, exists := r.repos[schemaName]; exists {
            r.repos[schemaName].Merge(records)
            if r.repos[schemaName].Schema.IsEmpty() && !schemas[schemaName].IsEmpty() {
                r.repos[schemaName].Schema = schemas[schemaName]
            }
        } else {
            schema := schemas[schemaName]
            r.repos[schemaName] = NewRecordRepository(records, schema)
            addedRepos = append(addedRepos, schemaName)
        }
    }
    return addedRepos
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
