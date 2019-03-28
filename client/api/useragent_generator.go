package api

import (
	"fmt"
	"giveaway/instagram"
)

type AppUserAgentGenerator struct {

}

func (u AppUserAgentGenerator) Get() string {
	return fmt.Sprintf("Instagram %s Android (22/5.1.1; 320dpi; 720x1280; samsung; SM-J320H; j3x3g; sc8830; en_US; %s)",
		instagram.Version,
		instagram.Constants[instagram.Version].VersionIncremental,
	)
}
