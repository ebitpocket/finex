package queries

import "github.com/nusa-exchange/finex/controllers/helpers"

type DepthQuery struct {
	Limit int64 `query:"limit" validate:"uint"`
}

func (t DepthQuery) Messages() map[string]string {
	return helpers.VaildateMessage("public.market_depth")
}

func (t DepthQuery) Translates() map[string]string {
	return helpers.VaildateTranslateFields()
}
