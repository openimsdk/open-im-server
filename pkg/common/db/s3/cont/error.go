package cont

import (
	"fmt"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/s3"
)

type HashAlreadyExistsError struct {
	Object *s3.ObjectInfo
}

func (e *HashAlreadyExistsError) Error() string {
	return fmt.Sprintf("hash already exists: %s", e.Object.Key)
}
