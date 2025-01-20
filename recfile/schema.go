package recfile

import (
    "fmt"
    "strings"
    "text/template"
)

type FieldType uint8

func (t FieldType) String() string {
    switch t {
    case FieldTypeString:
        return "line"
    case FieldTypeStringMultiline:
        return "multiline"
    case FieldTypeInt:
        return "int"
    case FieldTypeFloat:
        return "real"
    case FieldTypeBool:
        return "bool"
    case FieldTypeReference:
        return "reference"
    case FieldTypeEnum:
        return "enum"
    default:
        return "line"
    }
}

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

func (s RecordSchema) ToRecord() Record {
    rec := Record{
        Field{Name: "%rec", Value: s.RecordType},
    }
    if s.KeyFieldName != "" {
        rec = append(rec, Field{Name: "%key", Value: s.KeyFieldName})
    }
    if s.NameFormat != "" {
        rec = append(rec, Field{Name: "%label", Value: s.NameFormat})
    }
    for _, fieldName := range s.FieldOrder {
        fieldDetails := s.Fields[fieldName]
        if fieldDetails.Type == FieldTypeEnum {
            enumTypeDef := fmt.Sprintf("%s enum %s", fieldName, strings.Join(fieldDetails.EnumValues, " "))
            rec = append(rec, Field{Name: "%typedef", Value: enumTypeDef})
        } else if fieldDetails.Type == FieldTypeReference {
            rec = append(rec, Field{Name: "%ref", Value: fmt.Sprintf("%s %s", fieldName, strings.Join(fieldDetails.ReferencedSchemas, " "))})
        } else {
            rec = append(rec, Field{Name: "%type", Value: fmt.Sprintf("%s %s", fieldName, fieldDetails.Type.String())})
        }
        if fieldDetails.IsList {
            rec = append(rec, Field{Name: "%list", Value: fieldName})
        }
    }
    return rec
}
func (s RecordSchema) WithType(name string, fieldType FieldType) RecordSchema {
    if s.Fields == nil {
        s.Fields = make(map[string]FieldDetails)
    }
    if _, ok := s.Fields[name]; ok {
        s.Fields[name] = s.Fields[name].WithType(fieldType)
        return s
    }

    details := FieldDetails{
        Name: name,
        Type: fieldType,
    }
    s.Fields[name] = details
    s.FieldOrder = append(s.FieldOrder, name)

    return s
}

func (s RecordSchema) WithListType(fieldName string) RecordSchema {
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
    if s.Fields == nil {
        s.Fields = make(map[string]FieldDetails)
    }
    if _, ok := s.Fields[name]; ok {
        s.Fields[name] = s.Fields[name].WithEnum(name, enumValues)
        return s
    }

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
    return s.WithMultiReference(name, []string{referenceSchema})
}

func (s RecordSchema) WithMultiReference(name string, referencedSchemas []string) RecordSchema {
    if s.Fields == nil {
        s.Fields = make(map[string]FieldDetails)
    }
    if _, ok := s.Fields[name]; ok {
        s.Fields[name] = s.Fields[name].WithMultiReference(name, referencedSchemas)
        return s
    }
    details := FieldDetails{
        Name:              name,
        Type:              FieldTypeReference,
        ReferencedSchemas: referencedSchemas,
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

func (s RecordSchema) WithNameFormat(nameFormat string) RecordSchema {
    s.NameFormat = nameFormat
    return s
}

type FieldDetails struct {
    Name              string
    Type              FieldType
    EnumValues        []string
    ReferencedSchemas []string
    IsList            bool
}

func (d FieldDetails) WithType(fieldType FieldType) FieldDetails {
    d.Type = fieldType
    return d
}

func (d FieldDetails) WithEnum(name string, values []string) FieldDetails {
    d.Name = name
    d.Type = FieldTypeEnum
    d.EnumValues = values
    return d
}

func (d FieldDetails) WithReference(name string, schema string) FieldDetails {
    d.Name = name
    d.Type = FieldTypeReference
    d.ReferencedSchemas = []string{schema}
    return d
}

func (d FieldDetails) WithMultiReference(name string, schemas []string) FieldDetails {
    d.Name = name
    d.Type = FieldTypeReference
    d.ReferencedSchemas = schemas
    return d
}
