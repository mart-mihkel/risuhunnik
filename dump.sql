CREATE TABLE IF NOT EXISTS jokes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    joke TEXT UNIQUE NOT NULL
);

INSERT INTO jokes (joke) VALUES
    ("Ühel mehel oli lagi nii madal, et ta sai ainult lesta süüa!"),
    ("Kus elab Leedu hernes? Kauna sees!"),
    ("Kas Mona Lisa saab telefonile tulla? Ei saa, ta on maal!"),
    ("Parim aeg majoneesi ostmiseks on aprill! Sest mai on ees!"),
    ("Kuri tigu, ta nimi oli kuritegu!"),
    ("Meil on Kuusalust Gyproki vaja! Varustaja jõuab tagasi kuue aluse Gyprokiga!"),
    ("Mees seisab mööda teed."),
    ("Kuidas kutsutakse ümarate jalgadega venelast? Oleg!");
