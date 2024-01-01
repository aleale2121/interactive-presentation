-- Enable the uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the Presentation table with UUID primary key
CREATE TABLE presentations (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    currentPollIndex INTEGER DEFAULT -1
);

-- Create the Polls table with UUID primary key
CREATE TABLE polls (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    presentationID UUID NOT NULL,
    question TEXT NOT NULL,
    pollIndex int NOT NULL,
    createdAt TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (presentationID) REFERENCES presentations(id),
    UNIQUE(presentationID, pollIndex)
);

-- Create the Options table (for poll options) with UUID primary key
CREATE TABLE options (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    pollID UUID NOT NULL,
    optionKey TEXT NOT NULL,
    optionValue TEXT NOT NULL,
    FOREIGN KEY (pollID) REFERENCES polls(id),
    UNIQUE(pollID, optionKey)
);

-- Create the Votes table with UUID primary key
CREATE TABLE votes (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    pollID UUID NOT NULL,
    optionKey TEXT NOT NULL,
    clientID TEXT NOT NULL,
    FOREIGN KEY (pollID, optionKey) REFERENCES options(pollID, optionKey),
    UNIQUE(pollID, clientID)
);

