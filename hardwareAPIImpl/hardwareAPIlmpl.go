package hardwareAPIImpl

import (
    "log"
	"time"
	"sync"
	"fmt"
	"github.com/anuchandy/coffeemaker/hardwareAPI"
)

// UserAction groups the methods representing various actions that
// user can perform on the coffee machine.
type UserAction interface {
  // Fill water to boiler.
  FillWater()
  // Press the brewing button.
  PressBrewButton()
  // Place pot in the warmer plate.
  PutPot()
  // Remove pot from the warmer plate.
  RemovePot()
  // Show the current status, state of coffee-maker.
  ShowState()
}

// HardwareAPIImpl implements HardwareAPI and UserAction interface.
type HardwareAPIImpl struct {
  // Indicates brewing is in progress or not.
  brewingInProgress bool
  // Timer to track brewing time.
  brewingTicker     *time.Ticker
  // Boiler heater's current state [BOILER_ON, BOILER_OFF].
  boilerState       hardwareAPI.BoilerState
  // Boiler's water status [BOILER_EMPTY, BOILER_NOT_EMPTY].
  boilerStatus      hardwareAPI.BoilerStatus
  // Brew button status [BREW_BUTTON_PUSHED, BREW_BUTTON_NOT_PUSHED].
  brewButtonStatus  hardwareAPI.BrewButtonStatus
  // Warmer plate's heating element state [WARMER_ON, WARMER_OFF].
  warmerPlateState  hardwareAPI.WarmerPlateState
  // Warmer plate's status [WARMER_EMPTY, POT_EMPTY, POT_NOT_EMPTY].
  warmerStatus      hardwareAPI.WarmerStatus
  // Pressure relieving valve state [VALVE_OPEN, VALVE_CLOSED].
  reliefValveState  hardwareAPI.ReliefValveState
  // Indicator light state [INDICATOR_ON, INDICATOR_OFF].
  indicatorState    hardwareAPI.IndicatorState
  // Controls concurrent access to brewingInProgress flag.
  brewingMutex      sync.Mutex
}

// Reset reset the state of various hardware components of coffee maker.
func (h *HardwareAPIImpl) Reset() {
  h.setBrewingInProgress(false)
  if h.brewingTicker != nil {
    h.brewingTicker.Stop()
	h.brewingTicker = nil
  }

  h.boilerState      = hardwareAPI.BOILER_OFF
  h.boilerStatus     = hardwareAPI.BOILER_EMPTY

  h.brewButtonStatus = hardwareAPI.BREW_BUTTON_NOT_PUSHED

  h.warmerPlateState = hardwareAPI.WARMER_OFF
  h.warmerStatus     = hardwareAPI.WARMER_EMPTY

  h.reliefValveState = hardwareAPI.VALVE_CLOSED
  h.indicatorState   = hardwareAPI.INDICATOR_OFF
}

// GetBoilerStatus satisfies HardwareAPI::GetBoilerStatus.
func (h *HardwareAPIImpl) GetBoilerStatus() hardwareAPI.BoilerStatus {
  return h.boilerStatus
}

// GetBrewButtonStatus satisfies HardwareAPI::GetBrewButtonStatus.
func (h *HardwareAPIImpl) GetBrewButtonStatus() hardwareAPI.BrewButtonStatus {
  currentState := h.brewButtonStatus
  h.brewButtonStatus = hardwareAPI.BREW_BUTTON_NOT_PUSHED
  return currentState
}

// GetWarmerPlateStatus satisfies HardwareAPI::GetWarmerPlateStatus.
func (h *HardwareAPIImpl) GetWarmerPlateStatus() hardwareAPI.WarmerStatus {
  return h.warmerStatus
}

// SetBoilerState satisfies HardwareAPI::SetBoilerState.
func (h *HardwareAPIImpl) SetBoilerState(boilerState hardwareAPI.BoilerState) {
  h.boilerState = boilerState

  if h.boilerState == hardwareAPI.BOILER_ON {
    go func () {
		if h.isBrewingInProgress() {
		  return;
		}

		h.setBrewingInProgress(true)

		h.brewingTicker = time.NewTicker(1 * time.Second)
		for i := 0; i < 5 && h.isBrewingInProgress(); i++ {
		  select {
		    case <-h.brewingTicker.C:
		  }

		  if h.reliefValveState == hardwareAPI.VALVE_OPEN {
		    i = 0
		  }
		}

		if h.isBrewingInProgress() {
		  h.warmerStatus = hardwareAPI.POT_NOT_EMPTY
		  h.boilerStatus = hardwareAPI.BOILER_EMPTY
		}

		h.setBrewingInProgress(false)
		h.brewingTicker.Stop()
		h.brewingTicker = nil
	}()
  } else {
    if h.brewingInProgress {
		h.setBrewingInProgress(false)
	}
  }
}

// SetIndicatorState satisfies HardwareAPI::SetIndicatorState.
func (h *HardwareAPIImpl) SetIndicatorState(indicatorState hardwareAPI.IndicatorState) {
  h.indicatorState = indicatorState
}

// SetReliefValveState satisfies HardwareAPI::SetReliefValveState.
func (h *HardwareAPIImpl) SetReliefValveState(reliefValveState hardwareAPI.ReliefValveState) {
  h.reliefValveState = reliefValveState
}

// SetWarmerPlateState satisfies HardwareAPI::SetWarmerPlateState.
func (h *HardwareAPIImpl) SetWarmerPlateState(warmerPlateState hardwareAPI.WarmerPlateState) {
  h.warmerPlateState = warmerPlateState
}

// FillWater satisfies UserAction::FillWater.
func (h *HardwareAPIImpl) FillWater() {
	h.boilerStatus = hardwareAPI.BOILER_NOT_EMPTY
}

// PressBrewButton satisfies UserAction::PressBrewButton.
func (h *HardwareAPIImpl) PressBrewButton() {
	if h.isBrewingInProgress() {
		log.Println("NOP: There is ALREADY a brewing in progress!")
		return
	}

	h.brewButtonStatus = hardwareAPI.BREW_BUTTON_PUSHED
}

// PutPot satisfies UserAction::PutPot.
func (h *HardwareAPIImpl) PutPot() {
	h.warmerStatus = hardwareAPI.POT_EMPTY
}

// RemovePot satisfies UserAction::RemovePot.
func (h *HardwareAPIImpl) RemovePot() {
	h.warmerStatus = hardwareAPI.WARMER_EMPTY
}

// SetWarmerPlateState satisfies HardwareAPI::SetWarmerPlateState.
func (h *HardwareAPIImpl) ShowState() {
	bolierWaterStatus := "[There is no water in the boiler]"
	if h.boilerStatus == hardwareAPI.BOILER_NOT_EMPTY {
		bolierWaterStatus = "[Boiler has water]"
	}

	warmerStatus := "[There is no pot in the warmer plate]"
	if h.warmerStatus == hardwareAPI.POT_EMPTY {
      warmerStatus = "[Warmer plate holds an empty pot]"
	}

	if h.warmerStatus == hardwareAPI.POT_NOT_EMPTY {
		warmerStatus = "[The pot has coffee in it!]"
	}

	if h.isBrewingInProgress() {
      fmt.Printf("%-20s:%-8s\n", "Brewing", "InProgress")
	} else  {
      fmt.Printf("%-20s:%-8s\n", "Brewing", "NO")
	}

	if h.boilerState == hardwareAPI.BOILER_ON {
		fmt.Printf("%-20s:%-8s%s\n", "Boiler", "ON", bolierWaterStatus)
	} else {
		fmt.Printf("%-20s:%-8s%s\n", "Boiler", "OFF", bolierWaterStatus)
	}

	if h.warmerPlateState == hardwareAPI.WARMER_ON {
		fmt.Printf("%-20s:%-8s%s\n", "Warmer Plate", "ON", warmerStatus)
	} else {
		fmt.Printf("%-20s:%-8s%s\n", "Warmer Plate", "OFF", warmerStatus)
	}

	if h.reliefValveState == hardwareAPI.VALVE_OPEN {
		fmt.Printf("%-20s:%-8s\n", "Relief Valve", "OPEN")
	} else {
		fmt.Printf("%-20s:%-8s\n", "Relief Valve", "CLOSED")
	}

	if h.indicatorState == hardwareAPI.INDICATOR_ON {
		fmt.Printf("%-20s:%-8s\n", "Indicator", "ON")
	} else {
		fmt.Printf("%-20s:%-8s\n", "Indicator", "OFF")
	}
}

// isBrewingInProgress returns the true if brewing is in progress, false otherwise.
func (h *HardwareAPIImpl) isBrewingInProgress() bool {
  h.brewingMutex.Lock()
  defer h.brewingMutex.Unlock()
  return h.brewingInProgress
}

// setBrewingInProgress sets the brewing status, a value true for s indicates brewing
// is in progress, false when brewing is done or not started.
func (h *HardwareAPIImpl) setBrewingInProgress(s bool) {
  h.brewingMutex.Lock()
  defer h.brewingMutex.Unlock()
  h.brewingInProgress = s
}
