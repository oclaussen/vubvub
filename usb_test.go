package vubvub_test

import (
	"testing"

	"github.com/oclaussen/vubvub"
	"github.com/stretchr/testify/assert"
)

const testListOutput = `
Host USB Devices:

UUID:               XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
VendorId:           0x1234 (1234)
ProductId:          0x5678 (5678)
Revision:           2.1 (0201)
Port:               2
USB version/speed:  0/High
Manufacturer:       Foo Inc.
Product:            Bar Baz 1
SerialNumber:       0000000000000000
Current State:      Busy

UUID:               YYYYYYYY-YYYY-YYYY-YYYY-YYYYYYYYYYYY
VendorId:           0x90ab (90AB)
ProductId:          0xcdef (CDEF)
Revision:           2.1 (0201)
Port:               6
USB version/speed:  0/High
Manufacturer:       Foo Corp
Product:            Bla Blub 2
SerialNumber:       0000000000000000
Current State:      Captured

`

func TestListUSBDevices(t *testing.T) {
	t.Parallel()

	ds, err := vubvub.ParseDeviceList(testListOutput)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(ds))

	d := ds[0]
	assert.Equal(t, "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX", d.UUID)
	assert.Equal(t, "0x1234", d.VendorID)
	assert.Equal(t, "0x5678", d.ProductID)
	assert.Equal(t, "Busy", d.State)

	d = ds[1]
	assert.Equal(t, "YYYYYYYY-YYYY-YYYY-YYYY-YYYYYYYYYYYY", d.UUID)
	assert.Equal(t, "0x90ab", d.VendorID)
	assert.Equal(t, "0xcdef", d.ProductID)
	assert.Equal(t, "Captured", d.State)
}
