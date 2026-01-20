package steam

// Port describes Steam data access used by the service layer.
type Port interface {
	GetGamePages([]uint) ([]Page, error)
}
