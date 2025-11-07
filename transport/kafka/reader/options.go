package reader // import "go.microcore.dev/framework/transport/kafka/reader"

import (
	"time"

	_ "go.microcore.dev/framework"

	"github.com/segmentio/kafka-go"
)

type Option func(*kafka.ReaderConfig)

// The list of broker addresses used to connect to the kafka cluster.
func WithBrokers(brokers []string) Option {
	return func(k *kafka.ReaderConfig) {
		k.Brokers = brokers
	}
}

// GroupID holds the optional consumer group id.  If GroupID is specified, then
// Partition should NOT be specified e.g. 0
func WithGroupID(groupID string) Option {
	return func(k *kafka.ReaderConfig) {
		k.GroupID = groupID
	}
}

// GroupTopics allows specifying multiple topics, but can only be used in
// combination with GroupID, as it is a consumer-group feature. As such, if
// GroupID is set, then either Topic or GroupTopics must be defined.
func WithGroupTopics(groupTopics []string) Option {
	return func(k *kafka.ReaderConfig) {
		k.GroupTopics = groupTopics
	}
}

// The topic to read messages from.
func WithTopic(topic string) Option {
	return func(k *kafka.ReaderConfig) {
		k.Topic = topic
	}
}

// Partition to read messages from.  Either Partition or GroupID may
// be assigned, but not both
func WithPartition(partition int) Option {
	return func(k *kafka.ReaderConfig) {
		k.Partition = partition
	}
}

// An dialer used to open connections to the kafka server. This field is
// optional, if nil, the default dialer is used instead.
func WithDialer(dialer *kafka.Dialer) Option {
	return func(k *kafka.ReaderConfig) {
		k.Dialer = dialer
	}
}

// The capacity of the internal message queue, defaults to 100 if none is
// set.
func WithQueueCapacity(queueCapacity int) Option {
	return func(k *kafka.ReaderConfig) {
		k.QueueCapacity = queueCapacity
	}
}

// MinBytes indicates to the broker the minimum batch size that the consumer
// will accept. Setting a high minimum when consuming from a low-volume topic
// may result in delayed delivery when the broker does not have enough data to
// satisfy the defined minimum.
//
// Default: 1
func WithMinBytes(minBytes int) Option {
	return func(k *kafka.ReaderConfig) {
		k.MinBytes = minBytes
	}
}

// MaxBytes indicates to the broker the maximum batch size that the consumer
// will accept. The broker will truncate a message to satisfy this maximum, so
// choose a value that is high enough for your largest message size.
//
// Default: 1MB
func WithMaxBytes(maxBytes int) Option {
	return func(k *kafka.ReaderConfig) {
		k.MaxBytes = maxBytes
	}
}

// Maximum amount of time to wait for new data to come when fetching batches
// of messages from kafka.
//
// Default: 10s
func WithMaxWait(maxWait time.Duration) Option {
	return func(k *kafka.ReaderConfig) {
		k.MaxWait = maxWait
	}
}

// ReadBatchTimeout amount of time to wait to fetch message from kafka messages batch.
//
// Default: 10s
func WithReadBatchTimeout(readBatchTimeout time.Duration) Option {
	return func(k *kafka.ReaderConfig) {
		k.ReadBatchTimeout = readBatchTimeout
	}
}

// ReadLagInterval sets the frequency at which the reader lag is updated.
// Setting this field to a negative value disables lag reporting.
func WithReadLagInterval(readLagInterval time.Duration) Option {
	return func(k *kafka.ReaderConfig) {
		k.ReadLagInterval = readLagInterval
	}
}

// GroupBalancers is the priority-ordered list of client-side consumer group
// balancing strategies that will be offered to the coordinator.  The first
// strategy that all group members support will be chosen by the leader.
//
// Default: [Range, RoundRobin]
//
// Only used when GroupID is set
func WithGroupBalancers(groupBalancers []kafka.GroupBalancer) Option {
	return func(k *kafka.ReaderConfig) {
		k.GroupBalancers = groupBalancers
	}
}

// HeartbeatInterval sets the optional frequency at which the reader sends the consumer
// group heartbeat update.
//
// Default: 3s
//
// Only used when GroupID is set
func WithHeartbeatInterval(heartbeatInterval time.Duration) Option {
	return func(k *kafka.ReaderConfig) {
		k.HeartbeatInterval = heartbeatInterval
	}
}

// CommitInterval indicates the interval at which offsets are committed to
// the broker.  If 0, commits will be handled synchronously.
//
// Default: 0
//
// Only used when GroupID is set
func WithCommitInterval(commitInterval time.Duration) Option {
	return func(k *kafka.ReaderConfig) {
		k.CommitInterval = commitInterval
	}
}

// PartitionWatchInterval indicates how often a reader checks for partition changes.
// If a reader sees a partition change (such as a partition add) it will rebalance the group
// picking up new partitions.
//
// Default: 5s
//
// Only used when GroupID is set and WatchPartitionChanges is set.
func WithPartitionWatchInterval(partitionWatchInterval time.Duration) Option {
	return func(k *kafka.ReaderConfig) {
		k.PartitionWatchInterval = partitionWatchInterval
	}
}

// WatchForPartitionChanges is used to inform kafka-go that a consumer group should be
// polling the brokers and rebalancing if any partition changes happen to the topic.
func WithWatchPartitionChanges(watchPartitionChanges bool) Option {
	return func(k *kafka.ReaderConfig) {
		k.WatchPartitionChanges = watchPartitionChanges
	}
}

// SessionTimeout optionally sets the length of time that may pass without a heartbeat
// before the coordinator considers the consumer dead and initiates a rebalance.
//
// Default: 30s
//
// Only used when GroupID is set
func WithSessionTimeout(sessionTimeout time.Duration) Option {
	return func(k *kafka.ReaderConfig) {
		k.SessionTimeout = sessionTimeout
	}
}

// RebalanceTimeout optionally sets the length of time the coordinator will wait
// for members to join as part of a rebalance.  For kafka servers under higher
// load, it may be useful to set this value higher.
//
// Default: 30s
//
// Only used when GroupID is set
func WithRebalanceTimeout(rebalanceTimeout time.Duration) Option {
	return func(k *kafka.ReaderConfig) {
		k.RebalanceTimeout = rebalanceTimeout
	}
}

// JoinGroupBackoff optionally sets the length of time to wait between re-joining
// the consumer group after an error.
//
// Default: 5s
func WithJoinGroupBackoff(joinGroupBackoff time.Duration) Option {
	return func(k *kafka.ReaderConfig) {
		k.JoinGroupBackoff = joinGroupBackoff
	}
}

// RetentionTime optionally sets the length of time the consumer group will be saved
// by the broker. -1 will disable the setting and leave the
// retention up to the broker's offsets.retention.minutes property. By
// default, that setting is 1 day for kafka < 2.0 and 7 days for kafka >= 2.0.
//
// Default: -1
//
// Only used when GroupID is set
func WithRetentionTime(retentionTime time.Duration) Option {
	return func(k *kafka.ReaderConfig) {
		k.RetentionTime = retentionTime
	}
}

// StartOffset determines from whence the consumer group should begin
// consuming when it finds a partition without a committed offset.  If
// non-zero, it must be set to one of FirstOffset or LastOffset.
//
// Default: FirstOffset
//
// Only used when GroupID is set
func WithStartOffset(startOffset int64) Option {
	return func(k *kafka.ReaderConfig) {
		k.StartOffset = startOffset
	}
}

// BackoffDelayMin optionally sets the smallest amount of time the reader will wait before
// polling for new messages
//
// Default: 100ms
func WithReadBackoffMin(readBackoffMin time.Duration) Option {
	return func(k *kafka.ReaderConfig) {
		k.ReadBackoffMin = readBackoffMin
	}
}

// BackoffDelayMax optionally sets the maximum amount of time the reader will wait before
// polling for new messages
//
// Default: 1s
func WithReadBackoffMax(readBackoffMax time.Duration) Option {
	return func(k *kafka.ReaderConfig) {
		k.ReadBackoffMax = readBackoffMax
	}
}

// If not nil, specifies a logger used to report internal changes within the
// writer.
func WithLogger(logger kafka.Logger) Option {
	return func(k *kafka.ReaderConfig) {
		k.Logger = logger
	}
}

// ErrorLogger is the logger used to report errors. If nil, the writer falls
// back to using Logger instead.
func WithErrorLogger(errorLogger kafka.Logger) Option {
	return func(k *kafka.ReaderConfig) {
		k.ErrorLogger = errorLogger
	}
}

// IsolationLevel controls the visibility of transactional records.
// ReadUncommitted makes all records visible. With ReadCommitted only
// non-transactional and committed records are visible.
func WithIsolationLevel(isolationLevel kafka.IsolationLevel) Option {
	return func(k *kafka.ReaderConfig) {
		k.IsolationLevel = isolationLevel
	}
}

// Limit of how many attempts to connect will be made before returning the error.
//
// The default is to try 3 times.
func WithMaxAttempts(maxAttempts int) Option {
	return func(k *kafka.ReaderConfig) {
		k.MaxAttempts = maxAttempts
	}
}

// OffsetOutOfRangeError indicates that the reader should return an error in
// the event of an OffsetOutOfRange error, rather than retrying indefinitely.
// This flag is being added to retain backwards-compatibility, so it will be
// removed in a future version of kafka-go.
func WithOffsetOutOfRangeError(offsetOutOfRangeError bool) Option {
	return func(k *kafka.ReaderConfig) {
		k.OffsetOutOfRangeError = offsetOutOfRangeError
	}
}
