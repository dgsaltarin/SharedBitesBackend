-- Migration: Create groups and group_members tables
-- This file is for reference only. The actual migration is handled by GORM AutoMigrate.

-- Create groups table
CREATE TABLE IF NOT EXISTS groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description VARCHAR(500),
    owner_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create index on owner_id for faster queries
CREATE INDEX IF NOT EXISTS idx_groups_owner_id ON groups(owner_id);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_groups_created_at ON groups(created_at);

-- Create group_members table
CREATE TABLE IF NOT EXISTS group_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_group_members_group_id FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

-- Create index on group_id for faster queries
CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON group_members(group_id);

-- Add comments for documentation
COMMENT ON TABLE groups IS 'Stores group information for expense sharing';
COMMENT ON TABLE group_members IS 'Stores member information for each group';
COMMENT ON COLUMN groups.owner_id IS 'Reference to the user who owns this group';
COMMENT ON COLUMN group_members.group_id IS 'Reference to the group this member belongs to';
COMMENT ON COLUMN group_members.name IS 'Name of the group member'; 