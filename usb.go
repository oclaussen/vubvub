package vubvub

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

const (
	// TODO: find documentation which states there actually are and what they mean
	StateCaptured = "Captured"
)

type VMConfig map[string]USBDeviceList

type USBDeviceList []USBDevice

type USBDevice struct {
	UUID      string
	State     string
	Name      string `mapstructure:"name"`
	VendorID  string `mapstructure:"vendorid"`
	ProductID string `mapstructure:"productid"`
}

func ListUSBDevices() ([]*USBDevice, error) {
	stdout, err := vbm("list", "usbhost")
	if err != nil {
		return nil, err
	}

	return ParseDeviceList(stdout)
}

func ParseDeviceList(in string) ([]*USBDevice, error) {
	result := []*USBDevice{}
	current := &USBDevice{}
	re := regexp.MustCompile(`(.+):\s+(.*)`)
	for _, line := range strings.Split(in, "\n") {
		if line == "" {
			continue
		}

		groups := re.FindStringSubmatch(line)
		if groups == nil {
			continue
		}

		switch groups[1] {
		case "UUID":
			current = &USBDevice{UUID: groups[2]}
			result = append(result, current)
		case "Current State":
			current.State = groups[2]
		case "Name":
			current.Name = groups[2]
		case "VendorId":
			id, err := parseDeviceID(groups[2])
			if err != nil {
				return nil, err
			}

			current.VendorID = id
		case "ProductId":
			id, err := parseDeviceID(groups[2])
			if err != nil {
				return nil, err
			}
			current.ProductID = id
		}
	}

	return result, nil
}

func parseDeviceID(in string) (string, error) {
	re := regexp.MustCompile(`(0x[0-9a-f]{4})\s+\([0-9A-F]{4}\)`)
	groups := re.FindStringSubmatch(in)
	if groups == nil || len(groups) != 2 {
		return "", errors.New("Invalid ID")
	}

	return groups[1], nil
}

func FindUSBDevice(spec USBDevice) (*USBDevice, error) {
	ds, err := ListUSBDevices()
	if err != nil {
		return nil, err
	}

	for _, d := range ds {
		if d.Equals(spec) {
			return d, nil
		}
	}

	return nil, errors.New("Device not found")
}

func (d *USBDevice) Equals(other USBDevice) bool {
	return d.VendorID == other.VendorID && d.ProductID == other.ProductID
}

func CreateAll(vm string, ds USBDeviceList) error {
	for idx, d := range ds {
		if err := Create(vm, idx, d); err != nil {
			return err
		}
	}

	return nil
}

func Create(vm string, idx int, d USBDevice) error {
	_, err := vbm(
		"usbfilter", "add", strconv.Itoa(idx),
		"--target", vm,
		"--name", d.Name,
		"--vendorid", d.VendorID,
		"--productid", d.ProductID,
	)

	return err
}

func RemoveAll(vm string, ds USBDeviceList) error {
	for idx, _ := range ds {
		if err := Remove(vm, idx); err != nil {
			return err
		}
	}

	return nil
}

func Remove(vm string, idx int) error {
	_, err := vbm(
		"usbfilter", "remove", strconv.Itoa(idx),
		"--target", vm,
	)

	return err
}

func IsAttached(_ string, spec USBDevice) (bool, error) {
	d, err := FindUSBDevice(spec)
	if err != nil {
		return false, err
	}

	// TODO: figure out if it is captured by the right VM
	return d.State == StateCaptured, nil
}

func Attach(vm string, spec USBDevice) error {
	d, err := FindUSBDevice(spec)
	if err != nil {
		return err
	}

	if _, err := vbm("controlvm", vm, "usbattach", d.UUID); err != nil {
		return err
	}

	return nil
}

func Detach(vm string, spec USBDevice) error {
	d, err := FindUSBDevice(spec)
	if err != nil {
		return err
	}

	if _, err := vbm("controlvm", vm, "usbdetach", d.UUID); err != nil {
		return err
	}

	return nil
}

func Toggle(vm string, spec USBDevice) error {
	if attached, err := IsAttached(vm, spec); err != nil {
		return err
	} else if attached {
		return Detach(vm, spec)
	} else {
		return Attach(vm, spec)
	}
}
