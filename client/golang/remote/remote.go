package remote

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"myst/client/golang/core"
	"myst/client/golang/models"
	"myst/server/apis/generated"
)

var _ core.Remote = (*Remote)(nil)

type Remote struct {
	address     string
	client      *generated.ClientWithResponses
	bearerToken string
	user        *models.User
	username    string
	password    string
}

func New(opts ...Option) (*Remote, error) {
	r := &Remote{}

	for _, opt := range opts {
		if opt != nil {
			opt(r)
		}
	}

	var err error
	r.client, err = generated.NewClientWithResponses(
		r.address,
		generated.WithRequestEditorFn(r.authenticationMiddleware()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http client")
	}

	return r, nil
}

type Option func(*Remote)

func WithAddress(address string) Option {
	return func(r *Remote) {
		r.address = address
	}
}

func (r *Remote) Address() string {
	return r.address
}

func (r *Remote) authenticationMiddleware() generated.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+r.bearerToken)

		return nil
	}
}
