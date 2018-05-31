package gameuser

// Profile ...
type Profile struct {
	UserID       int // user id
	Name         string
	Gender       int
	Level        int
	Spine        string
	Intelligence int
	Intimacy     int
	Stamina      int
	Experience   int
	Coin		 int
	Diamond		 int

}

// NewProfile ...
func NewProfile(userID int) (profile *Profile) {
	profile = &Profile{}
	profile.UserID = userID
	//profile.Level = 1
	return
}
