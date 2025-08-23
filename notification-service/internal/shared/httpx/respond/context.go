package respond

// Ctx context for gin.Context low coupling respond
type Ctx interface {
	JSON(int, any)
	GetString(string) string
}
