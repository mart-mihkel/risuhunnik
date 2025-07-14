CREATE TABLE IF NOT EXISTS jokes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    joke TEXT UNIQUE NOT NULL,
    tag TEXT UNIQUE NOT NULL,
    attr TEXT,
    suggestion INTEGER DEFAULT 0
);

INSERT INTO jokes (joke, tag, attr) VALUES
    (
        "Ühel mehel oli lagi nii madal, et ta sai ainult lesta süüa!",
        "lest",
        "Dead Flounder by Quaternius"
    ),
    (
        "Kus elab Leedu hernes? Kauna sees!",
        "hernes",
        "Peapod by Poly by Google [CC-BY] via Poly Pizza"
    ),
    (
        "Kas Mona Lisa saab telefonile tulla? Ei saa, ta on maal!",
        "maal",
        "Wall painting by jeremy [CC-BY] via Poly Pizza"
    ),
    (
        "Puudega tudeng soovib ahjuküttega korterit.",
        "puud",
        "Log & Axe - Game Asset by Don Carson [CC-BY] via Poly Pizza"
    ),
    (
        "Parim aeg majoneesi ostmiseks on aprill, sest mai on ees.",
        "majonees",
        NULL
    ),
    (
        "Meil on Kuusalust Gyproki vaja. Varustaja jõuab tagasi kuue aluse Gyprokiga.",
        "gyprok",
        NULL
    ),
    (
        "Mille ostis lambakarjus peale lotovõitu? Villa!",
        "lambakarjus",
        NULL
    ),
    (
        "Mille sai rähn jõuludeks? Paku!",
        "rähn",
        NULL
    ),
    (
        "Mees seisab mööda teed.",
        "teed",
        NULL
    ),
    (
        "Kuidas kutsutakse ümarate jalgadega venelast? Oleg!",
        "oleg",
        NULL
    ),
    (
        "Kui rottid on närilised siis sääsed on imelised",
        "imelised",
        NULL
    );
