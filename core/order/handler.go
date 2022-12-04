package order

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ardanlabs/service/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"github.com/polldo/govod/api/web"
	"github.com/polldo/govod/api/weberr"
	"github.com/polldo/govod/core/cart"
	"github.com/polldo/govod/core/claims"
	"github.com/polldo/govod/core/course"
	"github.com/polldo/govod/database"
)

// Check if the user has already bought a course in the order?
func HandleCheckout(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		clm, err := claims.Get(ctx)
		if err != nil {
			return weberr.NotAuthorized(errors.New("user not authenticated"))
		}

		items, err := cart.FetchItems(ctx, db, clm.UserID)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		ord := Order{
			ID:        validate.GenerateID(),
			UserID:    clm.UserID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		var tot float64
		err = database.Transaction(db, func(tx sqlx.ExtContext) error {
			if err := Create(ctx, db, ord); err != nil {
				return err
			}

			for _, it := range items {
				c, err := course.Fetch(ctx, db, it.CourseID)
				if err != nil {
					return err
				}

				CreateItem(ctx, db, Item{
					OrderID:   ord.ID,
					CourseID:  c.ID,
					Price:     c.Price,
					CreatedAt: now,
				})

				tot += c.Price
			}

			return nil
		})

		return web.Respond(ctx, w, nil, http.StatusOK)
	}
}

// Remember to clean the cart after a successful payment.
func HandleWebhook(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return web.Respond(ctx, w, nil, http.StatusOK)
	}
}
