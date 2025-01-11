package recfile

import (
    "fmt"
    "strings"
    "text/template"
)

type FieldType uint8

const (
    FieldTypeString FieldType = iota
    FieldTypeStringMultiline
    FieldTypeInt
    FieldTypeFloat
    FieldTypeBool
    FieldTypeReference
    FieldTypeEnum
)

func FieldTypeFromString(fieldType string) FieldType {
    lower := strings.ToLower(fieldType)
    switch lower {
    case "line":
        return FieldTypeString
    case "multiline":
        return FieldTypeStringMultiline
    case "int":
        return FieldTypeInt
    case "real":
        return FieldTypeFloat
    case "bool":
        return FieldTypeBool
    case "reference":
        return FieldTypeReference
    default:
        return FieldTypeStringMultiline // default case, assuming string as default type
    }
}

type RecordSchema struct {
    RecordType   string
    NameFormat   string
    KeyFieldName string
    Fields       map[string]FieldDetails
    FieldOrder   []string
}

func (s RecordSchema) WithType(name string, fieldType FieldType) RecordSchema {
    if s.Fields == nil {
        s.Fields = make(map[string]FieldDetails)
    }
    name = strings.ToLower(name)
    if _, ok := s.Fields[name]; ok {
        s.Fields[name] = s.Fields[name].WithType(fieldType)
    } else {
        details := FieldDetails{
            Name: name,
            Type: fieldType,
        }
        s.Fields[name] = details
        s.FieldOrder = append(s.FieldOrder, name)
    }

    return s
}

func (s RecordSchema) WithListType(fieldName string) RecordSchema {
    fieldName = strings.ToLower(fieldName)
    newSchema := s
    if newSchema.Fields == nil {
        newSchema.Fields = make(map[string]FieldDetails)
    }
    if _, ok := newSchema.Fields[fieldName]; !ok {
        newSchema = newSchema.WithType(fieldName, FieldTypeString)
    }

    fieldDetails := newSchema.Fields[fieldName]
    fieldDetails.IsList = true
    newSchema.Fields[fieldName] = fieldDetails
    return newSchema
}

func (s RecordSchema) WithEnum(name string, enumValues []string) RecordSchema {
    name = strings.ToLower(name)
    details := FieldDetails{
        Name:       name,
        Type:       FieldTypeEnum,
        EnumValues: enumValues,
    }
    s.Fields[name] = details
    s.FieldOrder = append(s.FieldOrder, name)
    return s
}

func (s RecordSchema) WithKeyFieldName(keyFieldName string) RecordSchema {
    keyFieldName = strings.ToLower(keyFieldName)
    s.KeyFieldName = keyFieldName
    return s
}

func (s RecordSchema) WithReference(name string, referenceSchema string) RecordSchema {
    name = strings.ToLower(name)
    details := FieldDetails{
        Name:            name,
        Type:            FieldTypeReference,
        ReferenceSchema: referenceSchema,
    }
    s.Fields[name] = details
    s.FieldOrder = append(s.FieldOrder, name)
    return s
}

func (s RecordSchema) WithNameFormatFromKeyField() RecordSchema {
    s.NameFormat = fmt.Sprintf("{{ .%s}}", s.KeyFieldName)
    return s
}

func (s RecordSchema) DisplayNameFor(rec Record) string {
    asMap := rec.ToMap(",")
    parsedTemplate, err := template.New("text").Parse(s.NameFormat)
    if err != nil {
        panic(err)
    }
    var filledText strings.Builder
    err = parsedTemplate.Execute(&filledText, asMap)
    if err != nil {
        panic(err)
    }
    return filledText.String()
}

func (s RecordSchema) IsEmpty() bool {
    return len(s.Fields) == 0
}

type FieldDetails struct {
    Name            string
    Type            FieldType
    EnumValues      []string
    ReferenceSchema string
    IsList          bool
}

func (d FieldDetails) WithType(fieldType FieldType) FieldDetails {
    d.Type = fieldType
    return d
}
