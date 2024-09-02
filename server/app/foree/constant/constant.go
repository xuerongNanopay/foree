package foree_constant

const DefaultNBPIdPrefix = "NP"

const (
	DefaultFeeGroup              string = "FOREE_PERSONAL_FEE"
	DefaultRoleGroup             string = string(RoleGroupPersonal)
	DefaultTransactionLimitGroup string = string(TLPersonal1k)
)

type RoleGroup string

const (
	RoleGroupPersonal RoleGroup = "FOREE_PERSONAL"
	RoleGroupBO       RoleGroup = "FOREE_BO"
	RoleGroupAdmin    RoleGroup = "FOREE_ADMIN"
)
