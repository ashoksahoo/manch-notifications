package callbacks

import (
	"fmt"
	"notification-service/pkg/subscribers"
)

func MileStoneSubscriberCB(subj, reply string, m *subscribers.MileStone) {
	fmt.Printf("Received a new milestone on subject %s! with Value %+v\n", subj, m)
	// process mile stone
	fmt.Printf("Processed a new milestone on subject %s! with Id %s\n", subj, m.ProfileId)
}
