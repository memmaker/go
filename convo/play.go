package convo

import (
    "github.com/Knetic/govaluate"
    "github.com/memmaker/go/fxtools"
)

type Option struct {
    PlayerText     string
    Effect         string
    FollowUpBranch string
}
type ConversationFlow uint8

const (
    ConversationContinue ConversationFlow = iota
    ConversationEndNormal
    ConversationEndInstantlyWithChatter
)

type ConversationState struct {
    NodeName      string
    NPCText       string
    PlayerOptions []Option
    Effects       []string
    Flow          ConversationFlow
}

func (c *Conversation) OpenDialogueNode(currentNodeName string) ConversationState {
    endConversation := false
    instantEndWithChatter := false
    currentNode := c.GetNodeByName(currentNodeName)

    nodeText := c.handler.FillTemplatedText(currentNode.NpcText)

    var nodeOptions []Option
    var nodeEffects []string
    for _, effect := range currentNode.Effects {
        if fxtools.LooksLikeAFunction(effect) {
            name, args := fxtools.GetNameAndArgs(effect)
            switch name { // these effects are also only possible here, because they directly influence the conversation flow
            case "GotoNode":
                nodeName := args.Get(0)
                nextNode := c.GetNodeByName(nodeName)
                if !nextNode.IsEmpty() {
                    currentNode = nextNode
                }
            default: // parse as generic expression and effect
                if c.LoadScriptFuncs {
                    expr, parseErr := govaluate.NewEvaluableExpressionWithFunctions(effect, c.handler.GetScriptFuncs())
                    if parseErr != nil {
                        panic(parseErr)
                    }
                    _, evalErr := expr.Evaluate(c.Variables)
                    if evalErr != nil {
                        panic(evalErr)
                    }
                } else {
                    nodeEffects = append(nodeEffects, effect)
                }

            }
            continue
        }
        switch effect {
        case "EndConversation":
            endConversation = true
        case "EndWithChatter":
            instantEndWithChatter = true
        case "ReturnToPreviousNode":
            if c.prevNode != "" {
                prevNode := c.GetNodeByName(c.prevNode)
                if !prevNode.IsEmpty() {
                    currentNode = prevNode
                }
            }
        default:
            nodeEffects = append(nodeEffects, effect)
        }
    }

    for _, o := range currentNode.Options {
        option := o
        if !c.LoadScriptFuncs || option.canDisplay(c.Variables) {
            nodeOptions = append(nodeOptions, Option{
                PlayerText:     c.handler.FillTemplatedText(option.PlayerText),
                Effect:         option.Effect,
                FollowUpBranch: option.getFollowupBranch(c.Variables),
            })
        }
    }
    flow := ConversationContinue
    if endConversation {
        flow = ConversationEndNormal
    } else if instantEndWithChatter {
        flow = ConversationEndInstantlyWithChatter
    }

    c.prevNode = currentNodeName

    return ConversationState{
        NodeName:      currentNodeName,
        NPCText:       nodeText,
        PlayerOptions: nodeOptions,
        Effects:       nodeEffects,
        Flow:          flow,
    }
}
