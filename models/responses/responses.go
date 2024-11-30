package responses

type Quote struct {
	ID    int    `gorm:"type:BIGINT NOT NULL PRIMARY KEY"`
	Quote string `gorm:"type:VARCHAR NOT NULL"`
}

type Error struct {
	Code    int
	Message string
}
