package curve

type Curve interface{}
type Point interface{}

type HashToPoint interface {
	Hash([]byte) Point
}
