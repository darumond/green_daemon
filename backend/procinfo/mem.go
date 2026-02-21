package procinfo

const KiB uint64 = 1024
const MiB uint64 = 1024 * KiB
const GiB uint64 = 1024 * MiB

func BytesToMiB(bytes uint64) float64 {
	return float64(bytes) / float64(MiB)
}
