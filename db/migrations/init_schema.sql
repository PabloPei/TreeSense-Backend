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
    photo_url TEXT CHECK (photo_url ~* '^https?://.+') DEFAULT 'https://userphoto.png',
    language_code VARCHAR(10) DEFAULT 'es',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_auth_user_language FOREIGN KEY (language_code) REFERENCES conf.language(code)
);

COMMENT ON TABLE auth."user" IS 'Table of application users';
COMMENT ON COLUMN auth."user".user_id IS 'Unique identifier for the user';
COMMENT ON COLUMN auth."user".user_name IS 'Name of the user';
COMMENT ON COLUMN auth."user".photo_url IS 'URL of the user''s photo';
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
