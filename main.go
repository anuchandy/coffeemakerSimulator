package main
import (
	"github.com/anuchandy/coffeemakerSimulator/hardwareAPIImpl"
	"github.com/anuchandy/coffeemaker"
	"fmt"
)

func main() {
	// Creates a mock coffee-machine hardware
	var cmHardware *hardwareAPIImpl.HardwareAPIImpl = &hardwareAPIImpl.HardwareAPIImpl{}
	// Initializes the hardware
	cmHardware.Reset()

	// Switch-on the coffee-maker.
	coffeemaker.SwitchOn(cmHardware)

	var ui hardwareAPIImpl.UserAction = cmHardware
	for ;; {
	  var action int
	  fmt.Print("\nAction [1: Fill_Water 2: Place_Pot 3: Remove_Pot 4: Press_BrewButton 5: Show_Status 6: Exit] : ")
	  fmt.Scanf("%d", &action)

	  if action == 1 {
		  ui.FillWater()
	  }	else if action == 2 {
		  ui.PutPot()
	  } else if action == 3 {
		  ui.RemovePot()
	  } else if action == 4 {
		  ui.PressBrewButton()
	  } else if action == 5 {
		  ui.ShowState()
	  } else if action == 6 {
		break
	  } else  {
	    fmt.Println("Unknown action")
	  }
	}

	// Switch-off the coffee-maker.
	coffeemaker.SwitchOff()
}