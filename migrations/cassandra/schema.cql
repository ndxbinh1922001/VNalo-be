-- Cassandra Schema for VNalo Chat
-- Run this manually or through a script

-- Create keyspace
CREATE KEYSPACE IF NOT EXISTS vnalo_chat
WITH replication = {
    'class': 'SimpleStrategy',
    'replication_factor': 1
}
AND durable_writes = true;

-- Messages table (partitioned by conversation_id, sorted by timestamp)
CREATE TABLE IF NOT EXISTS vnalo_chat.messages (
    conversation_id BIGINT,
    message_id TIMEUUID,
    sender_id BIGINT,
    parent_message_id TIMEUUID,
    message_type TEXT,
    content TEXT,
    metadata MAP<TEXT, TEXT>,
    status TEXT,
    is_edited BOOLEAN,
    is_deleted BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (conversation_id, message_id)
) WITH CLUSTERING ORDER BY (message_id DESC)
  AND gc_grace_seconds = 864000
  AND compaction = {
    'class': 'TimeWindowCompactionStrategy',
    'compaction_window_size': '1',
    'compaction_window_unit': 'DAYS'
  };

-- Messages by user (for search across conversations)
CREATE TABLE IF NOT EXISTS vnalo_chat.messages_by_user (
    user_id BIGINT,
    message_id TIMEUUID,
    conversation_id BIGINT,
    sender_id BIGINT,
    content TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY (user_id, message_id)
) WITH CLUSTERING ORDER BY (message_id DESC);

-- Message delivery status
CREATE TABLE IF NOT EXISTS vnalo_chat.message_delivery (
    message_id TIMEUUID,
    user_id BIGINT,
    status TEXT,
    delivered_at TIMESTAMP,
    read_at TIMESTAMP,
    PRIMARY KEY (message_id, user_id)
);

-- Message reactions
CREATE TABLE IF NOT EXISTS vnalo_chat.message_reactions (
    message_id TIMEUUID,
    user_id BIGINT,
    emoji TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY (message_id, user_id, emoji)
);

-- Reaction counts (for quick display)
CREATE TABLE IF NOT EXISTS vnalo_chat.reaction_counts (
    message_id TIMEUUID,
    emoji TEXT,
    count COUNTER,
    PRIMARY KEY (message_id, emoji)
);

-- Message threads (replies)
CREATE TABLE IF NOT EXISTS vnalo_chat.message_threads (
    parent_message_id TIMEUUID,
    reply_message_id TIMEUUID,
    sender_id BIGINT,
    content TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY (parent_message_id, reply_message_id)
) WITH CLUSTERING ORDER BY (reply_message_id DESC);

-- User typing indicators (with TTL)
CREATE TABLE IF NOT EXISTS vnalo_chat.typing_indicators (
    conversation_id BIGINT,
    user_id BIGINT,
    is_typing BOOLEAN,
    updated_at TIMESTAMP,
    PRIMARY KEY (conversation_id, user_id)
) WITH default_time_to_live = 10;

