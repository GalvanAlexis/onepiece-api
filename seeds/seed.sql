-- ===========================
-- ONE PIECE API: Seed Data
-- ===========================

-- ===== CREWS =====
INSERT INTO crews (name, captain, affiliation) VALUES
    ('Straw Hat Pirates',    'Monkey D. Luffy',   'Straw Hat Grand Fleet'),
    ('Whitebeard Pirates',   'Edward Newgate',     'Yonko'),
    ('Heart Pirates',        'Trafalgar D. Water Law', NULL),
    ('Red Hair Pirates',     'Shanks',             'Yonko'),
    ('Donquixote Pirates',   'Donquixote Doflamingo', 'Shichibukai / Underworld')
ON CONFLICT (name) DO NOTHING;

-- ===== ARCS =====
INSERT INTO arcs (name, saga, description, order_index) VALUES
    ('East Blue Saga',        'East Blue',      'Luffy assembles his crew in the East Blue.',                1),
    ('Alabasta Arc',          'Alabasta Saga',  'Luffy battles Crocodile to save the Kingdom of Alabasta.', 2),
    ('Marineford Arc',        'Summit War',     'The war at Marineford to save Portgas D. Ace.',             3),
    ('Dressrosa Arc',         'Dressrosa Saga', 'Luffy takes down Doflamingo in the Kingdom of Dressrosa.', 4),
    ('Wano Country Arc',      'Wano Saga',      'Luffy challenges Kaido in the isolated country of Wano.',  5)
ON CONFLICT (name) DO NOTHING;

-- ===== DEVIL FRUITS =====
INSERT INTO devil_fruits (name, type, description, current_user) VALUES
    ('Gomu Gomu no Mi',    'Paramecia', 'Grants the user a rubber body.',                                      'Monkey D. Luffy'),
    ('Mera Mera no Mi',    'Logia',     'Grants the user the power to create and control fire.',               NULL),
    ('Gura Gura no Mi',    'Paramecia', 'Grants the user the ability to create quakes.',                       NULL),
    ('Ope Ope no Mi',      'Paramecia', 'Grants the user the ability to create a "Room" and manipulate it.',   'Trafalgar D. Water Law'),
    ('Yami Yami no Mi',    'Logia',     'Grants the power to control darkness and gravity.',                   'Marshall D. Teach'),
    ('Hana Hana no Mi',    'Paramecia', 'Allows the user to replicate body parts on any surface.',             'Nico Robin'),
    ('Hito Hito no Mi',    'Zoan',      'Grants a reindeer the ability to transform into a human.',            'Tony Tony Chopper'),
    ('Suke Suke no Mi',    'Paramecia', 'Grants the user and anything they touch the ability to turn invisible.', NULL),
    ('Fuwa Fuwa no Mi',    'Paramecia', 'Grants the user the ability to make non-living things float.',         NULL),
    ('Mochi Mochi no Mi',  'Paramecia', 'Special Paramecia that allows the user to create, control and transform into mochi.', 'Charlotte Katakuri')
ON CONFLICT (name) DO NOTHING;

-- ===== CHARACTERS =====
INSERT INTO characters (name, alias, status, bounty, origin, description, crew_id, devil_fruit_id) VALUES
    (
        'Monkey D. Luffy', 'Straw Hat Luffy', 'alive', 3000000000,
        'Foosha Village, East Blue',
        'The captain of the Straw Hat Pirates and the main protagonist of One Piece.',
        (SELECT id FROM crews WHERE name = 'Straw Hat Pirates'),
        (SELECT id FROM devil_fruits WHERE name = 'Gomu Gomu no Mi')
    ),
    (
        'Roronoa Zoro', 'Pirate Hunter Zoro', 'alive', 1111000000,
        'Shimotsuki Village, East Blue',
        'The swordsman of the Straw Hat Pirates, aiming to become the world''s greatest swordsman.',
        (SELECT id FROM crews WHERE name = 'Straw Hat Pirates'),
        NULL
    ),
    (
        'Nami', 'Cat Burglar Nami', 'alive', 366000000,
        'Cocoyasi Village, East Blue',
        'The navigator of the Straw Hat Pirates and an expert thief and pickpocket.',
        (SELECT id FROM crews WHERE name = 'Straw Hat Pirates'),
        NULL
    ),
    (
        'Usopp', 'God Usopp', 'alive', 500000000,
        'Syrup Village, East Blue',
        'The sniper of the Straw Hat Pirates and son of the pirate Yasopp.',
        (SELECT id FROM crews WHERE name = 'Straw Hat Pirates'),
        NULL
    ),
    (
        'Sanji', 'Black Leg Sanji', 'alive', 1032000000,
        'North Blue',
        'The cook of the Straw Hat Pirates and son of Vinsmoke Judge.',
        (SELECT id FROM crews WHERE name = 'Straw Hat Pirates'),
        NULL
    ),
    (
        'Tony Tony Chopper', 'Cotton Candy Lover Chopper', 'alive', 1000,
        'Drum Island, Grand Line',
        'The doctor of the Straw Hat Pirates and a reindeer who ate the Human-Human Fruit.',
        (SELECT id FROM crews WHERE name = 'Straw Hat Pirates'),
        (SELECT id FROM devil_fruits WHERE name = 'Hito Hito no Mi')
    ),
    (
        'Nico Robin', 'Devil Child Robin', 'alive', 930000000,
        'Ohara, West Blue',
        'The archaeologist of the Straw Hat Pirates, the only survivor of Ohara.',
        (SELECT id FROM crews WHERE name = 'Straw Hat Pirates'),
        (SELECT id FROM devil_fruits WHERE name = 'Hana Hana no Mi')
    ),
    (
        'Franky', 'Cyborg Franky', 'alive', 394000000,
        'South Blue',
        'The shipwright of the Straw Hat Pirates and a self-modified cyborg.',
        (SELECT id FROM crews WHERE name = 'Straw Hat Pirates'),
        NULL
    ),
    (
        'Brook', '"Soul King" Brook', 'alive', 383000000,
        'West Blue',
        'The musician of the Straw Hat Pirates, a living skeleton thanks to the Yomi Yomi no Mi.',
        (SELECT id FROM crews WHERE name = 'Straw Hat Pirates'),
        NULL
    ),
    (
        'Jinbe', 'Knight of the Sea Jinbe', 'alive', 1100000000,
        'Fish-Man Island',
        'The helmsman of the Straw Hat Pirates and a Fish-Man martial arts master.',
        (SELECT id FROM crews WHERE name = 'Straw Hat Pirates'),
        NULL
    ),
    (
        'Trafalgar D. Water Law', 'The Surgeon of Death', 'alive', 3000000000,
        'Flevance, North Blue',
        'The captain of the Heart Pirates, a Warlord of the Sea turned Revolutionary ally.',
        (SELECT id FROM crews WHERE name = 'Heart Pirates'),
        (SELECT id FROM devil_fruits WHERE name = 'Ope Ope no Mi')
    ),
    (
        'Shanks', 'Red-Haired Shanks', 'alive', NULL,
        'Grand Line',
        'One of the Four Emperors of the Sea and the captain of the Red Hair Pirates.',
        (SELECT id FROM crews WHERE name = 'Red Hair Pirates'),
        NULL
    ),
    (
        'Edward Newgate', 'Whitebeard', 'deceased', NULL,
        'Sphinx Island',
        'Formerly one of the Four Emperors of the Sea, called the strongest man in the world.',
        (SELECT id FROM crews WHERE name = 'Whitebeard Pirates'),
        (SELECT id FROM devil_fruits WHERE name = 'Gura Gura no Mi')
    ),
    (
        'Portgas D. Ace', 'Fire Fist Ace', 'deceased', 550000000,
        'Baterilla, South Blue',
        'Second commander of the Whitebeard Pirates and son of the Pirate King Gol D. Roger.',
        (SELECT id FROM crews WHERE name = 'Whitebeard Pirates'),
        (SELECT id FROM devil_fruits WHERE name = 'Mera Mera no Mi')
    ),
    (
        'Marshall D. Teach', 'Blackbeard', 'alive', NULL,
        'Unknown',
        'One of the Four Emperors of the Sea and the only person known to possess two Devil Fruit powers.',
        NULL,
        (SELECT id FROM devil_fruits WHERE name = 'Yami Yami no Mi')
    )
ON CONFLICT DO NOTHING;

-- ===== CHARACTER HAKI =====
-- Luffy: all 3
INSERT INTO character_haki (character_id, haki_id)
SELECT c.id, h.id FROM characters c, haki_types h
WHERE c.name = 'Monkey D. Luffy'
ON CONFLICT DO NOTHING;

-- Zoro: Busoshoku + Kenbunshoku
INSERT INTO character_haki (character_id, haki_id)
SELECT c.id, h.id FROM characters c, haki_types h
WHERE c.name = 'Roronoa Zoro' AND h.name IN ('Kenbunshoku Haki', 'Busoshoku Haki')
ON CONFLICT DO NOTHING;

-- Law: Busoshoku
INSERT INTO character_haki (character_id, haki_id)
SELECT c.id, h.id FROM characters c, haki_types h
WHERE c.name = 'Trafalgar D. Water Law' AND h.name = 'Busoshoku Haki'
ON CONFLICT DO NOTHING;

-- Shanks: all 3
INSERT INTO character_haki (character_id, haki_id)
SELECT c.id, h.id FROM characters c, haki_types h
WHERE c.name = 'Shanks'
ON CONFLICT DO NOTHING;

-- Newgate: Haoshoku + Busoshoku
INSERT INTO character_haki (character_id, haki_id)
SELECT c.id, h.id FROM characters c, haki_types h
WHERE c.name = 'Edward Newgate' AND h.name IN ('Haoshoku Haki', 'Busoshoku Haki')
ON CONFLICT DO NOTHING;

-- Ace: Kenbunshoku + Busoshoku
INSERT INTO character_haki (character_id, haki_id)
SELECT c.id, h.id FROM characters c, haki_types h
WHERE c.name = 'Portgas D. Ace' AND h.name IN ('Kenbunshoku Haki', 'Busoshoku Haki')
ON CONFLICT DO NOTHING;

-- Blackbeard: none (confirmed no Haki in early arcs)
