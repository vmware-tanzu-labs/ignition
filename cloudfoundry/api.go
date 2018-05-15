package cloudfoundry

// API is a Cloud Controller API
// counterfeiter ./cloudfoundry API
type API interface {
	OrganizationCreator
	OrganizationQuerier
	SpaceCreator
	RoleGrantor
	QuotaQuerier
	ISOSegmentQuerier
}
