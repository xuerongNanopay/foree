package time_util

/*
Time zones available for :  linux
------------------------
Retrieving time zones from :  /usr/share/zoneinfo/
GMT+00:00 Africa/Abidjan
GMT+00:00 Africa/Accra
GMT+03:00 Africa/Addis_Ababa
GMT+01:00 Africa/Algiers
GMT+03:00 Africa/Asmara
GMT+03:00 Africa/Asmera
GMT+00:00 Africa/Bamako
GMT+01:00 Africa/Bangui
GMT+00:00 Africa/Banjul
GMT+00:00 Africa/Bissau
GMT+02:00 Africa/Blantyre
GMT+01:00 Africa/Brazzaville
GMT+02:00 Africa/Bujumbura
GMT+02:00 Africa/Cairo
GMT+00:00 Africa/Casablanca
GMT+01:00 Africa/Ceuta
GMT+00:00 Africa/Conakry
GMT+00:00 Africa/Dakar
GMT+03:00 Africa/Dar_es_Salaam
GMT+03:00 Africa/Djibouti
GMT+01:00 Africa/Douala
GMT+00:00 Africa/El_Aaiun
GMT+00:00 Africa/Freetown
GMT+02:00 Africa/Gaborone
GMT+02:00 Africa/Harare
GMT+02:00 Africa/Johannesburg
GMT+03:00 Africa/Juba
GMT+03:00 Africa/Kampala
GMT+03:00 Africa/Khartoum
GMT+02:00 Africa/Kigali
GMT+01:00 Africa/Kinshasa
GMT+01:00 Africa/Lagos
GMT+01:00 Africa/Libreville
GMT+00:00 Africa/Lome
GMT+01:00 Africa/Luanda
GMT+02:00 Africa/Lubumbashi
GMT+02:00 Africa/Lusaka
GMT+01:00 Africa/Malabo
GMT+02:00 Africa/Maputo
GMT+02:00 Africa/Maseru
GMT+02:00 Africa/Mbabane
GMT+03:00 Africa/Mogadishu
GMT+00:00 Africa/Monrovia
GMT+03:00 Africa/Nairobi
GMT+01:00 Africa/Ndjamena
GMT+01:00 Africa/Niamey
GMT+00:00 Africa/Nouakchott
GMT+00:00 Africa/Ouagadougou
GMT+01:00 Africa/Porto-Novo
GMT+00:00 Africa/Sao_Tome
GMT+00:00 Africa/Timbuktu
GMT+02:00 Africa/Tripoli
GMT+01:00 Africa/Tunis
GMT+02:00 Africa/Windhoek
GMT-10:00 America/Adak
GMT-09:00 America/Anchorage
GMT-04:00 America/Anguilla
GMT-04:00 America/Antigua
GMT-03:00 America/Araguaina
GMT-03:00 America/Argentina/Buenos_Aires
GMT-03:00 America/Argentina/Catamarca
GMT-03:00 America/Argentina/ComodRivadavia
GMT-03:00 America/Argentina/Cordoba
GMT-03:00 America/Argentina/Jujuy
GMT-03:00 America/Argentina/La_Rioja
GMT-03:00 America/Argentina/Mendoza
GMT-03:00 America/Argentina/Rio_Gallegos
GMT-03:00 America/Argentina/Salta
GMT-03:00 America/Argentina/San_Juan
GMT-03:00 America/Argentina/San_Luis
GMT-03:00 America/Argentina/Tucuman
GMT-03:00 America/Argentina/Ushuaia
GMT-04:00 America/Aruba
GMT-03:00 America/Asuncion
GMT-05:00 America/Atikokan
GMT-10:00 America/Atka
GMT-03:00 America/Bahia
GMT-07:00 America/Bahia_Banderas
GMT-04:00 America/Barbados
GMT-03:00 America/Belem
GMT-06:00 America/Belize
GMT-04:00 America/Blanc-Sablon
GMT-04:00 America/Boa_Vista
GMT-05:00 America/Bogota
GMT-07:00 America/Boise
GMT-03:00 America/Buenos_Aires
GMT-07:00 America/Cambridge_Bay
GMT-03:00 America/Campo_Grande
GMT-06:00 America/Cancun
GMT-04:30 America/Caracas
GMT-03:00 America/Catamarca
GMT-03:00 America/Cayenne
GMT-05:00 America/Cayman
GMT-06:00 America/Chicago
GMT-07:00 America/Chihuahua
GMT-05:00 America/Coral_Harbour
GMT-03:00 America/Cordoba
GMT-06:00 America/Costa_Rica
GMT-07:00 America/Creston
GMT-03:00 America/Cuiaba
GMT-04:00 America/Curacao
GMT+00:00 America/Danmarkshavn
GMT-08:00 America/Dawson
GMT-07:00 America/Dawson_Creek
GMT-07:00 America/Denver
GMT-05:00 America/Detroit
GMT-04:00 America/Dominica
GMT-07:00 America/Edmonton
GMT-04:00 America/Eirunepe
GMT-06:00 America/El_Salvador
GMT-08:00 America/Ensenada
GMT-08:00 America/Fort_Nelson
GMT-05:00 America/Fort_Wayne
GMT-03:00 America/Fortaleza
GMT-04:00 America/Glace_Bay
GMT-03:00 America/Godthab
GMT-04:00 America/Goose_Bay
GMT-05:00 America/Grand_Turk
GMT-04:00 America/Grenada
GMT-04:00 America/Guadeloupe
GMT-06:00 America/Guatemala
GMT-05:00 America/Guayaquil
GMT-04:00 America/Guyana
GMT-04:00 America/Halifax
GMT-05:00 America/Havana
GMT-07:00 America/Hermosillo
GMT-05:00 America/Indiana/Indianapolis
GMT-06:00 America/Indiana/Knox
GMT-05:00 America/Indiana/Marengo
GMT-05:00 America/Indiana/Petersburg
GMT-06:00 America/Indiana/Tell_City
GMT-05:00 America/Indiana/Vevay
GMT-05:00 America/Indiana/Vincennes
GMT-05:00 America/Indiana/Winamac
GMT-05:00 America/Indianapolis
GMT-07:00 America/Inuvik
GMT-05:00 America/Iqaluit
GMT-05:00 America/Jamaica
GMT-03:00 America/Jujuy
GMT-09:00 America/Juneau
GMT-05:00 America/Kentucky/Louisville
GMT-05:00 America/Kentucky/Monticello
GMT-06:00 America/Knox_IN
GMT-04:00 America/Kralendijk
GMT-04:00 America/La_Paz
GMT-05:00 America/Lima
GMT-08:00 America/Los_Angeles
GMT-05:00 America/Louisville
GMT-04:00 America/Lower_Princes
GMT-03:00 America/Maceio
GMT-06:00 America/Managua
GMT-04:00 America/Manaus
GMT-04:00 America/Marigot
GMT-04:00 America/Martinique
GMT-06:00 America/Matamoros
GMT-07:00 America/Mazatlan
GMT-03:00 America/Mendoza
GMT-06:00 America/Menominee
GMT-06:00 America/Merida
GMT-08:00 America/Metlakatla
GMT-06:00 America/Mexico_City
GMT-03:00 America/Miquelon
GMT-04:00 America/Moncton
GMT-06:00 America/Monterrey
GMT-02:00 America/Montevideo
GMT-05:00 America/Montreal
GMT-04:00 America/Montserrat
GMT-05:00 America/Nassau
GMT-05:00 America/New_York
GMT-05:00 America/Nipigon
GMT-09:00 America/Nome
GMT-02:00 America/Noronha
GMT-07:00 America/North_Dakota/Beulah
GMT-06:00 America/North_Dakota/Center
GMT-06:00 America/North_Dakota/New_Salem
GMT-07:00 America/Ojinaga
GMT-05:00 America/Panama
GMT-05:00 America/Pangnirtung
GMT-03:00 America/Paramaribo
GMT-07:00 America/Phoenix
GMT-05:00 America/Port-au-Prince
GMT-04:00 America/Port_of_Spain
GMT-04:00 America/Porto_Acre
GMT-04:00 America/Porto_Velho
GMT-04:00 America/Puerto_Rico
GMT-03:00 America/Punta_Arenas
GMT-06:00 America/Rainy_River
GMT-06:00 America/Rankin_Inlet
GMT-03:00 America/Recife
GMT-06:00 America/Regina
GMT-06:00 America/Resolute
GMT-04:00 America/Rio_Branco
GMT-03:00 America/Rosario
GMT-08:00 America/Santa_Isabel
GMT-03:00 America/Santarem
GMT-03:00 America/Santiago
GMT-04:00 America/Santo_Domingo
GMT-02:00 America/Sao_Paulo
GMT-01:00 America/Scoresbysund
GMT-07:00 America/Shiprock
GMT-09:00 America/Sitka
GMT-04:00 America/St_Barthelemy
GMT-03:30 America/St_Johns
GMT-04:00 America/St_Kitts
GMT-04:00 America/St_Lucia
GMT-04:00 America/St_Thomas
GMT-04:00 America/St_Vincent
GMT-06:00 America/Swift_Current
GMT-06:00 America/Tegucigalpa
GMT-04:00 America/Thule
GMT-05:00 America/Thunder_Bay
GMT-08:00 America/Tijuana
GMT-05:00 America/Toronto
GMT-04:00 America/Tortola
GMT-08:00 America/Vancouver
GMT-04:00 America/Virgin
GMT-08:00 America/Whitehorse
GMT-06:00 America/Winnipeg
GMT-09:00 America/Yakutat
GMT-07:00 America/Yellowknife
GMT+11:00 Antarctica/Casey
GMT+05:00 Antarctica/Davis
GMT+10:00 Antarctica/DumontDUrville
GMT+11:00 Antarctica/Macquarie
GMT+05:00 Antarctica/Mawson
GMT+13:00 Antarctica/McMurdo
GMT-03:00 Antarctica/Palmer
GMT-03:00 Antarctica/Rothera
GMT+13:00 Antarctica/South_Pole
GMT+03:00 Antarctica/Syowa
GMT+00:00 Antarctica/Troll
GMT+06:00 Antarctica/Vostok
GMT+01:00 Arctic/Longyearbyen
GMT+03:00 Asia/Aden
GMT+06:00 Asia/Almaty
GMT+02:00 Asia/Amman
GMT+12:00 Asia/Anadyr
GMT+05:00 Asia/Aqtau
GMT+05:00 Asia/Aqtobe
GMT+05:00 Asia/Ashgabat
GMT+05:00 Asia/Ashkhabad
GMT+05:00 Asia/Atyrau
GMT+03:00 Asia/Baghdad
GMT+03:00 Asia/Bahrain
GMT+04:00 Asia/Baku
GMT+07:00 Asia/Bangkok
GMT+06:00 Asia/Barnaul
GMT+02:00 Asia/Beirut
GMT+06:00 Asia/Bishkek
GMT+08:00 Asia/Brunei
GMT+05:30 Asia/Calcutta
GMT+09:00 Asia/Chita
GMT+08:00 Asia/Choibalsan
GMT+08:00 Asia/Chongqing
GMT+08:00 Asia/Chungking
GMT+05:30 Asia/Colombo
GMT+07:00 Asia/Dacca
GMT+02:00 Asia/Damascus
GMT+07:00 Asia/Dhaka
GMT+09:00 Asia/Dili
GMT+04:00 Asia/Dubai
GMT+05:00 Asia/Dushanbe
GMT+02:00 Asia/Famagusta
GMT+02:00 Asia/Gaza
GMT+08:00 Asia/Harbin
GMT+02:00 Asia/Hebron
GMT+07:00 Asia/Ho_Chi_Minh
GMT+08:00 Asia/Hong_Kong
GMT+07:00 Asia/Hovd
GMT+08:00 Asia/Irkutsk
GMT+02:00 Asia/Istanbul
GMT+07:00 Asia/Jakarta
GMT+09:00 Asia/Jayapura
GMT+02:00 Asia/Jerusalem
GMT+04:30 Asia/Kabul
GMT+12:00 Asia/Kamchatka
GMT+05:00 Asia/Karachi
GMT+06:00 Asia/Kashgar
GMT+05:45 Asia/Kathmandu
GMT+05:45 Asia/Katmandu
GMT+10:00 Asia/Khandyga
GMT+05:30 Asia/Kolkata
GMT+07:00 Asia/Krasnoyarsk
GMT+08:00 Asia/Kuala_Lumpur
GMT+08:00 Asia/Kuching
GMT+03:00 Asia/Kuwait
GMT+08:00 Asia/Macao
GMT+08:00 Asia/Macau
GMT+11:00 Asia/Magadan
GMT+08:00 Asia/Makassar
GMT+08:00 Asia/Manila
GMT+04:00 Asia/Muscat
GMT+02:00 Asia/Nicosia
GMT+07:00 Asia/Novokuznetsk
GMT+06:00 Asia/Novosibirsk
GMT+06:00 Asia/Omsk
GMT+05:00 Asia/Oral
GMT+07:00 Asia/Phnom_Penh
GMT+07:00 Asia/Pontianak
GMT+09:00 Asia/Pyongyang
GMT+03:00 Asia/Qatar
GMT+06:00 Asia/Qostanay
GMT+06:00 Asia/Qyzylorda
GMT+06:30 Asia/Rangoon
GMT+03:00 Asia/Riyadh
GMT+07:00 Asia/Saigon
GMT+10:00 Asia/Sakhalin
GMT+05:00 Asia/Samarkand
GMT+09:00 Asia/Seoul
GMT+08:00 Asia/Shanghai
GMT+08:00 Asia/Singapore
GMT+11:00 Asia/Srednekolymsk
GMT+08:00 Asia/Taipei
GMT+05:00 Asia/Tashkent
GMT+04:00 Asia/Tbilisi
GMT+03:30 Asia/Tehran
GMT+02:00 Asia/Tel_Aviv
GMT+06:00 Asia/Thimbu
GMT+06:00 Asia/Thimphu
GMT+09:00 Asia/Tokyo
GMT+06:00 Asia/Tomsk
GMT+08:00 Asia/Ujung_Pandang
GMT+08:00 Asia/Ulaanbaatar
GMT+08:00 Asia/Ulan_Bator
GMT+06:00 Asia/Urumqi
GMT+11:00 Asia/Ust-Nera
GMT+07:00 Asia/Vientiane
GMT+10:00 Asia/Vladivostok
GMT+09:00 Asia/Yakutsk
GMT+06:30 Asia/Yangon
GMT+05:00 Asia/Yekaterinburg
GMT+04:00 Asia/Yerevan
GMT-01:00 Atlantic/Azores
GMT-04:00 Atlantic/Bermuda
GMT+00:00 Atlantic/Canary
GMT-01:00 Atlantic/Cape_Verde
GMT+00:00 Atlantic/Faeroe
GMT+00:00 Atlantic/Faroe
GMT+01:00 Atlantic/Jan_Mayen
GMT+00:00 Atlantic/Madeira
GMT+00:00 Atlantic/Reykjavik
GMT-02:00 Atlantic/South_Georgia
GMT+00:00 Atlantic/St_Helena
GMT-03:00 Atlantic/Stanley
GMT+11:00 Australia/ACT
GMT+10:30 Australia/Adelaide
GMT+10:00 Australia/Brisbane
GMT+10:30 Australia/Broken_Hill
GMT+11:00 Australia/Canberra
GMT+11:00 Australia/Currie
GMT+09:30 Australia/Darwin
GMT+08:45 Australia/Eucla
GMT+11:00 Australia/Hobart
GMT+11:00 Australia/LHI
GMT+10:00 Australia/Lindeman
GMT+11:00 Australia/Lord_Howe
GMT+11:00 Australia/Melbourne
GMT+11:00 Australia/NSW
GMT+09:30 Australia/North
GMT+08:00 Australia/Perth
GMT+10:00 Australia/Queensland
GMT+10:30 Australia/South
GMT+11:00 Australia/Sydney
GMT+11:00 Australia/Tasmania
GMT+11:00 Australia/Victoria
GMT+08:00 Australia/West
GMT+10:30 Australia/Yancowinna
GMT-04:00 Brazil/Acre
GMT-02:00 Brazil/DeNoronha
GMT-02:00 Brazil/East
GMT-04:00 Brazil/West
GMT+01:00 CET
GMT-06:00 CST6CDT
GMT-04:00 Canada/Atlantic
GMT-06:00 Canada/Central
GMT-05:00 Canada/Eastern
GMT-07:00 Canada/Mountain
GMT-03:30 Canada/Newfoundland
GMT-08:00 Canada/Pacific
GMT-06:00 Canada/Saskatchewan
GMT-08:00 Canada/Yukon
GMT-03:00 Chile/Continental
GMT-05:00 Chile/EasterIsland
GMT-05:00 Cuba
GMT+02:00 EET
GMT-05:00 EST
GMT-05:00 EST5EDT
GMT+02:00 Egypt
GMT+00:00 Eire
GMT+00:00 Etc/GMT
GMT+00:00 Etc/GMT+0
GMT-01:00 Etc/GMT+1
GMT-10:00 Etc/GMT+10
GMT-11:00 Etc/GMT+11
GMT-12:00 Etc/GMT+12
GMT-02:00 Etc/GMT+2
GMT-03:00 Etc/GMT+3
GMT-04:00 Etc/GMT+4
GMT-05:00 Etc/GMT+5
GMT-06:00 Etc/GMT+6
GMT-07:00 Etc/GMT+7
GMT-08:00 Etc/GMT+8
GMT-09:00 Etc/GMT+9
GMT+00:00 Etc/GMT-0
GMT+01:00 Etc/GMT-1
GMT+10:00 Etc/GMT-10
GMT+11:00 Etc/GMT-11
GMT+12:00 Etc/GMT-12
GMT+13:00 Etc/GMT-13
GMT+14:00 Etc/GMT-14
GMT+02:00 Etc/GMT-2
GMT+03:00 Etc/GMT-3
GMT+04:00 Etc/GMT-4
GMT+05:00 Etc/GMT-5
GMT+06:00 Etc/GMT-6
GMT+07:00 Etc/GMT-7
GMT+08:00 Etc/GMT-8
GMT+09:00 Etc/GMT-9
GMT+00:00 Etc/GMT0
GMT+00:00 Etc/Greenwich
GMT+00:00 Etc/UCT
GMT+00:00 Etc/UTC
GMT+00:00 Etc/Universal
GMT+00:00 Etc/Zulu
GMT+01:00 Europe/Amsterdam
GMT+01:00 Europe/Andorra
GMT+03:00 Europe/Astrakhan
GMT+02:00 Europe/Athens
GMT+00:00 Europe/Belfast
GMT+01:00 Europe/Belgrade
GMT+01:00 Europe/Berlin
GMT+01:00 Europe/Bratislava
GMT+01:00 Europe/Brussels
GMT+02:00 Europe/Bucharest
GMT+01:00 Europe/Budapest
GMT+01:00 Europe/Busingen
GMT+02:00 Europe/Chisinau
GMT+01:00 Europe/Copenhagen
GMT+00:00 Europe/Dublin
GMT+01:00 Europe/Gibraltar
GMT+00:00 Europe/Guernsey
GMT+02:00 Europe/Helsinki
GMT+00:00 Europe/Isle_of_Man
GMT+02:00 Europe/Istanbul
GMT+00:00 Europe/Jersey
GMT+02:00 Europe/Kaliningrad
GMT+02:00 Europe/Kiev
GMT+03:00 Europe/Kirov
GMT+00:00 Europe/Lisbon
GMT+01:00 Europe/Ljubljana
GMT+00:00 Europe/London
GMT+01:00 Europe/Luxembourg
GMT+01:00 Europe/Madrid
GMT+01:00 Europe/Malta
GMT+02:00 Europe/Mariehamn
GMT+02:00 Europe/Minsk
GMT+01:00 Europe/Monaco
GMT+03:00 Europe/Moscow
GMT+02:00 Europe/Nicosia
GMT+01:00 Europe/Oslo
GMT+01:00 Europe/Paris
GMT+01:00 Europe/Podgorica
GMT+01:00 Europe/Prague
GMT+02:00 Europe/Riga
GMT+01:00 Europe/Rome
GMT+04:00 Europe/Samara
GMT+01:00 Europe/San_Marino
GMT+01:00 Europe/Sarajevo
GMT+03:00 Europe/Saratov
GMT+02:00 Europe/Simferopol
GMT+01:00 Europe/Skopje
GMT+02:00 Europe/Sofia
GMT+01:00 Europe/Stockholm
GMT+02:00 Europe/Tallinn
GMT+01:00 Europe/Tirane
GMT+02:00 Europe/Tiraspol
GMT+03:00 Europe/Ulyanovsk
GMT+02:00 Europe/Uzhgorod
GMT+01:00 Europe/Vaduz
GMT+01:00 Europe/Vatican
GMT+01:00 Europe/Vienna
GMT+02:00 Europe/Vilnius
GMT+03:00 Europe/Volgograd
GMT+01:00 Europe/Warsaw
GMT+01:00 Europe/Zagreb
GMT+02:00 Europe/Zaporozhye
GMT+01:00 Europe/Zurich
GMT+00:00 Factory
GMT+00:00 GB
GMT+00:00 GB-Eire
GMT+00:00 GMT
GMT+00:00 GMT+0
GMT+00:00 GMT-0
GMT+00:00 GMT0
GMT+00:00 Greenwich
GMT-10:00 HST
GMT+08:00 Hongkong
GMT+00:00 Iceland
GMT+03:00 Indian/Antananarivo
GMT+06:00 Indian/Chagos
GMT+07:00 Indian/Christmas
GMT+06:30 Indian/Cocos
GMT+03:00 Indian/Comoro
GMT+05:00 Indian/Kerguelen
GMT+04:00 Indian/Mahe
GMT+05:00 Indian/Maldives
GMT+04:00 Indian/Mauritius
GMT+03:00 Indian/Mayotte
GMT+04:00 Indian/Reunion
GMT+03:30 Iran
GMT+02:00 Israel
GMT-05:00 Jamaica
GMT+09:00 Japan
GMT+12:00 Kwajalein
GMT+02:00 Libya
GMT+01:00 MET
GMT-07:00 MST
GMT-07:00 MST7MDT
GMT-08:00 Mexico/BajaNorte
GMT-07:00 Mexico/BajaSur
GMT-06:00 Mexico/General
GMT+13:00 NZ
GMT+13:45 NZ-CHAT
GMT-07:00 Navajo
GMT+08:00 PRC
GMT-08:00 PST8PDT
GMT-11:00 Pacific/Apia
GMT+13:00 Pacific/Auckland
GMT+10:00 Pacific/Bougainville
GMT+13:45 Pacific/Chatham
GMT+10:00 Pacific/Chuuk
GMT-05:00 Pacific/Easter
GMT+11:00 Pacific/Efate
GMT+13:00 Pacific/Enderbury
GMT-11:00 Pacific/Fakaofo
GMT+12:00 Pacific/Fiji
GMT+12:00 Pacific/Funafuti
GMT-06:00 Pacific/Galapagos
GMT-09:00 Pacific/Gambier
GMT+11:00 Pacific/Guadalcanal
GMT+10:00 Pacific/Guam
GMT-10:00 Pacific/Honolulu
GMT-10:00 Pacific/Johnston
GMT+14:00 Pacific/Kiritimati
GMT+11:00 Pacific/Kosrae
GMT+12:00 Pacific/Kwajalein
GMT+12:00 Pacific/Majuro
GMT-09:30 Pacific/Marquesas
GMT-11:00 Pacific/Midway
GMT+12:00 Pacific/Nauru
GMT-11:00 Pacific/Niue
GMT+11:30 Pacific/Norfolk
GMT+11:00 Pacific/Noumea
GMT-11:00 Pacific/Pago_Pago
GMT+09:00 Pacific/Palau
GMT-08:00 Pacific/Pitcairn
GMT+11:00 Pacific/Pohnpei
GMT+11:00 Pacific/Ponape
GMT+10:00 Pacific/Port_Moresby
GMT-10:00 Pacific/Rarotonga
GMT+10:00 Pacific/Saipan
GMT-11:00 Pacific/Samoa
GMT-10:00 Pacific/Tahiti
GMT+12:00 Pacific/Tarawa
GMT+13:00 Pacific/Tongatapu
GMT+10:00 Pacific/Truk
GMT+12:00 Pacific/Wake
GMT+12:00 Pacific/Wallis
GMT+10:00 Pacific/Yap
GMT+01:00 Poland
GMT+00:00 Portugal
GMT+08:00 ROC
GMT+09:00 ROK
GMT+08:00 Singapore
GMT-04:00 SystemV/AST4
GMT-04:00 SystemV/AST4ADT
GMT-06:00 SystemV/CST6
GMT-06:00 SystemV/CST6CDT
GMT
*/