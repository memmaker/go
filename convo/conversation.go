package convo

import "github.com/Knetic/govaluate"

type ConversationPartner interface {
    Name() string
}
type ConversationNode struct {
    Name    string
    NpcText string
    Effects []string
    Options []conversationOption
}

func (n ConversationNode) IsEmpty() bool {
    return n.Name == "" && n.NpcText == "" && len(n.Options) == 0
}

func (n ConversationNode) WithoutOptionByIndex(index int) ConversationNode {
    if index >= len(n.Options) {
        return n
    }
    n.Options = append(n.Options[:index], n.Options[index+1:]...)
    return n
}

func (n ConversationNode) MoveOptionUp(index int) ConversationNode {
    if index <= 0 || index >= len(n.Options) {
        return n
    }
    n.Options[index], n.Options[index-1] = n.Options[index-1], n.Options[index]
    return n
}

func (n ConversationNode) MoveOptionDown(index int) ConversationNode {
    if index < 0 || index >= len(n.Options)-1 {
        return n
    }
    n.Options[index], n.Options[index+1] = n.Options[index+1], n.Options[index]
    return n
}

type OpeningBranch struct {
    Name            string
    BranchCondition *govaluate.EvaluableExpression
    BranchName      string
}

type DialogueHandler interface {
    FillTemplatedText(text string) string
    GetScriptFuncs() map[string]govaluate.ExpressionFunction
}

type Conversation struct {
    handler         DialogueHandler
    openingBranches []OpeningBranch
    nodes           map[string]ConversationNode
    Variables       map[string]interface{}

    prevNode        string
    LoadScriptFuncs bool
}

func NewConversation(handler DialogueHandler) *Conversation {
    return &Conversation{nodes: make(map[string]ConversationNode), handler: handler, LoadScriptFuncs: true}
}

func (c *Conversation) MergeVariables(params map[string]interface{}) {
    if c.Variables == nil {
        c.Variables = params
        return
    }
    for key, value := range params {
        c.Variables[key] = value
    }
}
func (c *Conversation) GetRootNode() string {
    for _, branch := range c.openingBranches {
        if branch.Name != "" { // named branches are used for NPC initiated conversations
            continue
        }
        evaluateResult, err := branch.BranchCondition.Evaluate(c.Variables)
        asBool := evaluateResult.(bool)
        if err == nil && asBool {
            return branch.BranchName
        }
    }
    return ""
}

func (c *Conversation) GetOpeningBranches() []OpeningBranch {
    return c.openingBranches
}

func (c *Conversation) GetNodeByName(node string) ConversationNode {
    return c.nodes[node]
}

func (c *Conversation) GetAllNodes() map[string]ConversationNode {
    return c.nodes
}

func (c *Conversation) RemoveOptionByIndex(nodeName string, index int) {
    c.nodes[nodeName] = c.nodes[nodeName].WithoutOptionByIndex(index)
}

func (c *Conversation) MoveOptionUp(nodeName string, index int) {
    c.nodes[nodeName] = c.nodes[nodeName].MoveOptionUp(index)
}

func (c *Conversation) MoveOptionDown(nodeName string, index int) {
    c.nodes[nodeName] = c.nodes[nodeName].MoveOptionDown(index)
}

func (c *Conversation) GetOpeningBranchByName(name string) OpeningBranch {
    for _, branch := range c.openingBranches {
        if branch.Name == name {
            return branch
        }
    }
    return OpeningBranch{}
}
