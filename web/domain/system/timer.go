//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mock.go -package=$GOPACKAGE

package system

import "time"

type ITimer interface {
	Now() time.Time
}
