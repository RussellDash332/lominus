package api

import (
	"time"
)

type Announcement struct {
	Title       string // Announcement title
	Description string // Announcement description
	LastUpdated int64  // Announcement updated time
}

const ANNOUNCEMENT_URL_ENDPOINT = "https://luminus.nus.edu.sg/v2/api/announcement/NonArchived/%s"

func (req AnnouncementRequest) GetAnnouncements() ([]Announcement, error) {
	var announcements []Announcement

	rawResponse := RawResponse{}
	err := req.Request.GetRawResponse(&rawResponse)
	if err != nil {
		return announcements, err
	}

	for _, content := range rawResponse.Data {
		if _, exists := content["access"]; exists {
			title := content["title"].(string)
			description := content["description"].(string)
			lastUpdatedTime, err := time.Parse(time.RFC3339, content["lastUpdatedDate"].(string))
			if err != nil {
				return announcements, err
			}
			lastUpdated := lastUpdatedTime.Unix()

			announcement := Announcement{
				title,
				description,
				lastUpdated,
			}
			announcements = append(announcements, announcement)
		}
	}
	return announcements, nil
}
