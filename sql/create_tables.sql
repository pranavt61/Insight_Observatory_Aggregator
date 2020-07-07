CREATE TABLE IF NOT EXISTS blocks (
    height INT UNSIGNED,
    hash VARCHAR(64) UNIQUE,
    prev_hash VARCHAR(64),
    coinbase_tx VARCHAR(64),
    num_tx INT UNSIGNED,
    difficulty DOUBLE,
    block_size INT UNSIGNED,
    miner_time BIGINT UNSIGNED,
    network_time BIGINT UNSIGNED
);

CREATE TABLE IF NOT EXISTS inv (
    hash VARCHAR(64),
    peer_ip VARCHAR(20),
    network_time BIGINT UNSIGNED,
    session_id INT REFERENCES obs_sessions(session_id)
);

CREATE TABLE IF NOT EXISTS peer_conn (
    peer_ip VARCHAR(20),
    version INT,
    subversion VARCHAR(64),
    start_height INT,
    services BIGINT UNSIGNED,
    peer_time BIGINT UNSIGNED,
    network_time BIGINT UNSIGNED,
    disconnect_time BIGINT UNSIGNED,
    session_id INT REFERENCES obs_sessions(session_id)
);

CREATE TABLE IF NOT EXISTS obs_sessions (
    session_id INT NOT NULL AUTO_INCREMENT,
    ip VARCHAR(32),
    name VARCHAR(32),
    software_version VARCHAR(32),
    start_time BIGINT UNSIGNED,
    end_time BIGINT UNSIGNED,
    PRIMARY KEY (session_id)
);
