package payment

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
	"github.com/polldo/govod/core/order"
	"github.com/polldo/govod/database"
)

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
		ord := order.Order{
			ID:        validate.GenerateID(),
			UserID:    clm.UserID,
			CreatedAt: now,
		}

		var tot float64
		err = database.Transaction(db, func(tx sqlx.ExtContext) error {
			if err := order.Create(ctx, db, ord); err != nil {
				return err
			}

			for _, it := range items {
				c, err := course.Fetch(ctx, db, it.CourseID)
				if err != nil {
					return err
				}

				order.CreateItem(ctx, db, order.Item{
					OrderID:   ord.ID,
					CourseID:  c.ID,
					Price:     c.Price,
					CreatedAt: now,
				})

				tot += c.Price
			}

			pay := Payment{
				ID:         validate.GenerateID(),
				OrderID:    ord.ID,
				ProviderID: "",
				Amount:     tot,
				CreatedAt:  now,
				UpdatedAt:  now,
			}

			if err := Create(ctx, db, pay); err != nil {
				return err
			}

			return nil
		})

		return web.Respond(ctx, w, nil, http.StatusOK)
	}
}

func HandleWebhook(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return web.Respond(ctx, w, nil, http.StatusOK)
	}
}
