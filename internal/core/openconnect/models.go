package openconnect

// Models for OpenConnect integration
// These data structures represent the messages exchanged with OpenConnect-compatible simulators
// Extends the standard GSPro OpenConnect format with PlayerData for handedness

// Message represents the base message structure from a simulator
type Message struct {
	Message string `json:"Message"`
}

// PlayerInfo represents player information from a simulator
type PlayerInfo struct {
	Message string `json:"Message"`
	Player  Player `json:"Player"`
}

// Player represents player details from a simulator
type Player struct {
	Club   string `json:"Club"`
	Handed string `json:"Handed"`
}

// ShotData represents the shot data sent to a simulator
type ShotData struct {
	DeviceID        string      `json:"DeviceID"`
	Units           string      `json:"Units"`
	APIversion      string      `json:"APIversion"`
	ShotNumber      int         `json:"ShotNumber"`
	ShotDataOptions ShotOptions `json:"ShotDataOptions"`
	BallData        *BallData   `json:"BallData,omitempty"`
	ClubData        *ClubData   `json:"ClubData,omitempty"`
	PlayerData      *PlayerData `json:"PlayerData,omitempty"`
}

// ShotOptions represents shot data options
type ShotOptions struct {
	ContainsBallData          bool `json:"ContainsBallData"`
	ContainsClubData          bool `json:"ContainsClubData"`
	LaunchMonitorIsReady      bool `json:"LaunchMonitorIsReady,omitempty"`
	LaunchMonitorBallDetected bool `json:"LaunchMonitorBallDetected,omitempty"`
}

// BallData represents ball data sent to a simulator
type BallData struct {
	Speed     float64 `json:"Speed"`
	SpinAxis  float64 `json:"SpinAxis"`
	TotalSpin int16   `json:"TotalSpin"`
	BackSpin  int16   `json:"BackSpin"`
	SideSpin  int16   `json:"SideSpin"`
	HLA       float64 `json:"HLA"`
	VLA       float64 `json:"VLA"`
}

// ClubData represents club data sent to a simulator
type ClubData struct {
	Speed                float64 `json:"Speed"`
	AngleOfAttack        float64 `json:"AngleOfAttack"`
	FaceToTarget         float64 `json:"FaceToTarget"`
	Lie                  float64 `json:"Lie"`
	Loft                 float64 `json:"Loft"`
	Path                 float64 `json:"Path"`
	SpeedAtImpact        float64 `json:"SpeedAtImpact"`
	VerticalFaceImpact   float64 `json:"VerticalFaceImpact"`
	HorizontalFaceImpact float64 `json:"HorizontalFaceImpact"`
	ClosureRate          float64 `json:"ClosureRate"`
}

// PlayerData represents enriched player data included with every shot
// This extends the standard OpenConnect format so receiving software
// knows the player's handedness and current club
type PlayerData struct {
	Handed string `json:"Handed"`
	Club   string `json:"Club"`
}
