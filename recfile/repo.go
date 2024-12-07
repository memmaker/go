package recfile

import (
    "fmt"
    "io"
    "strings"
)

type RecordRepository struct {
    Schema  RecordSchema
    Records map[string]Record
}

func (r *RecordRepository) Merge(records []Record) {
    for _, record := range records {
        id := record.FindValueForKeyIgnoreCase(r.Schema.KeyFieldName)
        r.Records[id] = record
    }
}

func (r *RecordRepository) Count() int {
    return len(r.Records)
}

func NewRecordRepository(recordsAsList []Record, schema RecordSchema) *RecordRepository {
    repo := &RecordRepository{
        Schema:  schema,
        Records: make(map[string]Record, len(recordsAsList)),
    }
    for _, record := range recordsAsList {
        id := record.FindValueForKeyIgnoreCase(schema.KeyFieldName)
        repo.Records[id] = record
    }
    return repo
}

type RepoMan struct {
    Repos map[string]*RecordRepository
}

func NewRepoMan() *RepoMan {
    return &RepoMan{
        Repos: make(map[string]*RecordRepository),
    }
}

func (r *RepoMan) AddRepo(recordFile io.ReadCloser) {
    recordsByType, schemas := ReadMultiAndClose(recordFile)
    for schemaName, records := range recordsByType {
        if _, exists := r.Repos[schemaName]; exists {
            r.Repos[schemaName].Merge(records)
            if r.Repos[schemaName].Schema.IsEmpty() && !schemas[schemaName].IsEmpty() {
                r.Repos[schemaName].Schema = schemas[schemaName]
            }
        } else {
            schema := schemas[schemaName]
            r.Repos[schemaName] = NewRecordRepository(records, schema)
        }
    }
}

func (r *RepoMan) Status() string {
    if len(r.Repos) == 0 {
        return "No records in repo"
    }

    var lines []string
    for schemaName, repo := range r.Repos {
        reportLine := fmt.Sprintf("%s: %d records", schemaName, repo.Count())
        lines = append(lines, reportLine)
    }

    return fmt.Sprintf("RepoMan with %d schemas:\n%s", len(r.Repos), strings.Join(lines, "\n"))
}
