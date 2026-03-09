package openconnect

import (
	"log"

	"github.com/brentyates/squaregolf-connector/internal/core"
)

func (o *Integration) registerStateListeners() {
	o.stateManager.RegisterBallReadyCallback(o.onBallReadyChanged)
	o.stateManager.RegisterLastBallMetricsCallback(o.onLastBallMetricsChanged)
	o.stateManager.RegisterLastClubMetricsCallback(o.onLastClubMetricsChanged)
}

func (o *Integration) onBallReadyChanged(oldValue, newValue bool) {
	if oldValue == newValue {
		return
	}

	if !o.Base.Connected || o.Base.Socket == nil {
		return
	}

	emptyShotData := ShotData{
		DeviceID:   "SquareGolf",
		Units:      "Yards",
		APIversion: "1",
		ShotNumber: o.lastShotNumber,
		ShotDataOptions: ShotOptions{
			ContainsBallData:          false,
			ContainsClubData:          false,
			LaunchMonitorIsReady:      newValue,
			LaunchMonitorBallDetected: newValue,
		},
		PlayerData: o.getPlayerData(),
	}

	if err := o.sendData(emptyShotData); err != nil {
		log.Printf("Error sending empty shot data to OpenConnect: %v", err)
	}
}

func (o *Integration) onLastBallMetricsChanged(oldValue, newValue *core.BallMetrics) {
	if oldValue == newValue {
		return
	}

	if !o.Base.Connected || o.Base.Socket == nil {
		return
	}

	if newValue == nil {
		return
	}

	shotData := o.convertToShotFormat(*newValue, true)
	if err := o.sendData(shotData); err != nil {
		log.Printf("Error sending shot data to OpenConnect: %v", err)
	}
}

func (o *Integration) onLastClubMetricsChanged(oldValue, newValue *core.ClubMetrics) {
	if oldValue == newValue {
		return
	}

	if !o.Base.Connected || o.Base.Socket == nil {
		return
	}

	if newValue == nil {
		zeroedClubData := &ClubData{
			Speed:                0,
			AngleOfAttack:        0,
			FaceToTarget:         0,
			Lie:                  0,
			Loft:                 0,
			Path:                 0,
			SpeedAtImpact:        0,
			VerticalFaceImpact:   0,
			HorizontalFaceImpact: 0,
			ClosureRate:          0,
		}

		shotData := o.convertToShotFormat(core.BallMetrics{}, false)
		shotData.ShotDataOptions.ContainsBallData = false
		shotData.ShotDataOptions.ContainsClubData = true
		shotData.ClubData = zeroedClubData
		if err := o.sendData(shotData); err != nil {
			log.Printf("Error sending zeroed club data to OpenConnect: %v", err)
		}
		return
	}

	shotData := o.convertToShotFormat(core.BallMetrics{}, false)
	shotData.ShotDataOptions.ContainsBallData = false
	shotData.ShotDataOptions.ContainsClubData = true
	shotData.ClubData = o.convertClubData(*newValue)
	if err := o.sendData(shotData); err != nil {
		log.Printf("Error sending club data to OpenConnect: %v", err)
	}
}
