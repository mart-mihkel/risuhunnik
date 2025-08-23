BEGIN TRANSACTION;

INSERT INTO conundrums (text, author, verified, stars, date) VALUES
    ('Ühel mehel oli lagi nii madal, et ta sai ainult lesta süüa!', 'Kermo Kahveltõstuk', 1, 6, '2025-07-29 12:23:13'),
    ('Kui hull kukub aknast alla siis pole hullu!', 'Juhan Juustuburger', 1, 0, '2025-07-29 12:23:13'),
    ('Mees seisab mööda teed.', 'Juhan Juustuburger', 1, 2, '2025-07-29 12:23:13'),
    ('Kuri tigu, ta nimi oli kuritegu!', 'Juhan Juustuburger', 1, 0, '2025-07-29 12:23:13'),
    ('Kas Mona Lisa saab telefonile tulla? Ei saa, ta on maal!', 'Juhan Juustuburger', 1, 0, '2025-07-29 12:23:13'),
    ('Parim aeg majoneesi ostmiseks on aprill, sest mai on ees!', 'Juhan Juustuburger', 1, 0, '2025-07-29 12:23:13'),
    ('Leedu hernes elab kauna sees!', 'Juhan Juustuburger', 1, 0, '2025-07-29 12:23:13'),
    ('Kala sureb lahtise suuga!', 'Juhan Juustuburger', 0, 0, '2025-07-29 12:23:13'),
    ('Meil on Kuusalust Gyproki vaja! Varustaja jõuab tagasi kuue aluse Gyprokiga!', 'Juhan Juustuburger', 1, 0, '2025-07-29 12:23:13'),
    ('Saatan käis maal ja ta sõitis Cadillaciga!', 'Juhan Juustuburger', 0, 0, '2025-07-29 12:23:13'),
    (''';DROP TABLE conundrums;--', 'Gustav Geenius', 0, 0, '2025-07-29 12:23:13'),
    ('Ümarate jalgadega venelane, ta nimi oli Oleg!', 'Juhan Juustuburger', 1, 0, '2025-07-29 12:23:13'),
    ('Näljane mõtleb ikka leivast!', 'Artur Atsetoon', 0, 0, '2025-07-29 12:23:13'),
    ('Magavale hiirele kass suhu ei jookse', 'Andreas Akudrell', 0, 0, '2025-07-29 12:23:13'),
    ('Miks läks kann üle tee? Ta oli teekann!', 'Mihkel Majameister', 1, 0, '2025-07-29 12:23:13'),
    ('Nüüd ma olen küll endale karuteene teinud!', 'Juhan Juustuburger', 0, 0, '2025-07-30 13:18:16'),
    ('Väljas tuiskas nii kõvasti, et Kingpool ja Sammalhabe ei näinud Muhvigi.', 'Kermo Kahveltõstuk', 1, 0, '2025-08-01 16:39:24'),
    ('Mis on Aserbaidźaani pealinn?', 'Kermo Kahveltõstuk', 1, 0, '2025-08-01 17:27:31'),
    ('Piknikul meepoti kaotamine on paras karuteene!', 'Kermo Kahveltõstuk', 1, 1, '2025-08-01 23:10:43'),
    ('ACDC kontserti tehnika oli peamiselt heavy metal!', 'Kermo Kahveltõstuk', 1, 0, '2025-08-01 23:22:01'),
    ('Sokide puudumisel tuleks teavitada loomaaeda.', 'Kermo Kahveltõstuk', 0, 0, '2025-08-01 23:49:31'),
    ('Vesi läks vihast keema!', 'Janar Kootimum', 1, 0, '2025-08-15 23:19:31'),
    ('Toitlustus kompleks? Mulle paistab küll lihtne!', 'Juhan Juustuburger', 1, 0, '2025-08-25 21:42:11'),
    ('The man who stole my magnets started spinning when I pressed charges against him!', 'Juhan Juustuburger', 1, 0, '2025-08-22 12:31:15'),
    ('There''s been an ''rm -rf'' incident, oopsie daisy haha!', 'Mihkel Majameister', 0, 0, '2025-08-24 12:17:15');

INSERT INTO comments (cid, comment, author, date) VALUES
    (11, 'Idikas!', 'Juhan Juustuburger', '2025-07-29 12:23:13'),
    (1, 'Wow! Suurepärane!', 'Juhan Juustuburger','2025-07-29 15:21:27'),
    (17, 'Haha, päris hea Fred!', 'Juhan Juustuburger', '2025-08-01 16:39:52'),
    (17, 'Kes on Fred?', 'Kermo Kahveltõstuk', '2025-08-01 16:40:06'),
    (17, 'Mihkel miks sa ei vasta', 'Kermo Kahveltõstuk', '2025-08-01 17:26:58'),
    (18, 'Paku...', 'Kermo Kahveltõstuk', '2025-08-01 17:27:43'),
    (17, 'Piim läks halvaks, pidin keskenduma!', 'Juhan Juustuburger', '2025-08-02 13:22:27');

COMMIT;
