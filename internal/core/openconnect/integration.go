package openconnect

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/brentyates/squaregolf-connector/internal/config"
	"github.com/brentyates/squaregolf-connector/internal/core"
	"github.com/brentyates/squaregolf-connector/internal/core/simulator"
)

var (
	openConnectInstance *Integration
	openConnectOnce     sync.Once
)

type Integration struct {
	*simulator.Base
	stateManager   *core.StateManager
	launchMonitor  *core.LaunchMonitor
	shotNumber     int
	lastShotNumber int
	lastPlayerInfo *PlayerInfo
}

func GetInstance(stateManager *core.StateManager, launchMonitor *core.LaunchMonitor, host string, port int) *Integration {
	openConnectOnce.Do(func() {
		openConnectInstance = &Integration{
			stateManager:  stateManager,
			launchMonitor: launchMonitor,
		}
		openConnectInstance.Base = simulator.NewBase(openConnectInstance, host, port)
		openConnectInstance.registerStateListeners()
	})
	return openConnectInstance
}

func (o *Integration) Name() string {
	return "OpenConnect"
}

func (o *Integration) DefaultPort() int {
	return 922
}

func (o *Integration) GetStateManager() *core.StateManager {
	return o.stateManager
}

func (o *Integration) GetLaunchMonitor() *core.LaunchMonitor {
	return o.launchMonitor
}

func (o *Integration) SetStatus(status simulator.ConnectionStatus) {
	switch status {
	case simulator.StatusDisconnected:
		o.stateManager.SetOpenConnectStatus(core.OpenConnectStatusDisconnected)
	case simulator.StatusConnecting:
		o.stateManager.SetOpenConnectStatus(core.OpenConnectStatusConnecting)
	case simulator.StatusConnected:
		o.stateManager.SetOpenConnectStatus(core.OpenConnectStatusConnected)
	case simulator.StatusError:
		o.stateManager.SetOpenConnectStatus(core.OpenConnectStatusError)
	}
}

func (o *Integration) SetError(err error) {
	o.stateManager.SetOpenConnectError(err)
}

func (o *Integration) OnConnected() {
	// Send initial ready message
	initMessage := ShotData{
		DeviceID:   "SquareGolf",
		Units:      "Yards",
		APIversion: "1",
		ShotNumber: 0,
		ShotDataOptions: ShotOptions{
			ContainsBallData:          false,
			ContainsClubData:          false,
			LaunchMonitorIsReady:      true,
			LaunchMonitorBallDetected: false,
		},
		PlayerData: o.getPlayerData(),
	}
	if err := o.sendData(initMessage); err != nil {
		log.Printf("Error sending init message to OpenConnect: %v", err)
	}

	// Activate ball detection immediately on connect
	if err := o.launchMonitor.ActivateBallDetection(); err != nil {
		log.Printf("Failed to activate ball detection for OpenConnect: %v", err)
	}
}

func (o *Integration) OnDisconnected() {
}

func (o *Integration) ProcessMessage(rawMessage string) {
	var baseMsg Message
	if err := json.Unmarshal([]byte(rawMessage), &baseMsg); err != nil {
		log.Printf("Invalid JSON from OpenConnect client: %v", err)
		return
	}

	switch baseMsg.Message {
	case "GSPro ready", "OpenConnect ready":
		o.handleReadyMessage()
	case "GSPro Player Information", "OpenConnect Player Information":
		var playerInfo PlayerInfo
		if err := json.Unmarshal([]byte(rawMessage), &playerInfo); err != nil {
			log.Printf("Error parsing player info: %v", err)
			return
		}
		o.handlePlayerMessage(&playerInfo)
		o.handleReadyMessage()
	case "Ball Data received", "Club & Ball Data received", "Shot received successfully":
		log.Printf("Received shot confirmation from OpenConnect: %s", baseMsg.Message)
		if err := o.launchMonitor.ActivateBallDetection(); err != nil {
			log.Printf("Failed to re-arm ball detection after shot: %v", err)
		}
	default:
		log.Printf("Unknown OpenConnect message type: %s", baseMsg.Message)
	}
}

func (o *Integration) handleReadyMessage() {
	err := o.launchMonitor.ActivateBallDetection()
	if err != nil {
		log.Printf("Failed to activate ball detection: %v", err)
		return
	}
}

func (o *Integration) handlePlayerMessage(playerInfo *PlayerInfo) {
	o.lastPlayerInfo = playerInfo

	if clubName := playerInfo.Player.Club; clubName != "" {
		clubType := o.mapClubToInternal(clubName)
		if clubType != nil {
			log.Printf("OpenConnect selected club: %s (mapped to %v)", clubName, clubType)
			o.stateManager.SetClub(clubType)
		} else {
			log.Printf("Unmapped OpenConnect club: %s", clubName)
		}

		friendlyName := mapClubToFriendlyName(clubName)
		o.stateManager.SetClubName(&friendlyName)
	}

	if handed := playerInfo.Player.Handed; handed != "" {
		var handednessType core.HandednessType
		var handednessStr string
		if handed == "LH" {
			handednessType = core.LeftHanded
			handednessStr = "left"
			log.Printf("OpenConnect selected handedness: Left-handed")
		} else {
			handednessType = core.RightHanded
			handednessStr = "right"
			log.Printf("OpenConnect selected handedness: Right-handed")
		}
		o.stateManager.SetHandedness(&handednessType)

		// Persist handedness to config so it survives restarts
		config.GetInstance().SetHandedness(handednessStr)
	}
}

func (o *Integration) sendData(shotData ShotData) error {
	jsonData, err := json.Marshal(shotData)
	if err != nil {
		return err
	}
	return o.Base.SendMessage(jsonData)
}
