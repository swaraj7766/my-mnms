package mnms

import (
	"reflect"
	"testing"
)

// TestMNMSConfigReadWrite test mnms config read write
func TestMNMSConfigReadWrite(t *testing.T) {
	var err error
	QC.OwnPublicKeys, err = GenerateOwnPublickey()
	if err != nil {
		t.Fatal("generate own public key fail", err)
	}
	// clear mnms config
	err = cleanMNMSConfig()
	if err != nil {
		t.Log("clean mnms config fail", err)
	}

	// Init default mnms config
	err = InitDefaultMNMSConfigIfNotExist()
	if err != nil {
		t.Fatal("init default mnms config fail", err)
	}
	// expect config
	adminPass, err := GenPassword(QC.Name, "admin")
	if err != nil {
		t.Fatal("gen admin password fail", err)
	}

	expectConfig := &MNMSConfig{
		Users: []UserConfig{
			{
				Name:     "admin",
				Role:     MNMSAdminRole,
				Password: adminPass,
			},
		},
	}
	// Test readconfig
	readedConfig, err := GetMNMSConfig()
	if err != nil {
		t.Fatal("read mnms config fail", err)
	}
	if !reflect.DeepEqual(readedConfig, expectConfig) {
		t.Fatal("expect mnms config", expectConfig, "but got", readedConfig)
	}

	// modify config
	expectConfig.Users = append(expectConfig.Users, UserConfig{
		Name:     "test",
		Role:     MNMSSuperUserRole,
		Password: "testrawpass",
	})

}
