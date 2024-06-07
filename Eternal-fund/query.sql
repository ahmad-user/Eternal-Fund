-- Table structure for table `users`
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    occupation VARCHAR(255),
    email VARCHAR(255),
    password_hash VARCHAR(255),
    avatar_file_name VARCHAR(255),
    role VARCHAR(255),  
    -- token VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
-- Table structure for table `campaigns`
CREATE TABLE campaigns (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    name VARCHAR(255),
    short_description VARCHAR(255),
    description TEXT,
    perks TEXT,
    backer_count INTEGER,
    goal_amount INTEGER,
    current_amount INTEGER,
    slug VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);
-- Table structure for table `campaign_images`
CREATE TABLE campaign_images (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER,
    file_name VARCHAR(255),
    is_primary SMALLINT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE SET NULL
);
-- Table structure for table `transactions`
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER,
    user_id INTEGER,
    amount INTEGER,
    status VARCHAR(255),    
    code VARCHAR(255),
    payment_url VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

