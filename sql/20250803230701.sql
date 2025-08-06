BEGIN TRANSACTION;

CREATE TABLE authors(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    author TEXT UNIQUE NOT NULL
);

INSERT INTO authors (author) VALUES
    ('Juhan Juustuburger'),
    ('Mihkel Majameister'),
    ('Simon Segane'),
    ('Peeter Pannkook'),
    ('Artur Atsetoon'),
    ('Kalev Krõps'),
    ('Märt Makaron'),
    ('Tiina Tiramisu'),
    ('Eino Espresso'),
    ('Riho Ravioli'),
    ('Heli Hotdog'),
    ('Anu Avokaado'),
    ('Andreas Akudrell'),
    ('Kermo Kahveltõstuk'),
    ('Kalle Kahtlane'),
    ('Joonas Jäätis'),
    ('Signe Soolõhe'),
    ('Mauno Muffin'),
    ('Keiu Kookos'),
    ('Oskar Omelett'),
    ('Kadri Kartul'),
    ('Jaagup Juustukook'),
    ('Gustav Geenius'),
    ('Rita Röst'),
    ('Heikki Hakkliha'),
    ('Liisa Lattapannkook'),
    ('Andres Apelsin'),
    ('Marje Majonees'),
    ('Ulvi Uba'),
    ('Linda Lihapall'),
    ('Lembit Lest'),
    ('Elina Eelroog'),
    ('Tanel Tomat'),
    ('Gregor Granaat');

COMMIT;
