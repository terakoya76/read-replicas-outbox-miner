package tracker

// TrackerClient provide accesses for DataSource Progress
// by implements TrackerStorage interface
type TrackerClient struct {
	TrackerStorage
}

// BuildMySQLClient builds MySQL specific SourceClient
func BuildTrackerClient() (*TrackerClient, error) {
	ts, err := BuildTrackerStorage()
	if err != nil {
		return nil, err
	}

	tc := &TrackerClient{ts}
	return tc, nil
}

// GetNextPosition returns next id after the last-published record
func (tc *TrackerClient) GetNextPosition(dbName string, tableName string) (Position, error) {
	prgs, err := tc.GetProgress(dbName, tableName)
	if err != nil {
		return 0, err
	}

	return prgs.Position + 1, nil
}

// UpdatePosition updates Progress's Position
func (tc *TrackerClient) UpdatePosition(dbName string, tableName string, position Position) error {
	prgs, err := tc.GetProgress(dbName, tableName)
	if err != nil {
		return err
	}

	prgs.Position += position
	prgs, err = tc.UpdateProgress(prgs)
	if err != nil {
		return err
	}

	return nil
}
