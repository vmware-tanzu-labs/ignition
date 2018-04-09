package cloudfoundry

// API is a Cloud Controller API
type API interface {
	OrganizationCreator
	OrganizationQuerier
	SpaceCreator
	RoleGrantor
	QuotaQuerier
}
