package bridges

import "context"

type Bridge interface {
	Run(ctx context.Context)
}
