package vubvub_test

import (
	"bytes"
	"testing"

	"github.com/oclaussen/vubvub"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const testConfig = `
---
vms:
  testvm:
    - name: somedevice
      vendorid: '0x1234'
      productid: '0x5678'
`

func TestReadConfig(t *testing.T) {
	t.Parallel()

	viper.SetConfigType("yaml")
        err := viper.ReadConfig(bytes.NewBuffer([]byte(testConfig)))
        assert.Nil(t, err)

	config := map[string]vubvub.USBDeviceList{}

	err = viper.UnmarshalKey("vms", &config)
	assert.Nil(t, err)

	ds, ok := config["testvm"]
	assert.True(t, ok)
	assert.Equal(t, 1, len(ds))

	d := ds[0]
	assert.Equal(t, "somedevice", d.Name)
	assert.Equal(t, "0x1234", d.VendorID)
	assert.Equal(t, "0x5678", d.ProductID)
}
