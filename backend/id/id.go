package id

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

/**
 * A unique identifier that is both space efficient and provides good sharding behavior.
 */
const (
	// version is the version bit for the ID. Only `0` is a valid version number at
	// this time.
	version = 0

	versionSize  = 1
	sequenceSize = 10
	workerSize   = 21

	sequenceShift  = versionSize
	workerShift    = sequenceShift + sequenceSize
	timestampShift = workerShift + workerSize

	// workerIdMask extracts the upper 21 bits of the workerId
	workerIDMask = 0xfffff800

	// sequenceMask extracts the lower 10 bits of the sequence
	sequenceMask = 0x3ff

	// versionMask extracts the version number from an ID
	versionMask = 0x1
)

var (
	// idEpoch describes the ID epoch (Jan 1, 2010)
	idEpoch = time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC)

	// workerID stores the worker id of this process used when generating a new ID
	workerID uint64 = initWorkerID()

	// sequence is automatically incremented when generating a new LyftId
	// using the NewLyftId() function.
	sequence uint64 = initSequenceNumber()
)

// An InvalidIDError is returned from ParseID if the string cannot be parsed or
// resolves to an invalid ID.
type InvalidIDError string

// Error satisfies the error built-in interface
func (err InvalidIDError) Error() string { return fmt.Sprintf("invalid ID: %s", string(err)) }

// ID describes a globally unique ID
//
// The ID is a unique 64-bit identifier composed of a (MSB first):
// 32-bit unix timestamp
// 21-bit worker id - truncated sha256 hash of hostname and pid
// 10-bit sequence number - randomly offset per-process counter
// 01-bit version identifier, to allow future changes to address space
type ID uint64

// NewID creates a new Lyft ID from the current time.
func NewID() ID { return NewIDWithTime(time.Now()) }

// NewIDWithTime creates a new Lyft ID with the specified time.
func NewIDWithTime(t time.Time) ID {
	id := uint64(t.Unix()-idEpoch.Unix()) << timestampShift
	id |= workerID
	id |= (atomic.AddUint64(&sequence, 1) << 1) & sequenceMask
	return ID(id)
}

// ParseID attempts to parse an ID from a numeric string. The string is
// expected to be a base-10 64-bit unsigned integer. An error is returned if
// the string cannot be converted to an integer or if the resulting value is an
// invalid ID.
func ParseID(s string) (ID, error) {
	_id, err := strconv.ParseUint(s, 10, 64)
	id := ID(_id)

	if err != nil {
		return id, err
	}

	return id, id.Validate()
}

// Valid returns true if the ID is at least 32-bit and has a version of 0.
func (id ID) Valid() bool { return id > math.MaxUint32 && id.Version() == version }

// Validate returns an InvalidIDError if id is invalid.
func (id ID) Validate() error {
	if !id.Valid() {
		return InvalidIDError(id.String())
	}

	return nil
}

// Time returns the timestamp encoded into the ID. The returned value has
// second precision.
func (id ID) Time() time.Time { return time.Unix(int64(id>>timestampShift)+idEpoch.Unix(), 0) }

// Worker returns the worker ID associated with the ID.
func (id ID) Worker() uint32 { return uint32(id&workerIDMask) >> workerShift }

// Sequence returns the sequence number associated with the ID.
func (id ID) Sequence() uint32 { return uint32(id>>sequenceShift) & sequenceMask }

// Version returns the version number identifying the Lyft ID format. The only
// valid value currently is `0`.
func (id ID) Version() uint8 { return uint8(id & versionMask) }

// String satisfies the fmt.Stringer interface. IDs are represented as base-10
// 64-bit unsigned strings.
func (id ID) String() string { return strconv.FormatUint(uint64(id), 10) }

// MarshalJSON satisfies the json.Marshaler interface. Since IDs are up to
// 64-bit values, IDs are encoded as strings.
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

// UnmarshalJSON satisfies the json.Unmarshaler interface. It supports IDs as
// both JSON number values or strings. UnmarshalJSON returns an InvalidIDError
// if the value is invalid.
func (id *ID) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, (*uint64)(id)); err != nil {
		var str string
		if err = json.Unmarshal(b, &str); err != nil {
			return err
		} else if *id, err = ParseID(str); err != nil {
			return err
		}
	}

	return id.Validate()
}

// initWorkerID generates and returns this worker's id. If this function fails to get the hostname, it will panic.
func initWorkerID() uint64 {
	hostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Errorf("cannot get hostname: %v", err))
	}

	pid := make([]byte, 2)
	binary.BigEndian.PutUint16(pid, uint16(os.Getpid()))

	hash := sha256.New()

	_, err = hash.Write([]byte(hostname))
	if err != nil {
		panic(fmt.Errorf("cannot write hostname: %v", err))
	}
	_, err = hash.Write(pid)
	if err != nil {
		panic(fmt.Errorf("cannot write pid: %v", err))
	}

	sum := hash.Sum(nil)
	id := binary.BigEndian.Uint32(sum[len(sum)-4:]) // extract last 4 bytes
	return uint64(id & workerIDMask)
}

// initSequenceNumber returns a random sequence number. If this function fails to read from the random source, it will
// panic.
func initSequenceNumber() uint64 {
	buf := make([]byte, 8)

	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Errorf("cannot read random object id: %v", err))
	}

	return binary.BigEndian.Uint64(buf)
}
