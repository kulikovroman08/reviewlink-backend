package controller

import (
	"github.com/kulikovroman08/reviewlink-backend/internal/service"
)

type Application struct {
	UserService        service.UserService
	PlaceService       service.PlaceService
	publicNewsService  service.PublicNewsService
	ReviewService      service.ReviewService
	TokenService       service.TokenService
	AdminService       service.AdminService
	LeaderboardService service.LeaderboardService
	BonusService       service.BonusService
	ReviewReplyService service.ReviewReplyService
	ReviewVoteService  service.ReviewVoteService
}

func NewApplication(
	user service.UserService,
	place service.PlaceService,
	publicNewsService service.PublicNewsService,
	review service.ReviewService,
	token service.TokenService,
	admin service.AdminService,
	leaderboard service.LeaderboardService,
	bonus service.BonusService,
	reviewReply service.ReviewReplyService,
	reviewVote service.ReviewVoteService,
) *Application {
	return &Application{
		UserService:        user,
		PlaceService:       place,
		publicNewsService:  publicNewsService,
		ReviewService:      review,
		TokenService:       token,
		AdminService:       admin,
		LeaderboardService: leaderboard,
		BonusService:       bonus,
		ReviewReplyService: reviewReply,
		ReviewVoteService:  reviewVote,
	}
}
