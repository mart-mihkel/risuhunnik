CREATE TABLE IF NOT EXISTS jokes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    joke TEXT UNIQUE NOT NULL,
    tags TEXT,
    suggestion INTEGER DEFAULT 0
);

INSERT INTO jokes (joke, tags) VALUES
    ("Ühel mehel oli lagi nii madal, et ta sai ainult lesta süüa!", "tag1"),
    ("Kus elab Leedu hernes? Kauna sees", "tag1"),
    ("Tere siin Louvre, kuidas saab aidata? Kas Mona Lisa saab telefonile tulla? Ei saa, ta on maal", "tag1"),
    ("Puudega tudeng soovib üürida ahiküttega korterit.", "tag1"),
    ("Parim aeg majoneesi ostmiseks on aprill, sest siis mai on ees.", "tag1"),
    ("Meil on Kuusalust Gyproki vaja. Varustaja jõuab tagasi kuue aluse Gyprokiga.", "tag1"),
    ("Lambakarjus võitis lotoga. Mis ta esimene asjana ostis? Villa.", "tag1"),
    ("Mis sai rähn endale jõuludeks? Paku", "tag1"),
    ("Mees seisab mööda teed.", "tag1"),
    ("Kuidas kutsutakse ümarate jalgadega venelast? Oleg", "tag1"),
    ("Kui rottid on närilised siis sääsed on imelised", "tag1"),
    ("Miks me ütleme 'Ta sai peksa', miks me ei ütle 'Ta leib peksa'", "tag1");
