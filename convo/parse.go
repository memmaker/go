package convo

import (
    "github.com/Knetic/govaluate"
    "github.com/memmaker/go/recfile"
    "os"
    "strings"
)

func ParseConversation(filename string, handler DialogueHandler) (*Conversation, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    records, _ := recfile.ReadMulti(file)
    conversation := NewConversation(handler)
    conversation.LoadScriptFuncs = len(handler.GetScriptFuncs()) > 0

    openingBranches := make([]OpeningBranch, 0)
    variables := make(map[string]interface{})

    if records["Variables"] != nil && conversation.LoadScriptFuncs {
        for _, variableRecord := range records["Variables"] {
            for _, field := range variableRecord {
                expression, parseErr := govaluate.NewEvaluableExpressionWithFunctions(field.Value, handler.GetScriptFuncs())
                if parseErr != nil {
                    return nil, parseErr
                }
                varValue, evalErr := expression.Evaluate(nil)
                if evalErr != nil {
                    return nil, evalErr
                }
                variables[field.Name] = varValue
            }
        }
        conversation.Variables = variables
    }

    for _, branchRecords := range records["OpeningBranch"] {
        var branch OpeningBranch
        for _, fields := range branchRecords {
            fieldName := strings.ToLower(fields.Name)
            if fieldName == "cond" && conversation.LoadScriptFuncs {
                cond, parseErr := govaluate.NewEvaluableExpressionWithFunctions(fields.Value, handler.GetScriptFuncs())
                if parseErr != nil {
                    return nil, parseErr
                }
                branch.BranchCondition = cond
            } else if fieldName == "goto" {
                branch.BranchName = fields.Value
            } else if fieldName == "name" {
                branch.Name = fields.Value
            }
        }
        openingBranches = append(openingBranches, branch)
    }
    conversation.openingBranches = openingBranches
    allNodes := make(map[string]ConversationNode)
    for _, nodeRecord := range records["Nodes"] {
        var conversationNode ConversationNode
        var currentOption conversationOption
        for _, field := range nodeRecord {
            fieldName := strings.ToLower(field.Name)
            if fieldName == "name" {
                conversationNode.Name = field.Value
            } else if fieldName == "npc" {
                conversationNode.NpcText = strings.TrimSpace(field.Value)
            } else if fieldName == "effect" {
                conversationNode.Effects = append(conversationNode.Effects, field.Value)
            } else if strings.HasPrefix(fieldName, "o_") {
                if fieldName == "o_text" {
                    if currentOption.PlayerText != "" {
                        conversationNode.Options = append(conversationNode.Options, currentOption)
                    }
                    currentOption.PlayerText = strings.TrimSpace(field.Value)
                    currentOption.branchCondition = nil
                    currentOption.successBranch = ""
                    currentOption.failureBranch = ""
                    currentOption.displayCondition = nil
                    currentOption.Effect = ""
                } else if fieldName == "o_cond" && conversation.LoadScriptFuncs {
                    dispCond, parseErr := govaluate.NewEvaluableExpressionWithFunctions(field.Value, handler.GetScriptFuncs())
                    if parseErr != nil {
                        return nil, parseErr
                    }
                    currentOption.displayCondition = dispCond
                } else if fieldName == "o_goto" || fieldName == "o_succ" {
                    currentOption.successBranch = field.Value
                } else if fieldName == "o_effect" {
                    currentOption.Effect = field.Value
                } else if fieldName == "o_fail" {
                    currentOption.failureBranch = field.Value
                } else if fieldName == "o_test" && conversation.LoadScriptFuncs {
                    branchCond, parseErr := govaluate.NewEvaluableExpressionWithFunctions(field.Value, handler.GetScriptFuncs())
                    if parseErr != nil {
                        return nil, parseErr
                    }
                    currentOption.branchCondition = branchCond
                }
            }
        }
        if currentOption.PlayerText != "" {
            conversationNode.Options = append(conversationNode.Options, currentOption)
        }
        allNodes[conversationNode.Name] = conversationNode
    }
    conversation.nodes = allNodes
    return conversation, nil
}
