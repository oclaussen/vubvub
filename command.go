package vubvub

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Execute() int {
	if err := NewRootCommand().Execute(); err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", err.Error())

		return 1
	}

	return 0
}

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "vubvub",
	}

	cmd.AddCommand(NewApplyCommand())
	cmd.AddCommand(NewRemoveCommand())
	cmd.AddCommand(NewToggleCommand())

	return cmd
}

func NewApplyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "apply",
		Short: "creates all configured USB filters for a VM",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			vm := args[0]

			config, err := ReadConfig()
			if err != nil {
				return err
			}

			list, ok := config[vm]
			if !ok {
				return errors.New("VM not found")
			}

			return CreateAll(vm, list)
		},
	}
}

func NewRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "removes all configured USB filters from a VM",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			vm := args[0]

			config, err := ReadConfig()
			if err != nil {
				return err
			}

			list, ok := config[vm]
			if !ok {
				return errors.New("VM not found")
			}

			return RemoveAll(vm, list)
		},
	}
}

func NewToggleCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "toggle",
		Short: "toggles a captured device between host and VM",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			vm := args[0]
			name := args[1]

			config, err := ReadConfig()
			if err != nil {
				return err
			}

			ds, ok := config[vm]
			if !ok {
				return errors.New("VM not found")
			}

			for _, d := range ds {
				if d.Name == name {
					return Toggle(vm, d)
				}
			}

			return errors.New("USB Device not found")
		},
	}
}

func ReadConfig() (VMConfig, error) {
	viper.SetConfigName(".vubvub")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := map[string]USBDeviceList{}
	if err := viper.UnmarshalKey("vms", &config); err != nil {
		return nil, err
	}

	return config, nil
}
