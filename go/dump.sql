CREATE TABLE IF NOT EXISTS jokes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    joke TEXT UNIQUE NOT NULL,
    tags TEXT NOT NULL
);

INSERT INTO jokes (joke, tags) VALUES
    ("Ühel mehel oli lagi nii madal, et ta sai ainult lesta süüa!", "lest,lagi,mees"),
    ("Kui hull kukub aknast alla siis pole hullu!", "aken,hull,Valkla"),
    ("Mees seisab mööda teed.", "mees,tee"),
    ("Kuri tigu, ta nimi oli kuritegu!", "tigu,kuri"),
    ("Kas Mona Lisa saab telefonile tulla? Ei saa, ta on maal!", "maal,naine"),
    ("Parim aeg majoneesi ostmiseks on aprill, sest mai on ees!", "majonees,aprill"),
    ("Leedu hernes elab kauna sees!", "Kaunas,hernes,kaun"),
    ("Meil on Kuusalust Gyproki vaja! Varustaja jõuab tagasi kuue aluse Gyprokiga!", "Kuusalu,kuus"),
    ("Saatan käis maal ja ta sõitis Cadillaciga", "saatan,maal,Cadillac"),
    ("Ümarate jalgadega venelane, ta nimi oli Oleg!", "mees,jalg");
