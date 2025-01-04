package convo

import "github.com/Knetic/govaluate"

type conversationOption struct {
    displayCondition *govaluate.EvaluableExpression
    PlayerText       string
    branchCondition  *govaluate.EvaluableExpression
    successBranch    string // will default to the current node if not set
    failureBranch    string
    Effect           string
}

func (o *conversationOption) canDisplay(params map[string]interface{}) bool {
    if o.displayCondition == nil {
        return true
    }
    evaluateResult, err := o.displayCondition.Evaluate(params)
    if err != nil {
        panic(err)
    }
    asBool := evaluateResult.(bool)
    return err == nil && asBool
}

func (o *conversationOption) getFollowupBranch(params map[string]interface{}) string {
    if o.branchCondition == nil {
        return o.successBranch
    }
    evaluateResult, err := o.branchCondition.Evaluate(params)
    asBool := evaluateResult.(bool)
    if err == nil && asBool {
        return o.successBranch
    }
    return o.failureBranch
}

func (o *conversationOption) GetAllPossibleBranches() []string {
    if o.branchCondition == nil {
        return []string{o.successBranch}
    }
    return []string{o.successBranch, o.failureBranch}
}

func (o *conversationOption) GetDisplayCondition() *govaluate.EvaluableExpression {
    return o.displayCondition
}

func (o *conversationOption) GetBranchCondition() *govaluate.EvaluableExpression {
    return o.branchCondition
}

func (o *conversationOption) GetSuccessBranch() string {
    return o.successBranch
}

func (o *conversationOption) GetFailureBranch() string {
    return o.failureBranch
}

func (o *conversationOption) GetGotoBranch() string {
    return o.successBranch
}
