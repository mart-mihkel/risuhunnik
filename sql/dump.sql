CREATE TABLE IF NOT EXISTS conundrums  (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    text TEXT UNIQUE NOT NULL,
    tags TEXT NOT NULL,
    stars INTEGER NOT NULL DEFAULT 0,
    verified INTEGER NOT NULL DEFAULT 0
);

INSERT INTO conundrums (text, tags, verified) VALUES
    ("Ühel mehel oli lagi nii madal, et ta sai ainult lesta süüa!", "lest lagi mees", 1),
    ("Kui hull kukub aknast alla siis pole hullu!", "aken hull Valkla", 1),
    ("Mees seisab mööda teed.", "mees tee", 1),
    ("Kuri tigu, ta nimi oli kuritegu!", "tigu kuri", 1),
    ("Kas Mona Lisa saab telefonile tulla? Ei saa, ta on maal!", "maal naine", 1),
    ("Parim aeg majoneesi ostmiseks on aprill, sest mai on ees!", "majonees aprill", 1),
    ("Leedu hernes elab kauna sees!", "Kaunas hernes kaun", 1),
    ("Meil on Kuusalust Gyproki vaja! Varustaja jõuab tagasi kuue aluse Gyprokiga!", "Kuusalu kuus", 1),
    ("Saatan käis maal ja ta sõitis Cadillaciga", "saatan maal Cadillac", 1),
    ("Ümarate jalgadega venelane, ta nimi oli Oleg!", "mees jalg", 1);
