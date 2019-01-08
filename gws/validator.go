package gws

import "errors"

func validateSubscription(s Subscription) []error {
	var errs = make([]error, 0)
	if s.ID() == "" {
		errs = append(errs, errors.New("Subscription ID is empty"))
	}
	if s.Connection() == nil {
		errs = append(errs, errors.New("Subscription is not associated with a connection"))
	}
	if s.Query() == "" {
		errs = append(errs, errors.New("Subscription query is empty"))
	}
	return errs
}
