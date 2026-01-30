package migration

func upgradeDatabaseToVersion3_4_0(fromVersion string) error {
	return updateSchemaVersion(version3_4_0)
}
