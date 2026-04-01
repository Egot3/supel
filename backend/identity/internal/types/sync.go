package types

type SyncData struct {
	UUID string `json:"uuid"`
	Role UserRole `json:"role"`
}