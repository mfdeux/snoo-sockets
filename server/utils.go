package server

import "github.com/labstack/echo"

func getSubscriptionQuery(c echo.Context, subscription string) string {
	switch subscription {
	case SubscriptionThingsAll:
		return ""
	case SubscriptionSubmissionsAll:
		return "all"
	case SubscriptionSubmissionsSubreddit:
		return c.Param("subreddit")
	case SubscriptionSubmissionsUser:
		return c.Param("username")
	case SubscriptionCommentsAll:
		return "all"
	case SubscriptionCommentsSubreddit:
		return c.Param("subreddit")
	case SubscriptionCommentsSubmission:
		return c.Param("submissionId")
	case SubscriptionCommentsUser:
		return c.Param("username")
	default:
		return ""
	}
}

func shouldSendMessage(c *Client, m []byte) bool {
	return true
}
