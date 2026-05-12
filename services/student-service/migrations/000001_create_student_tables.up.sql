CREATE TABLE IF NOT EXISTS students (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL,
    email TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL,
    license_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS groups (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL,
    name TEXT NOT NULL,
    parent_id UUID REFERENCES groups(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS student_group_members (
    student_id UUID REFERENCES students(id) ON DELETE CASCADE,
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    PRIMARY KEY (student_id, group_id)
);
