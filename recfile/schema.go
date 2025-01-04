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
    details := FieldDetails{
        Name: name,
        Type: fieldType,
    }
    if s.Fields == nil {
        s.Fields = make(map[string]FieldDetails)
    }
    s.Fields[name] = details
    s.FieldOrder = append(s.FieldOrder, name)
    return s
}

func (s RecordSchema) WithEnum(name string, enumValues []string) RecordSchema {
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
    s.KeyFieldName = keyFieldName
    return s
}

func (s RecordSchema) WithReference(name string, referenceSchema string) RecordSchema {
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
}
