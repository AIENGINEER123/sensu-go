package handlers

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/sensu/sensu-go/backend/apid/actions"
	"github.com/sensu/sensu-go/backend/store"
	storev2 "github.com/sensu/sensu-go/backend/store/v2"
)

func (h Handlers) DeleteV3Resource(r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	name, err := url.PathUnescape(params["id"])
	if err != nil {
		return nil, actions.NewError(actions.InvalidArgument, err)
	}

	ctx := r.Context()
	namespace := store.NewNamespaceFromContext(ctx)
	storeName := h.V3Resource.StoreName()

	req := storev2.NewResourceRequest(ctx, namespace, name, storeName)
	if err := h.StoreV2.Delete(req); err != nil {
		switch err := err.(type) {
		case *store.ErrNotFound:
			return nil, actions.NewErrorf(actions.NotFound)
		case *store.ErrNotValid:
			return nil, actions.NewError(actions.InvalidArgument, err)
		default:
			return nil, actions.NewError(actions.InternalErr, err)
		}
	}
	return nil, nil
}
