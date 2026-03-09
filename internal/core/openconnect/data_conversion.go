package openconnect

import (
	"github.com/brentyates/squaregolf-connector/internal/core"
)

// convertToShotFormat converts internal shot data format to OpenConnect format
func (o *Integration) convertToShotFormat(ballMetrics core.BallMetrics, incrementShot bool) ShotData {
	if incrementShot {
		o.shotNumber++
		o.lastShotNumber = o.shotNumber
	}

	return ShotData{
		DeviceID:   "SquareGolf",
		Units:      "Yards",
		APIversion: "1",
		ShotNumber: o.lastShotNumber,
		ShotDataOptions: ShotOptions{
			ContainsBallData: true,
			ContainsClubData: false,
		},
		BallData: &BallData{
			Speed:     ballMetrics.BallSpeedMPS * 2.23694, // Convert m/s to mph
			SpinAxis:  ballMetrics.SpinAxis * -1,
			TotalSpin: ballMetrics.TotalspinRPM,
			BackSpin:  ballMetrics.BackspinRPM,
			SideSpin:  ballMetrics.SidespinRPM * -1,
			HLA:       ballMetrics.HorizontalAngle,
			VLA:       ballMetrics.VerticalAngle,
		},
		ClubData:   &ClubData{},
		PlayerData: o.getPlayerData(),
	}
}

// convertClubData converts internal club data format to OpenConnect format
func (o *Integration) convertClubData(clubMetrics core.ClubMetrics) *ClubData {
	return &ClubData{
		Speed:                0,
		AngleOfAttack:        clubMetrics.AttackAngle,
		FaceToTarget:         clubMetrics.FaceAngle,
		Lie:                  0,
		Loft:                 clubMetrics.DynamicLoftAngle,
		Path:                 clubMetrics.PathAngle,
		SpeedAtImpact:        0,
		VerticalFaceImpact:   0,
		HorizontalFaceImpact: 0,
		ClosureRate:          0,
	}
}

// getPlayerData reads current handedness and club from state for enriched shot payload
func (o *Integration) getPlayerData() *PlayerData {
	handed := "RH"
	if h := o.stateManager.GetHandedness(); h != nil && *h == core.LeftHanded {
		handed = "LH"
	}

	club := ""
	if clubName := o.stateManager.GetClubName(); clubName != nil {
		club = *clubName
	}

	return &PlayerData{
		Handed: handed,
		Club:   club,
	}
}

// mapClubToInternal maps simulator club name to internal ClubType
func (o *Integration) mapClubToInternal(clubName string) *core.ClubType {
	clubMap := map[string]core.ClubType{
		"DR": core.ClubDriver,
		"W2": core.ClubWood3,
		"W3": core.ClubWood3,
		"W4": core.ClubWood5,
		"W5": core.ClubWood5,
		"W6": core.ClubWood7,
		"W7": core.ClubWood7,
		"H2": core.ClubWood3,
		"H3": core.ClubWood3,
		"H4": core.ClubWood3,
		"H5": core.ClubWood3,
		"H6": core.ClubWood5,
		"H7": core.ClubIron4,
		"I1": core.ClubWood3,
		"I2": core.ClubWood3,
		"I3": core.ClubWood5,
		"I4": core.ClubIron4,
		"I5": core.ClubIron5,
		"I6": core.ClubIron6,
		"I7": core.ClubIron7,
		"I8": core.ClubIron8,
		"I9": core.ClubIron9,
		"PW": core.ClubPitchingWedge,
		"AW": core.ClubApproachWedge,
		"GW": core.ClubApproachWedge,
		"SW": core.ClubSandWedge,
		"LW": core.ClubSandWedge,
		"PT": core.ClubPutter,
	}

	if club, ok := clubMap[clubName]; ok {
		return &club
	}
	return nil
}

// mapClubToFriendlyName converts club codes to short readable names
func mapClubToFriendlyName(clubCode string) string {
	nameMap := map[string]string{
		"DR": "DR", "W2": "2W", "W3": "3W", "W4": "4W", "W5": "5W", "W6": "6W", "W7": "7W",
		"H2": "2H", "H3": "3H", "H4": "4H", "H5": "5H", "H6": "6H", "H7": "7H",
		"I1": "1I", "I2": "2I", "I3": "3I", "I4": "4I", "I5": "5I", "I6": "6I", "I7": "7I", "I8": "8I", "I9": "9I",
		"PW": "PW", "AW": "AW", "GW": "GW", "SW": "SW", "LW": "LW",
		"PT": "PUTT",
	}

	if name, ok := nameMap[clubCode]; ok {
		return name
	}
	return clubCode
}
