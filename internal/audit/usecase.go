package audit

import "context"

type Usecase interface {
	Create(ctx context.Context, who, what string)
}
