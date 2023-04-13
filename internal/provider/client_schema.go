package provider

type Project struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	PlatformId     string `json:"platform_id"`
	RegionId       string `json:"region_id"`
	PgVersion      int64  `json:"pg_version"`
	StorePasswords bool   `json:"store_passwords"`
}

type Branch struct {
	Id           string `json:"id"`
	ProjectId    string `json:"project_id"`
	Name         string `json:"name"`
	Primary      bool   `json:"primary"`
	CurrentState string `json:"current_state"`
}

type Role struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	BranchId string `json:"branch_id"`
}

type Database struct {
	Id        int64  `json:"id"`
	BranchId  string `json:"branch_id"`
	Name      string `json:"name"`
	OwnerName string `json:"owner_name"`
}

type Endpoint struct {
	Id                    string  `json:"id"`
	Host                  string  `json:"host"`
	BranchId              string  `json:"branch_id"`
	ProjectId             string  `json:"project_id"`
	RegionId              string  `json:"region_id"`
	AutoscalingLimitMinCu float64 `json:"autoscaling_limit_min_cu"`
	AutoscalingLimitMaxCu float64 `json:"autoscaling_limit_max_cu"`
	CurrentState          string  `json:"current_state"`
}

type Operation struct {
	Id         string `json:"id"`
	Action     string `json:"action"`
	Status     string `json:"status"`
	EndpointId string `json:"endpoint_id"`
	BranchId   string `json:"branch_id"`
	ProjectId  string `json:"project_id"`
}

type ProjectOutput struct {
	Project Project `json:"project"`
}

type ProjectCreateInputProjectBranch struct {
	Name string `json:"name"`
}

type ProjectCreateInputProject struct {
	Name                  string                          `json:"name"`
	RegionId              string                          `json:"region_id"`
	PgVersion             int64                           `json:"pg_version"`
	StorePasswords        bool                            `json:"store_passwords"`
	Branch                ProjectCreateInputProjectBranch `json:"branch"`
	AutoscalingLimitMinCu float64                         `json:"autoscaling_limit_min_cu"`
	AutoscalingLimitMaxCu float64                         `json:"autoscaling_limit_max_cu"`
}

type ProjectCreateInput struct {
	Project ProjectCreateInputProject `json:"project"`
}

type ProjectCreateOutput struct {
	Project   Project    `json:"project"`
	Roles     []Role     `json:"roles"`
	Databases []Database `json:"databases"`
	Branch    Branch     `json:"branch"`
	Endpoints []Endpoint `json:"endpoints"`
}

type ProjectUpdateInputProject struct {
	Name string `json:"name"`
}

type ProjectUpdateInput struct {
	Project ProjectUpdateInputProject `json:"project"`
}

type BranchListOutput struct {
	Branches []Branch `json:"branches"`
}

type BranchOutput struct {
	Branch Branch `json:"branch"`
}

type BranchUpdateInputBranch struct {
	Name string `json:"name"`
}

type BranchUpdateInput struct {
	Branch BranchUpdateInputBranch `json:"branch"`
}

type EndpointListOutput struct {
	Endpoints []Endpoint `json:"endpoints"`
}

type EndpointOutput struct {
	Endpoint Endpoint `json:"endpoint"`
}

type EndpointUpdateInputEndpoint struct {
	AutoscalingLimitMinCu float64 `json:"autoscaling_limit_min_cu"`
	AutoscalingLimitMaxCu float64 `json:"autoscaling_limit_max_cu"`
}

type EndpointUpdateInput struct {
	Endpoint EndpointUpdateInputEndpoint `json:"endpoint"`
}

type OperationListOutput struct {
	Operations []Operation `json:"operations"`
}

type RoleOutput struct {
	Role Role `json:"role"`
}

type RoleCreateInputRole struct {
	Name string `json:"name"`
}

type RoleCreateInput struct {
	Role RoleCreateInputRole `json:"role"`
}

type DatabaseOutput struct {
	Database Database `json:"database"`
}

type DatabaseCreateInputDatabase struct {
	Name      string `json:"name"`
	OwnerName string `json:"owner_name"`
}

type DatabaseCreateInput struct {
	Database DatabaseCreateInputDatabase `json:"database"`
}

type DatabaseUpdateInputDatabase struct {
	Name      string `json:"name"`
	OwnerName string `json:"owner_name"`
}

type DatabaseUpdateInput struct {
	Database DatabaseUpdateInputDatabase `json:"database"`
}
