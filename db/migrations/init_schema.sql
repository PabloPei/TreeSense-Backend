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
    role_id VARCHAR(1) PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE auth.role IS 'Table of defined roles in the application';
COMMENT ON COLUMN auth.role.role_id IS 'Unique identifier for the role';
COMMENT ON COLUMN auth.role.role_name IS 'Name of the role';
COMMENT ON COLUMN auth.role.description IS 'Description of the role';

INSERT INTO auth.role (role_id, role_name, description)
VALUES
    ('V', 'VIEWER', 'User with read-only access'),
    ('E', 'EDITOR', 'User with editing capabilities'),
    ('A', 'ADMIN', 'User with full administrative privileges');


CREATE TABLE auth.user_role (
    user_id UUID,
    role_id VARCHAR(1),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID,
    valid_until TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_user_role_user FOREIGN KEY (user_id) REFERENCES auth."user"(user_id),
    CONSTRAINT fk_user_role_role FOREIGN KEY (role_id) REFERENCES auth.role(role_id)
);

COMMENT ON TABLE auth.user_role IS 'Table for assigning roles to users';
COMMENT ON COLUMN auth.user_role.user_id IS 'Identifier of the user';
COMMENT ON COLUMN auth.user_role.role_id IS 'Identifier of the assigned role';
COMMENT ON COLUMN auth.user_role.valid_until IS 'Date until the assignment is valid';

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


CREATE TABLE treesense."tree_specie" (
    tree_specie_id VARCHAR(10) PRIMARY KEY,  
    name VARCHAR(100) NOT NULL,  
    description TEXT  
);

COMMENT ON TABLE treesense."tree_specie" IS 'Table storing different tree species';
COMMENT ON COLUMN treesense."tree_specie".tree_specie_id IS 'Unique identifier for the tree species (code)';
COMMENT ON COLUMN treesense."tree_specie".name IS 'Common name of the tree species';
COMMENT ON COLUMN treesense."tree_specie".description IS 'Additional information about the species';

INSERT INTO treesense."tree_specie" (tree_specie_id, name, description) VALUES
('QRO', 'Quercus robur', 'Commonly known as English oak, native to Europe'),
('PIN', 'Pinus sylvestris', 'Scots pine, widely distributed across Eurasia'),
('ACR', 'Acer rubrum', 'Red maple, native to North America');



CREATE TABLE treesense."tree_state" (
    tree_state_id VARCHAR(10) PRIMARY KEY,  
    name VARCHAR(50) NOT NULL,  
    description TEXT  
);

COMMENT ON TABLE treesense."tree_state" IS 'Table storing different health states of trees';
COMMENT ON COLUMN treesense."tree_state".tree_state_id IS 'Unique identifier for the tree state';
COMMENT ON COLUMN treesense."tree_state".name IS 'Name of the tree state';
COMMENT ON COLUMN treesense."tree_state".description IS 'Additional information about the state';

INSERT INTO treesense."tree_state" (tree_state_id, name, description) VALUES
('HLT', 'Healthy', 'Tree is in good condition with no visible issues'),
('SCK', 'Sick', 'Tree shows signs of disease or infestation'),
('DRY', 'Dry', 'Tree appears to be dry or dying');


CREATE TABLE treesense."tree" (
    tree_id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    route_id UUID NOT NULL,
    specie VARCHAR(3),
    state VARCHAR(3),
    location GEOMETRY(Point, 4326) NOT NULL,
    antique INT,
    height INT,
    diameter INT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_tree_route FOREIGN KEY (route_id) REFERENCES treesense."route"(route_id),
    CONSTRAINT fk_tree_specie FOREIGN KEY (specie) REFERENCES treesense."tree_specie"(tree_specie_id),
    CONSTRAINT fk_tree_state FOREIGN KEY (state) REFERENCES treesense."tree_state"(tree_state_id)
);
 

COMMENT ON TABLE treesense."tree" IS 'Table storing scanned trees along different user routes';
COMMENT ON COLUMN treesense."tree".tree_id IS 'Unique identifier for the tree';
COMMENT ON COLUMN treesense."tree".route_id IS 'Reference to the route where the tree was scanned';
COMMENT ON COLUMN treesense."tree".specie IS 'Species code of the tree (check if 3 characters are enough)';
COMMENT ON COLUMN treesense."tree".state IS 'State of the tree (e.g., healthy, sick, dry)';
COMMENT ON COLUMN treesense."tree".location IS 'Geographic location of the tree stored as a point (WGS 84 - SRID 4326)';
COMMENT ON COLUMN treesense."tree".antique IS 'Approximate age of the tree in years';
COMMENT ON COLUMN treesense."tree".height IS 'Height of the tree in meters';
COMMENT ON COLUMN treesense."tree".diameter IS 'Diameter of the tree trunk in centimeters';
COMMENT ON COLUMN treesense."tree".description IS 'Additional information about the tree';
COMMENT ON COLUMN treesense."tree".created_at IS 'Timestamp of when the record was created';
