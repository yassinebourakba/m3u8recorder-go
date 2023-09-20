package modules

type RoomDossier struct {
	HlsSource  string `json:"hls_source"`
	RoomStatus string `json:"room_status"`
}

type TopTenResponse struct {
	Top []struct {
		RoomSlug string `json:"room_slug"`
	} `json:"top" validate:"required,dive"`
}
