package repositories

import (
	"context"

	"github.com/bestreads/Backend/internal/database"
	"github.com/bestreads/Backend/internal/middlewares"
	"gorm.io/gorm"
)

// wichtige informationen:
//  1. den namen der preload dinger *RICHTIG* zu schreiben (auch groß- und kleinschreibung)
//  2. gorm wird (sporadisch) andere sachen preloaden, keine ahnung wieso?
//  3. eigentlich wollte ich eine api machen, mit der man nur die metadaten von einem user lädt.
//     würde ein bisschen hübscher in json aussehen, das habe ich aber mal nicht gemacht
func GetDbPost(ctx *fasthttp.RequestCtx, uid uint, bid uint) ([]database.Post, error) {
	return gorm.G[database.Post](middlewares.DB(ctx)).
		Preload("User", nil).
		Preload("Book", nil).
		Where("user_id = ? AND book_id = ?", uid, bid).
		Find(ctx)
}

func CreateDbPost(ctx context.Context, post database.Post) error {
	return gorm.G[database.Post](middlewares.DB(ctx)).Create(ctx, &post)
}
