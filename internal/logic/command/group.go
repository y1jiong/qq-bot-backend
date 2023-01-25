package command

import (
	"context"
	"qq-bot-backend/internal/service"
)

func tryGroup(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "approval":
			catch = tryGroupApproval(ctx, next[2])
		case "bind":
			// /group bind <namespace>
			service.Group().BindNamespace(ctx, service.Bot().GetGroupId(ctx), next[2])
			catch = true
		}
	case endBranchRe.MatchString(cmd):
		switch endBranchRe.FindString(cmd) {
		case "query":
			// /group query
			service.Group().QueryGroup(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		case "unbind":
			// /group unbind
			service.Group().Unbind(ctx, service.Bot().GetGroupId(ctx))
			catch = true
		}
	}
	return
}

func tryGroupApproval(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "enable":
			catch = tryGroupApprovalEnable(ctx, next[2])
		case "disable":
			catch = tryGroupApprovalDisable(ctx, next[2])
		}
	}
	return
}

func tryGroupApprovalEnable(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "regexp":
			service.Group().AddApprovalProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		case "whitelist":
			service.Group().AddApprovalProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		case "blacklist":
			service.Group().AddApprovalProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		}
	case endBranchRe.MatchString(cmd):
		end := endBranchRe.FindString(cmd)
		switch end {
		case "mc":
			service.Group().AddApprovalProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		case "regexp":
			service.Group().AddApprovalProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		case "whitelist":
			service.Group().AddApprovalProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		case "blacklist":
			service.Group().AddApprovalProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		}
	}
	return
}

func tryGroupApprovalDisable(ctx context.Context, cmd string) (catch bool) {
	switch {
	case nextBranchRe.MatchString(cmd):
		next := nextBranchRe.FindStringSubmatch(cmd)
		switch next[1] {
		case "regexp":
			service.Group().RemoveApprovalProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		case "whitelist":
			service.Group().RemoveApprovalProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		case "blacklist":
			service.Group().RemoveApprovalProcess(ctx, service.Bot().GetGroupId(ctx), next[1], next[2])
			catch = true
		}
	case endBranchRe.MatchString(cmd):
		end := endBranchRe.FindString(cmd)
		switch end {
		case "mc":
			service.Group().RemoveApprovalProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		case "regexp":
			service.Group().RemoveApprovalProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		case "whitelist":
			service.Group().RemoveApprovalProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		case "blacklist":
			service.Group().RemoveApprovalProcess(ctx, service.Bot().GetGroupId(ctx), end)
			catch = true
		}
	}
	return
}
