package cmd_geoip

func Download() {
	lib, err := StartGeoIpComponent()
	if err != nil {
		panic(err)
	}

	upgrade_err := lib.Upgrade(true)
	if upgrade_err != nil {
		panic(upgrade_err)
	}
}
