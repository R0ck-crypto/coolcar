package id

type AccountID string
type TripID string
type IdentityID string
type CarID string
type BlobID string

func (a AccountID) String() string {
	return string(a)
}

func (t TripID) String() string {
	return string(t)
}

func (i IdentityID) String() string {
	return string(i)
}

func (c CarID) String() string {
	return string(c)
}

func (b BlobID) String() string {
	return string(b)
}
