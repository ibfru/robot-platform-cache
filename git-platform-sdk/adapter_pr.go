package sdkadapter

import (
	"context"
	"strconv"

	"github.com/opensourceways/go-gitee/gitee"
)

func (c *ClientTarget) DeletePRComment(pr *PRParameter) error {
	var id int32
	if len(pr.CommentID) > 0 {
		i, e := strconv.Atoi(pr.CommentID)
		if e != nil {
			id = int32(i)
		}
	}

	_, err := c.ac.PullRequestsApi.DeleteV5ReposOwnerRepoPullsCommentsId(
		context.Background(), pr.Org, pr.Repo, id, nil)
	return formatErr(err, "delete comment of pr")
}

func (c *ClientTarget) AddPRComment(pr *PRParameter) error {
	opt := gitee.PullRequestCommentPostParam{Body: pr.Comment}
	_, _, err := c.ac.PullRequestsApi.PostV5ReposOwnerRepoPullsNumberComments(
		context.Background(), pr.Org, pr.Repo, pr.Number, opt)
	return formatErr(err, "create comment of pr")
}
