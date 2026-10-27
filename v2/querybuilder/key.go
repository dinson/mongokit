package querybuilder

type KeyMongoDB string

func (k KeyMongoDB) String() string {
	return string(k)
}
