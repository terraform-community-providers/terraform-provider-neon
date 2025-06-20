package provider

type Project struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	PlatformId     string `json:"platform_id"`
	RegionId       string `json:"region_id"`
	OrgId          string `json:"org_id,omitempty"`
	PgVersion      int64  `json:"pg_version"`
	StorePasswords bool   `json:"store_passwords"`
}

type Branch struct {
	Id           string  `json:"id"`
	ProjectId    string  `json:"project_id"`
	ParentId     *string `json:"parent_id"`
	Name         string  `json:"name"`
	Default      bool    `json:"default"`
	Protected    bool    `json:"protected"`
	CurrentState string  `json:"current_state"`
}

type Role struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	BranchId  string `json:"branch_id"`
	Protected bool   `json:"protected"`
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
	ComputeProvisioner    string  `json:"provisioner"`
	SuspendTimeoutSeconds int64   `json:"suspend_timeout_seconds"`
	Type                  string  `json:"type"`
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

type ProjectCreateInputProjectDefaultEndpointSettings struct {
	AutoscalingLimitMinCu float64 `json:"autoscaling_limit_min_cu"`
	AutoscalingLimitMaxCu float64 `json:"autoscaling_limit_max_cu"`
	SuspendTimeoutSeconds int64   `json:"suspend_timeout_seconds"`
}

type ProjectCreateInputProject struct {
	Name                    string                                           `json:"name"`
	RegionId                string                                           `json:"region_id"`
	OrgId                   *string                                          `json:"org_id,omitempty"`
	PgVersion               int64                                            `json:"pg_version"`
	StorePasswords          bool                                             `json:"store_passwords"`
	Branch                  ProjectCreateInputProjectBranch                  `json:"branch"`
	DefaultEndpointSettings ProjectCreateInputProjectDefaultEndpointSettings `json:"default_endpoint_settings"`
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

type BranchCreateInputBranch struct {
	Name      string `json:"name"`
	ParentId  string `json:"parent_id,omitempty"`
	Protected *bool  `json:"protected,omitempty"`
}

type BranchCreateInput struct {
	Branch BranchCreateInputBranch `json:"branch"`
}

type BranchUpdateInputBranch struct {
	Name      *string `json:"name,omitempty"`
	Protected *bool   `json:"protected,omitempty"`
}

type BranchUpdateInput struct {
	Branch BranchUpdateInputBranch `json:"branch"`
}

type BranchEndpointListOutput struct {
	Endpoints []Endpoint `json:"endpoints"`
}

type EndpointOutput struct {
	Endpoint Endpoint `json:"endpoint"`
}

type EndpointCreateInputEndpoint struct {
	BranchId              string  `json:"branch_id"`
	Type                  string  `json:"type"`
	AutoscalingLimitMinCu float64 `json:"autoscaling_limit_min_cu"`
	AutoscalingLimitMaxCu float64 `json:"autoscaling_limit_max_cu"`
	SuspendTimeoutSeconds int64   `json:"suspend_timeout_seconds"`
}

type EndpointCreateInput struct {
	Endpoint EndpointCreateInputEndpoint `json:"endpoint"`
}

type EndpointUpdateInputEndpoint struct {
	AutoscalingLimitMinCu float64 `json:"autoscaling_limit_min_cu"`
	AutoscalingLimitMaxCu float64 `json:"autoscaling_limit_max_cu"`
	SuspendTimeoutSeconds int64   `json:"suspend_timeout_seconds"`
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

type RolePasswordOutput struct {
	Password string `json:"password"`
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
