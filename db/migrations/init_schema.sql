-- ===============================================
-- Configuration Schema: reference data
-- ===============================================
CREATE SCHEMA IF NOT EXISTS conf;

CREATE TABLE conf.language (
    code VARCHAR(10) PRIMARY KEY,
    name VARCHAR(50)
);

COMMENT ON TABLE conf.language IS 'Table of available languages for the application';
COMMENT ON COLUMN conf.language.code IS 'Language code (e.g., en, es)';
COMMENT ON COLUMN conf.language.name IS 'Full name of the language';

INSERT INTO conf.language (code, name) VALUES
    ('en', 'English'),
    ('es', 'Español'),
    ('zh', '中文 (Chinese)');
    
    

-- ===============================================
-- Authorization Schema: users, roles
-- ===============================================
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth."user" (
    user_id UUID DEFAULT gen_random_uuid() UNIQUE NOT NULL,
    user_name VARCHAR(50) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    photo BYTEA DEFAULT decode('iVBORw0KGgoAAAANSUhEUgAAAOEAAADhCAMAAAAJbSJIAAAAV1BMVEX6+vqPj4////+Li4u5ubn8/PyIiIiFhYWJiYnk5OShoaGnp6fT09Pn5+eRkZHu7u7Z2dn19fXCwsKamprHx8exsbHOzs7X19eurq6/v7+jo6Pe3t6WlpZaNtXmAAAE3UlEQVR4nO2d25aqOhBFsUIRbgqI4AX//zsP0fa0vUfbBoKm4ljzpfvROapIIGSFKAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIEWamG+P/vn/Owhi5Juu3XZHnp6Lblutm1PT9q5aDKRriVulEqZVBqUSr9pjxh0gyrWOlr273KL05Vh/gyDTkv+jdJIsscEemrNUP9K7oU0W+f6UD1Bz+9rs4xuEOrFSrR/15T7rJwiwjU/y8gF9l3IWoyHxKLAVHxS68AYej1qZDbyRFaIocbaYIjhNHHlajTqygIS2CUqRiquDYqHFAinS0H2S+0WUwijzYThP/KFahjDY8vUWvtIEUkeK5hkkYMz9X83rUoJsQ+pTy2YIrFcJ4ytn8EoZRRCocBEMoostVeFH0LfAUOs4dSK8kpfQ2pbOT4Gp1Et6mvHZr0vEOXPhYQ7vU0TCphRueHAXFj6bsKij95pSrOY9N/xQxktymPLgbJqKfobh3HWhGw0GyIW3d5vuLoeg5f/6j4TdpL9qwczdUoh+DYWhDuhPdpY5PFhdD2dfhboGxdC/ZkMsFZvxMtOH64+9pGnfDjWTBBR7xxT/ku08XqejpcGzTvWub6rXsLnW/EIVfhu7LGNIXMdxnRC16NjRw5FZD2as0F9xuTWU//l7hxmVNeCO/hKaI89dqdAAljBxe4wdxFRp4P7dPpc/2/zNnv5AhFT8X3uBonuE5FMG57/IT4e/VfkDldEU9hFPCyCx+T1XU+6AEzaw4TVH3gQmaZbcpisFV0DDlWkzD3K1Pa8ud0EnbBClotut3NmXUx9B2sd9B2fmZo86DjgVFTOXmr4d+fa4DLuAV4rJ9EF5TOg/fz2ACiBud/rRUiT5vPyF+eIWJ1v3hnGidGMY/566sPione00CR1U21HU9rCs2YWffP+kV8A3fPwQAAIAP7k/1WApJkwpTM/THeFmOfRYJuelhGgo13nYuTaJX3VqCI1W5awDhIUof/K+hzlkZneKY+F7Bmb4uOhXPq3DUv1rQ85t916CaHcrjtegSF51gePDWp1y/o4Q+X5y+p4RjETtPRVxiq6UlnmrovkvPFl9tusS2dTt87SNaInpgh68IBh3eJLhSWxjCcK7h265DX4afP9IsEDa0w1cUaomQkx2+olBLhJwsDT09IrqfEGFt6CkKxY17cNsOb3ujqX2Tobfj+N41mCbeUqVzT56bis+T6t4i6HN/+3va1Gde7z3zhdfd0e4H7jzHb5rN7fg5OzwfUjc3WmGPOvp9NeOW47Iy9P16jXavvf3W/o/+ovyVfeptufsO19Do34IiwmxLnO/1EP8vuQ30sttTJeWIjFcpihE0W/Jf0KhqI0fQbDmZeIz+c9JWxjV4g7lYtlN1LGGz0A+of/jBnOkoJTGMSM1iZdSdzNMhmYbzEiOObkVsZ/sVpv7PDJCdn+wcDfH+UQbIhiByQkzZQc8qpEqSWG5/3sMUlYVOJn5nRieHOpxPzfEoWXcbW0uT8oqHcPS+GH9wVXZ33wT81c18JzCP96F+DfGS5lrvt4d8oy65tTS9bJZOr/k1dc67XV1Foae8Lrv4uamqoS77frfd7nZ9X9ZZ1TQsbEe+E1+Zte+gARJsAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACJP/AAFSQ7wNy+LTAAAAAElFTkSuQmCC', 'base64'),
    language_code VARCHAR(10) DEFAULT 'es',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_auth_user_language FOREIGN KEY (language_code) REFERENCES conf.language(code)
);

COMMENT ON TABLE auth."user" IS 'Table of application users';
COMMENT ON COLUMN auth."user".user_id IS 'Unique identifier for the user';
COMMENT ON COLUMN auth."user".user_name IS 'Name of the user';
COMMENT ON COLUMN auth."user".photo IS 'User''s photo';
COMMENT ON COLUMN auth."user".language_code IS 'Identifier of the user''s preferred language';

CREATE TABLE auth.role (
    role_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE auth.role IS 'Table containing the defined roles in the application';
COMMENT ON COLUMN auth.role.role_id IS 'Unique identifier for the role';
COMMENT ON COLUMN auth.role.role_name IS 'Name of the role';
COMMENT ON COLUMN auth.role.description IS 'Description of the role';

INSERT INTO auth.role (role_name, description)
VALUES
    ( 'FIELD AGENT', 'User responsible for collecting and uploading data from the field with limited access to the application'),
    ( 'VIEWER', 'User with read-only access to view dashboards and reports'),
    ( 'EDITOR', 'User with the ability to edit content and update records'),
    ( 'MANAGER', 'User with the ability to manage users'),
    ( 'ADMIN', 'User with full administrative privileges, including managing roles and system configurations');



CREATE TABLE auth.user_role (
    user_id UUID,
    role_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID,
    valid_until TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_user_role_user FOREIGN KEY (user_id) REFERENCES auth."user"(user_id),
    CONSTRAINT fk_user_role_role FOREIGN KEY (role_id) REFERENCES auth.role(role_id),
    CONSTRAINT fk_user_role_created_by FOREIGN KEY (created_by) REFERENCES auth."user"(user_id),  
    CONSTRAINT fk_user_role_updated_by FOREIGN KEY (updated_by) REFERENCES auth."user"(user_id)
);


COMMENT ON TABLE auth.user_role IS 'Table for assigning roles to users';
COMMENT ON COLUMN auth.user_role.user_id IS 'Identifier of the user';
COMMENT ON COLUMN auth.user_role.role_id IS 'Identifier of the assigned role';
COMMENT ON COLUMN auth.user_role.valid_until IS 'Date until the assignment is valid';

CREATE TABLE auth.permission (
    permission_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    permission_name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE auth.permission IS 'Table for permissions in the application';
COMMENT ON COLUMN auth.permission.permission_id IS 'Identifier of the permission';
COMMENT ON COLUMN auth.permission.permission_name IS 'Name of the permission';
COMMENT ON COLUMN auth.permission.description IS 'Description of the permission';

INSERT INTO auth.permission (permission_name, description)
VALUES
    ('READ',   'User with permission to view tree census data and metrics, but cannot modify any records.'),
    ('SURVEY',  'Field technician responsible for collecting and uploading tree data from the field, with limited application access.'),
    ('EDIT',   'User with permission to modify existing tree data and update records.'),
    ('DELETE', 'User with permission to delete tree records from the system.'),
    ('MANAGE', 'User with permission to manage roles and permissions within the application.'),
    ('CONFIG',  'User with permissions for system configuration.');


CREATE TABLE auth.role_permission (
    role_name VARCHAR(50),
    permission_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_name, permission_name),
    CONSTRAINT fk_role_permission_permission FOREIGN KEY (permission_name) REFERENCES auth."permission"(permission_name),
    CONSTRAINT fk_permission_role_role FOREIGN KEY (role_name) REFERENCES auth.role(role_name)
);

COMMENT ON TABLE auth.role_permission IS 'Table for role assigment permissions in the application';

INSERT INTO auth.role_permission (role_name, permission_name)
VALUES
    ('FIELD AGENT', 'SURVEY'),
    ('VIEWER', 'READ'),
    ('EDITOR', 'READ'),
    ('EDITOR', 'EDIT'),
    ('MANAGER', 'READ'),
    ('MANAGER', 'MANAGE'),
    ('ADMIN', 'READ'),
    ('ADMIN', 'SURVEY'),
    ('ADMIN', 'EDIT'),
    ('ADMIN', 'DELETE'),
    ('ADMIN', 'MANAGE'),
    ('ADMIN', 'CONFIG');



-- ===============================================
-- Main Schema: core application data
-- ===============================================
CREATE SCHEMA IF NOT EXISTS treesense;

CREATE TABLE treesense."route" (
    route_id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    user_id UUID NOT NULL,
    route GEOMETRY(LineString, 4326) NOT NULL, -- TODO ver si es el mejor formato para almacenar la ruta en formato geoespacial
    start_at TIMESTAMP NOT NULL,
    end_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_route_user FOREIGN KEY (user_id) REFERENCES auth."user"(user_id)
);

COMMENT ON TABLE treesense."route" IS 'Table of user routes';
COMMENT ON COLUMN treesense."route".route_id IS 'Unique identifier for the route';
COMMENT ON COLUMN treesense."route".user_id IS 'Unique identifier for the user than made that route';
COMMENT ON COLUMN treesense."route".route IS 'Route (latitude, longitude)';


CREATE TABLE treesense."tree_species" (
    tree_species_id VARCHAR(100) PRIMARY KEY,  
    description TEXT  
);

COMMENT ON TABLE treesense."tree_species" IS 'Table storing different tree species';
COMMENT ON COLUMN treesense."tree_species".tree_species_id IS 'Unique identifier for the tree species (code)';
COMMENT ON COLUMN treesense."tree_species".description IS 'Additional information about the species';

INSERT INTO treesense."tree_species" (tree_species_id, description) VALUES
('Quercus robur', 'Commonly known as English oak, native to Europe'),
('Pinus sylvestris', 'Scots pine, widely distributed across Eurasia'),
('Acer rubrum', 'Red maple, native to North America');



CREATE TABLE treesense."tree_state" (
    tree_state_id VARCHAR(100) PRIMARY KEY,  
    description TEXT  
);

COMMENT ON TABLE treesense."tree_state" IS 'Table storing different health states of trees';
COMMENT ON COLUMN treesense."tree_state".tree_state_id IS 'Unique identifier for the tree state';
COMMENT ON COLUMN treesense."tree_state".description IS 'Additional information about the state';

INSERT INTO treesense."tree_state" (tree_state_id, description) VALUES
('Healthy', 'Tree is in good condition with no visible issues'),
('Sick', 'Tree shows signs of disease or infestation'),
('Dry', 'Tree appears to be dry or dying');


CREATE TABLE treesense."tree" (
    tree_id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    route_id UUID, --TODO NOT NULL,
    species VARCHAR(100),
    state VARCHAR(100),
    location GEOMETRY(Point, 4326) NOT NULL,
    age INT,
    height FLOAT,
    diameter FLOAT,
    photo_url TEXT CHECK (photo_url ~* '^https?://.+') DEFAULT 'https://userphoto.png',
    description TEXT,
    created_by UUID,
    updated_by UUID,    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_tree_route FOREIGN KEY (route_id) REFERENCES treesense."route"(route_id),
    CONSTRAINT fk_tree_species FOREIGN KEY (species) REFERENCES treesense."tree_species"(tree_species_id),
    CONSTRAINT fk_tree_state FOREIGN KEY (state) REFERENCES treesense."tree_state"(tree_state_id),
    CONSTRAINT fk_tree_created_by FOREIGN KEY (created_by) REFERENCES auth."user"(user_id),  
    CONSTRAINT fk_tree_updated_by FOREIGN KEY (updated_by) REFERENCES auth."user"(user_id)
);
 

COMMENT ON TABLE treesense."tree" IS 'Table storing scanned trees along different user routes';
COMMENT ON COLUMN treesense."tree".tree_id IS 'Unique identifier for the tree';
COMMENT ON COLUMN treesense."tree".route_id IS 'Reference to the route where the tree was scanned';
COMMENT ON COLUMN treesense."tree".species IS 'Species name of the tree';
COMMENT ON COLUMN treesense."tree".state IS 'State of the tree (e.g., healthy, sick, dry)';
COMMENT ON COLUMN treesense."tree".location IS 'Geographic location of the tree stored as a point (WGS 84 - SRID 4326)';
COMMENT ON COLUMN treesense."tree".age IS 'Approximate age of the tree in years';
COMMENT ON COLUMN treesense."tree".height IS 'Height of the tree in meters';
COMMENT ON COLUMN treesense."tree".diameter IS 'Diameter of the tree trunk in centimeters';
COMMENT ON COLUMN treesense."tree".description IS 'Additional information about the tree';
COMMENT ON COLUMN treesense."tree".created_at IS 'Timestamp of when the record was created';


-- ===============================================
-- Audit Schema: Schema for audit purpose
-- ===============================================

CREATE SCHEMA IF NOT EXISTS audit;

CREATE TABLE audit."activity_log" (
    activity_log_id UUID DEFAULT gen_random_uuid() UNIQUE NOT NULL,
    action_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id UUID,
    CONSTRAINT fk_activity_log_user_user FOREIGN KEY (user_id) REFERENCES auth."user"(user_id)
);

COMMENT ON TABLE audit."activity_log" IS 'Table of audit for activitys of users';
COMMENT ON COLUMN audit."activity_log".user_id IS 'Unique identifier for the user who made the action';
COMMENT ON COLUMN audit."activity_log".action_name IS 'Action Name';
