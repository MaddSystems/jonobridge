package capabilities

type Capability interface {
    Name() string
    Version() string
    GetDataContextName() string
    Initialize(imei string) error
    GetSnapshot() map[string]interface{}
}

// SnapshotProvider allows capabilities to contribute their own audit data
// Each capability implements this to self-report its snapshot without
// modifying the central audit code.
type SnapshotProvider interface {
    // GetSnapshotData returns capability-specific data for audit snapshots
    // Keys should be descriptive (e.g., "buffer_circular", "jammer_metrics")
    // Returns nil if capability has no data to contribute
    GetSnapshotData(imei string) map[string]interface{}
}
