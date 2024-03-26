package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kusionstack.io/kusion-module-framework/pkg/module"
	"kusionstack.io/kusion-module-framework/pkg/server"
	apiv1 "kusionstack.io/kusion/pkg/apis/core/v1"
	"kusionstack.io/kusion/pkg/log"
	"kusionstack.io/kusion/pkg/workspace"
)

const (
	CloudDBType = "cloud"
	LocalDBType = "local"
)

const (
	dbEngine         = "mysql"
	dbResSuffix      = "-mysql"
	dbHostAddressEnv = "KUSION_DB_HOST"
	dbUsernameEnv    = "KUSION_DB_USERNAME"
	dbPasswordEnv    = "KUSION_DB_PASSWORD"
)

var (
	ErrEmptyInstanceTypeForCloudDB = errors.New("empty instance type for cloud managed mysql instance")
	ErrEmptyCloudProviderType      = errors.New("empty cloud provider type in mysql module config")
)

var (
	localDeploymentSuffix = "-db-local-deployment"
	localSecretSuffix     = "-db-local-secret"
	localPVCSuffix        = "-db-local-pvc"
	localServiceSuffix    = "-db-local-service"
)

var (
	defaultUsername       string   = "root"
	defaultCategory       string   = "Basic"
	defaultSecurityIPs    []string = []string{"0.0.0.0/0"}
	defaultPrivateRouting bool     = true
	defaultSize           int      = 10
)

var defaultRandomProviderCfg = apiv1.ProviderConfig{
	Source:  "hashicorp/random",
	Version: "3.6.0",
}

var randomPassword = "random_password"

// MySQL describes the attributes to locally deploy or create a cloud provider
// managed MySQL database instance for the workload.
type MySQL struct {
	// The deployment mode of the MySQL database.
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// The MySQL database version to use.
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	// The type of the MySQL instance.
	InstanceType string `json:"instanceType,omitempty" yaml:"instanceType,omitempty"`
	// The allocated storage size of the MySQL instance.
	Size int `json:"size,omitempty" yaml:"size,omitempty"`
	// The edition of the MySQL instance provided by the cloud vendor.
	Category string `json:"category,omitempty" yaml:"category,omitempty"`
	// The operation account for the MySQL database.
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	// The list of IP addresses allowed to access the MySQL instance provided by the cloud vendor.
	SecurityIPs []string `json:"securityIPs,omitempty" yaml:"securityIPs,omitempty"`
	// The virtual subnet ID associated with the VPC that the cloud MySQL instance will be created in.
	SubnetID string `json:"subnetID,omitempty" yaml:"subnetID,omitempty"`
	// Whether the host address of the cloud MySQL instance for the workload to connect with is via
	// public network or private network of the cloud vendor.
	PrivateRouting bool `json:"privateRouting,omitempty" yaml:"privateRouting,omitempty"`
	// The specified name of the MySQL database instance.
	DatabaseName string `json:"databaseName,omitempty" yaml:"databaseName,omitempty"`
}

func (mysql *MySQL) Generate(_ context.Context, request *module.GeneratorRequest) (*module.GeneratorResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Debugf("failed to generate mysql module: %v", r)
		}
	}()

	// MySQL does not exist in AppConfiguration configs.
	if request.DevModuleConfig == nil {
		log.Info("MySQL does not exist in AppConfig config")

		return nil, nil
	}

	// Get the complete configs of the MySQL instance.
	err := mysql.GetCompleteConfig(request.DevModuleConfig, request.PlatformModuleConfig)
	if err != nil {
		return nil, err
	}

	// Set the database name.
	if mysql.DatabaseName == "" {
		mysql.DatabaseName = GenerateDefaultMySQLName(request.Project, request.Stack, request.App)
	}

	// Generate the MySQL intance resources based on the type and the cloud provider config.
	var resources []apiv1.Resource
	var patchers []apiv1.Patcher
	switch strings.ToLower(mysql.Type) {
	case LocalDBType:
		resources, patchers, err = mysql.GenerateLocalResources(request)
	case CloudDBType:
		providerType, err := GetCloudProviderType(request.PlatformModuleConfig)
		if err != nil {
			return nil, err
		}

		switch strings.ToLower(providerType) {
		case "aws":
			resources, patchers, err = mysql.GenerateAWSResources(request)
		case "alicloud":
			resources, patchers, err = mysql.GenerateAlicloudResources(request)
		default:
			return nil, fmt.Errorf("unsupported cloud provider type: %s", providerType)
		}
	default:
		return nil, fmt.Errorf("unsupported mysql type: %s", mysql.Type)
	}

	if err != nil {
		return nil, err
	}

	return &module.GeneratorResponse{
		Resources: resources,
		Patchers:  patchers,
	}, nil
}

// GetCompleteConfig combines the configs in devModuleConfig and platformModuleConfig to form a complete
// configuration for the MySQL instance.
func (mysql *MySQL) GetCompleteConfig(devConfig apiv1.Accessory, platformConfig apiv1.GenericConfig) error {
	// Set the default values for MySQL instance if platformConfig not exists.
	if platformConfig == nil {
		mysql.Username = defaultUsername
		mysql.Category = defaultCategory
		mysql.SecurityIPs = defaultSecurityIPs
		mysql.PrivateRouting = defaultPrivateRouting
		mysql.Size = defaultSize
	}

	// Get the type and version of the MySQL instance in devConfig.
	if mysqlType, ok := devConfig["type"]; ok {
		mysql.Type = mysqlType.(string)
	}
	if mysqlVersion, ok := devConfig["version"]; ok {
		mysql.Version = mysqlVersion.(string)
	}

	// Get the other configs of the MySQL instance in platformConfig,
	// and use the default values if some of them don't exist.
	if username, ok := platformConfig["username"]; ok {
		mysql.Username = username.(string)
	} else {
		mysql.Username = defaultUsername
	}

	if category, ok := platformConfig["category"]; ok {
		mysql.Category = category.(string)
	} else {
		mysql.Category = defaultCategory
	}

	if securityIPs, ok := platformConfig["securityIPs"]; ok {
		mysql.SecurityIPs = securityIPs.([]string)
	} else {
		mysql.SecurityIPs = defaultSecurityIPs
	}

	if privateRouting, ok := platformConfig["privateRouting"]; ok {
		mysql.PrivateRouting = privateRouting.(bool)
	} else {
		mysql.PrivateRouting = defaultPrivateRouting
	}

	if size, ok := platformConfig["size"]; ok {
		mysql.Size = size.(int)
	} else {
		mysql.Size = defaultSize
	}

	if instanceType, ok := platformConfig["instanceType"]; ok {
		mysql.InstanceType = instanceType.(string)
	}

	if subnetID, ok := platformConfig["subnetID"]; ok {
		mysql.SubnetID = subnetID.(string)
	}

	if databaseName, ok := platformConfig["databaseName"]; ok {
		mysql.DatabaseName = databaseName.(string)
	}

	return mysql.Validate()
}

// GenerateDBSecret generates Kubernetes Secret resource to store the host address, username
// and password of the local MySQL database instance.
func (mysql *MySQL) GenerateDBSecret(request *module.GeneratorRequest, hostAddress, username, password string) (
	*apiv1.Resource, *apiv1.Patcher, error,
) {
	// Create the data map of Kubernetes Secret storing the database host address, username
	// and password.
	data := make(map[string]string)
	data["hostAddress"] = hostAddress
	data["username"] = username
	data["password"] = password

	// Create the Kubernetes Secret.
	secret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysql.DatabaseName + dbResSuffix,
			Namespace: request.Project,
		},
		StringData: data,
	}

	resourceID := module.KubernetesResourceID(secret.TypeMeta, secret.ObjectMeta)
	resource, err := module.WrapK8sResourceToKusionResource(resourceID, secret)
	if err != nil {
		return nil, nil, err
	}

	// Inject the database credentials into the workload as the environment variables with
	// Kusion resource patcher.
	hostAddressKey := dbHostAddressEnv + "_" + strings.ToUpper(strings.ReplaceAll(mysql.DatabaseName, "-", "_"))
	usernameKey := dbUsernameEnv + "_" + strings.ToUpper(strings.ReplaceAll(mysql.DatabaseName, "-", "_"))
	passwordKey := dbPasswordEnv + "_" + strings.ToUpper(strings.ReplaceAll(mysql.DatabaseName, "-", "_"))

	envVars := []v1.EnvVar{
		{
			Name: hostAddressKey,
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: secret.Name,
					},
					Key: "hostAddress",
				},
			},
		},
		{
			Name: usernameKey,
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: secret.Name,
					},
					Key: "username",
				},
			},
		},
		{
			Name: passwordKey,
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: secret.Name,
					},
					Key: "password",
				},
			},
		},
	}

	patcher := &apiv1.Patcher{
		Environments: envVars,
	}

	return resource, patcher, nil
}

// GenerateTFRandomPassword generates the terraform random_password resource as the password
// of the cloud provided MySQL database instance.
func (mysql *MySQL) GenerateTFRandomPassword(request *module.GeneratorRequest) (*apiv1.Resource, string, error) {
	resAttrs := map[string]any{
		"length":           16,
		"special":          true,
		"override_special": "_",
	}

	// Set the random_password provider with the default provider config.
	randomPasswordProvider := defaultRandomProviderCfg

	id, err := module.TerraformResourceID(randomPasswordProvider, randomPassword, mysql.DatabaseName+dbResSuffix)
	if err != nil {
		return nil, "", err
	}

	resExts, err := module.TerraformProviderExtensions(randomPasswordProvider, nil, randomPassword)
	if err != nil {
		return nil, "", err
	}

	resource, err := module.WrapTFResourceToKusionResource(id, resAttrs, resExts, nil)
	if err != nil {
		return nil, "", err
	}

	return resource, id, nil
}

// Validate validates whether the input of a MySQL database instance is valid.
func (mysql *MySQL) Validate() error {
	if mysql.Type == CloudDBType && mysql.InstanceType == "" {
		return ErrEmptyInstanceTypeForCloudDB
	}

	return nil
}

// GenerateDefaultMySQLName generates the default name of the MySQL instance.
func GenerateDefaultMySQLName(projectName, stackName, appName string) string {
	strs := []string{projectName, stackName, appName, dbEngine}

	return strings.Join(strs, "-")
}

// GetCloudProviderType returns the cloud provider type of the MySQL instance.
func GetCloudProviderType(platformConfig apiv1.GenericConfig) (string, error) {
	if platformConfig == nil {
		return "", workspace.ErrEmptyModuleConfigBlock
	}

	if cloud, ok := platformConfig["cloud"]; ok {
		return cloud.(string), nil
	}

	return "", ErrEmptyCloudProviderType
}

// IsPublicAccessible returns whether the mysql database instance is publicly
// accessible according to the securityIPs.
func IsPublicAccessible(securityIPs []string) bool {
	var parsedIP net.IP
	for _, ip := range securityIPs {
		if IsIPAddress(ip) {
			parsedIP = net.ParseIP(ip)
		} else if IsCIDR(ip) {
			parsedIP, _, _ = net.ParseCIDR(ip)
		}

		if parsedIP != nil && !parsedIP.IsPrivate() {
			return true
		}
	}

	return false
}

// IsIPAddress returns whether the input string is a valid ip address.
func IsIPAddress(ipStr string) bool {
	ip := net.ParseIP(ipStr)

	return ip != nil
}

// IsCIDR returns whether the input string is a valid CIDR record.
func IsCIDR(cidrStr string) bool {
	_, _, err := net.ParseCIDR(cidrStr)

	return err == nil
}

func main() {
	server.Start(&MySQL{})
}
