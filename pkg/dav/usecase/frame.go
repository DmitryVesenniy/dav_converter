package usecase

type (
	// FrameHeader заголовок каждого кадра
	FrameHeader struct {
		CountKadr uint32
		TimeTM    uint32
		SizeKadr  uint32
		Top       int16
		Left      int16
		Width     int16
		Height    int16
		CountAll  uint32
		TFront    uint32
		TSpad     uint32
	}

	// TagFrameIDX одна запись из таблицы индексов dav
	TagFrameIDX struct {
		TimeTM     uint32
		KadrIndex  uint32
		KadrOffset int64
		SizeData   uint32
	}
)
