package timezone

import (
	"fmt"
	"strings"
	"time"
)

// Location represents a timezone location.
type Location struct {
	Name     string
	IanaName string
	Offset   int // offset in hours from UTC
}

// System manages timezone operations.
type System struct {
	locations map[string]*Location
}

// NewSystem creates a new timezone system.
func NewSystem() *System {
	s := &System{
		locations: make(map[string]*Location),
	}
	s.initLocations()
	return s
}


func (s *System) initLocations() {
	// Comprehensive world locations including countries and major cities with IANA time zones
	var locations = []Location{
		// ===== COUNTRIES (UN members + extras) =====
		
		// Afghanistan
		{"Afghanistan", "Asia/Kabul", 4},
		{"Kabul", "Asia/Kabul", 4},
		
		// Albania
		{"Albania", "Europe/Tirane", 1},
		{"Tirana", "Europe/Tirane", 1},
		
		// Algeria
		{"Algeria", "Africa/Algiers", 1},
		{"Algiers", "Africa/Algiers", 1},
		
		// Andorra
		{"Andorra", "Europe/Andorra", 1},
		{"Andorra la Vella", "Europe/Andorra", 1},
		
		// Angola
		{"Angola", "Africa/Luanda", 1},
		{"Luanda", "Africa/Luanda", 1},
		
		// Antigua and Barbuda
		{"Antigua and Barbuda", "America/Antigua", -4},
		{"Saint John's", "America/Antigua", -4},
		
		// Argentina
		{"Argentina", "America/Argentina/Buenos_Aires", -3},
		{"Buenos Aires", "America/Argentina/Buenos_Aires", -3},
		{"Cordoba", "America/Argentina/Cordoba", -3},
		{"Rosario", "America/Argentina/Buenos_Aires", -3},
		{"Mendoza", "America/Argentina/Mendoza", -3},
		
		// Armenia
		{"Armenia", "Asia/Yerevan", 4},
		{"Yerevan", "Asia/Yerevan", 4},
		
		// Australia
		{"Australia", "Australia/Sydney", 10},
		{"Canberra", "Australia/Sydney", 10},
		{"Sydney", "Australia/Sydney", 10},
		{"Melbourne", "Australia/Melbourne", 10},
		{"Brisbane", "Australia/Brisbane", 10},
		{"Perth", "Australia/Perth", 8},
		{"Adelaide", "Australia/Adelaide", 9},
		{"Hobart", "Australia/Hobart", 10},
		{"Darwin", "Australia/Darwin", 9},
		
		// Austria
		{"Austria", "Europe/Vienna", 1},
		{"Vienna", "Europe/Vienna", 1},
		{"Graz", "Europe/Vienna", 1},
		{"Salzburg", "Europe/Vienna", 1},
		
		// Azerbaijan
		{"Azerbaijan", "Asia/Baku", 4},
		{"Baku", "Asia/Baku", 4},
		
		// Bahamas
		{"Bahamas", "America/Nassau", -5},
		{"Nassau", "America/Nassau", -5},
		
		// Bahrain
		{"Bahrain", "Asia/Bahrain", 3},
		{"Manama", "Asia/Bahrain", 3},
		
		// Bangladesh
		{"Bangladesh", "Asia/Dhaka", 6},
		{"Dhaka", "Asia/Dhaka", 6},
		{"Chittagong", "Asia/Dhaka", 6},
		
		// Barbados
		{"Barbados", "America/Barbados", -4},
		{"Bridgetown", "America/Barbados", -4},
		
		// Belarus
		{"Belarus", "Europe/Minsk", 3},
		{"Minsk", "Europe/Minsk", 3},
		
		// Belgium
		{"Belgium", "Europe/Brussels", 1},
		{"Brussels", "Europe/Brussels", 1},
		{"Antwerp", "Europe/Brussels", 1},
		{"Bruges", "Europe/Brussels", 1},
		
		// Belize
		{"Belize", "America/Belize", -6},
		{"Belmopan", "America/Belize", -6},
		{"Belize City", "America/Belize", -6},
		
		// Benin
		{"Benin", "Africa/Porto-Novo", 1},
		{"Porto-Novo", "Africa/Porto-Novo", 1},
		{"Cotonou", "Africa/Porto-Novo", 1},
		
		// Bhutan
		{"Bhutan", "Asia/Thimphu", 6},
		{"Thimphu", "Asia/Thimphu", 6},
		
		// Bolivia
		{"Bolivia", "America/La_Paz", -4},
		{"La Paz", "America/La_Paz", -4},
		{"Sucre", "America/La_Paz", -4},
		{"Santa Cruz", "America/La_Paz", -4},
		
		// Bosnia and Herzegovina
		{"Bosnia and Herzegovina", "Europe/Sarajevo", 1},
		{"Sarajevo", "Europe/Sarajevo", 1},
		
		// Botswana
		{"Botswana", "Africa/Gaborone", 2},
		{"Gaborone", "Africa/Gaborone", 2},
		
		// Brazil
		{"Brazil", "America/Sao_Paulo", -3},
		{"Brasilia", "America/Sao_Paulo", -3},
		{"Sao Paulo", "America/Sao_Paulo", -3},
		{"Rio de Janeiro", "America/Sao_Paulo", -3},
		{"Salvador", "America/Bahia", -3},
		{"Recife", "America/Recife", -3},
		{"Belo Horizonte", "America/Sao_Paulo", -3},
		{"Fortaleza", "America/Fortaleza", -3},
		{"Manaus", "America/Manaus", -4},
		
		// Brunei
		{"Brunei", "Asia/Brunei", 8},
		{"Bandar Seri Begawan", "Asia/Brunei", 8},
		
		// Bulgaria
		{"Bulgaria", "Europe/Sofia", 2},
		{"Sofia", "Europe/Sofia", 2},
		{"Plovdiv", "Europe/Sofia", 2},
		
		// Burkina Faso
		{"Burkina Faso", "Africa/Ouagadougou", 0},
		{"Ouagadougou", "Africa/Ouagadougou", 0},
		
		// Burundi
		{"Burundi", "Africa/Bujumbura", 2},
		{"Bujumbura", "Africa/Bujumbura", 2},
		{"Gitega", "Africa/Bujumbura", 2},
		
		// Cabo Verde
		{"Cabo Verde", "Atlantic/Cape_Verde", -1},
		{"Praia", "Atlantic/Cape_Verde", -1},
		
		// Cambodia
		{"Cambodia", "Asia/Phnom_Penh", 7},
		{"Phnom Penh", "Asia/Phnom_Penh", 7},
		{"Siem Reap", "Asia/Phnom_Penh", 7},
		
		// Cameroon
		{"Cameroon", "Africa/Douala", 1},
		{"Yaounde", "Africa/Douala", 1},
		{"Douala", "Africa/Douala", 1},
		
		// Canada
		{"Canada", "America/Toronto", -5},
		{"Ottawa", "America/Toronto", -5},
		{"Toronto", "America/Toronto", -5},
		{"Montreal", "America/Toronto", -5},
		{"Vancouver", "America/Vancouver", -8},
		{"Calgary", "America/Edmonton", -7},
		{"Edmonton", "America/Edmonton", -7},
		{"Winnipeg", "America/Winnipeg", -6},
		{"Quebec City", "America/Toronto", -5},
		{"Halifax", "America/Halifax", -4},
		{"St John's", "America/St_Johns", -3},
		
		// Central African Republic
		{"Central African Republic", "Africa/Bangui", 1},
		{"Bangui", "Africa/Bangui", 1},
		
		// Chad
		{"Chad", "Africa/Ndjamena", 1},
		{"N'Djamena", "Africa/Ndjamena", 1},
		
		// Chile
		{"Chile", "America/Santiago", -4},
		{"Santiago", "America/Santiago", -4},
		{"Valparaiso", "America/Santiago", -4},
		{"Concepcion", "America/Santiago", -4},
		
		// China
		{"China", "Asia/Shanghai", 8},
		{"Beijing", "Asia/Shanghai", 8},
		{"Shanghai", "Asia/Shanghai", 8},
		{"Guangzhou", "Asia/Shanghai", 8},
		{"Shenzhen", "Asia/Shanghai", 8},
		{"Chengdu", "Asia/Shanghai", 8},
		{"Wuhan", "Asia/Shanghai", 8},
		{"Chongqing", "Asia/Shanghai", 8},
		{"Tianjin", "Asia/Shanghai", 8},
		{"Nanjing", "Asia/Shanghai", 8},
		{"Xi'an", "Asia/Shanghai", 8},
		
		// Colombia
		{"Colombia", "America/Bogota", -5},
		{"Bogota", "America/Bogota", -5},
		{"Medellin", "America/Bogota", -5},
		{"Cali", "America/Bogota", -5},
		{"Barranquilla", "America/Bogota", -5},
		
		// Comoros
		{"Comoros", "Indian/Comoro", 3},
		{"Moroni", "Indian/Comoro", 3},
		
		// Congo
		{"Congo", "Africa/Brazzaville", 1},
		{"Brazzaville", "Africa/Brazzaville", 1},
		
		// Costa Rica
		{"Costa Rica", "America/Costa_Rica", -6},
		{"San Jose", "America/Costa_Rica", -6},
		
		// Côte d'Ivoire
		{"Côte d'Ivoire", "Africa/Abidjan", 0},
		{"Ivory Coast", "Africa/Abidjan", 0},
		{"Yamoussoukro", "Africa/Abidjan", 0},
		{"Abidjan", "Africa/Abidjan", 0},
		
		// Croatia
		{"Croatia", "Europe/Zagreb", 1},
		{"Zagreb", "Europe/Zagreb", 1},
		{"Split", "Europe/Zagreb", 1},
		
		// Cuba
		{"Cuba", "America/Havana", -5},
		{"Havana", "America/Havana", -5},
		
		// Cyprus
		{"Cyprus", "Asia/Nicosia", 2},
		{"Nicosia", "Asia/Nicosia", 2},
		{"Limassol", "Asia/Nicosia", 2},
		
		// Czechia
		{"Czechia", "Europe/Prague", 1},
		{"Czech Republic", "Europe/Prague", 1},
		{"Prague", "Europe/Prague", 1},
		{"Brno", "Europe/Prague", 1},
		
		// Democratic Republic of the Congo
		{"Democratic Republic of the Congo", "Africa/Kinshasa", 1},
		{"DRC", "Africa/Kinshasa", 1},
		{"Kinshasa", "Africa/Kinshasa", 1},
		{"Lubumbashi", "Africa/Lubumbashi", 2},
		
		// Denmark
		{"Denmark", "Europe/Copenhagen", 1},
		{"Copenhagen", "Europe/Copenhagen", 1},
		{"Aarhus", "Europe/Copenhagen", 1},
		
		// Djibouti
		{"Djibouti", "Africa/Djibouti", 3},
		{"Djibouti City", "Africa/Djibouti", 3},
		
		// Dominica
		{"Dominica", "America/Dominica", -4},
		{"Roseau", "America/Dominica", -4},
		
		// Dominican Republic
		{"Dominican Republic", "America/Santo_Domingo", -4},
		{"Santo Domingo", "America/Santo_Domingo", -4},
		{"Santiago de los Caballeros", "America/Santo_Domingo", -4},
		
		// Ecuador
		{"Ecuador", "America/Guayaquil", -5},
		{"Quito", "America/Guayaquil", -5},
		{"Guayaquil", "America/Guayaquil", -5},
		
		// Egypt
		{"Egypt", "Africa/Cairo", 2},
		{"Cairo", "Africa/Cairo", 2},
		{"Alexandria", "Africa/Cairo", 2},
		{"Giza", "Africa/Cairo", 2},
		
		// El Salvador
		{"El Salvador", "America/El_Salvador", -6},
		{"San Salvador", "America/El_Salvador", -6},
		
		// Equatorial Guinea
		{"Equatorial Guinea", "Africa/Malabo", 1},
		{"Malabo", "Africa/Malabo", 1},
		
		// Eritrea
		{"Eritrea", "Africa/Asmara", 3},
		{"Asmara", "Africa/Asmara", 3},
		
		// Estonia
		{"Estonia", "Europe/Tallinn", 2},
		{"Tallinn", "Europe/Tallinn", 2},
		{"Tartu", "Europe/Tallinn", 2},
		
		// Eswatini
		{"Eswatini", "Africa/Mbabane", 2},
		{"Swaziland", "Africa/Mbabane", 2},
		{"Mbabane", "Africa/Mbabane", 2},
		
		// Ethiopia
		{"Ethiopia", "Africa/Addis_Ababa", 3},
		{"Addis Ababa", "Africa/Addis_Ababa", 3},
		
		// Fiji
		{"Fiji", "Pacific/Fiji", 12},
		{"Suva", "Pacific/Fiji", 12},
		
		// Finland
		{"Finland", "Europe/Helsinki", 2},
		{"Helsinki", "Europe/Helsinki", 2},
		{"Tampere", "Europe/Helsinki", 2},
		
		// France
		{"France", "Europe/Paris", 1},
		{"Paris", "Europe/Paris", 1},
		{"Marseille", "Europe/Paris", 1},
		{"Lyon", "Europe/Paris", 1},
		{"Toulouse", "Europe/Paris", 1},
		{"Nice", "Europe/Paris", 1},
		{"Bordeaux", "Europe/Paris", 1},
		
		// Gabon
		{"Gabon", "Africa/Libreville", 1},
		{"Libreville", "Africa/Libreville", 1},
		
		// Gambia
		{"Gambia", "Africa/Banjul", 0},
		{"Banjul", "Africa/Banjul", 0},
		
		// Georgia
		{"Georgia", "Asia/Tbilisi", 4},
		{"Tbilisi", "Asia/Tbilisi", 4},
		
		// Germany
		{"Germany", "Europe/Berlin", 1},
		{"Berlin", "Europe/Berlin", 1},
		{"Hamburg", "Europe/Berlin", 1},
		{"Munich", "Europe/Berlin", 1},
		{"Frankfurt", "Europe/Berlin", 1},
		{"Cologne", "Europe/Berlin", 1},
		{"Stuttgart", "Europe/Berlin", 1},
		{"Dusseldorf", "Europe/Berlin", 1},
		
		// Ghana
		{"Ghana", "Africa/Accra", 0},
		{"Accra", "Africa/Accra", 0},
		{"Kumasi", "Africa/Accra", 0},
		
		// Greece
		{"Greece", "Europe/Athens", 2},
		{"Athens", "Europe/Athens", 2},
		{"Thessaloniki", "Europe/Athens", 2},
		
		// Grenada
		{"Grenada", "America/Grenada", -4},
		{"St. George's", "America/Grenada", -4},
		
		// Guatemala
		{"Guatemala", "America/Guatemala", -6},
		{"Guatemala City", "America/Guatemala", -6},
		
		// Guinea
		{"Guinea", "Africa/Conakry", 0},
		{"Conakry", "Africa/Conakry", 0},
		
		// Guinea-Bissau
		{"Guinea-Bissau", "Africa/Bissau", 0},
		{"Bissau", "Africa/Bissau", 0},
		
		// Guyana
		{"Guyana", "America/Guyana", -4},
		{"Georgetown", "America/Guyana", -4},
		
		// Haiti
		{"Haiti", "America/Port-au-Prince", -5},
		{"Port-au-Prince", "America/Port-au-Prince", -5},
		
		// Holy See
		{"Holy See", "Europe/Vatican", 1},
		{"Vatican City", "Europe/Vatican", 1},
		
		// Honduras
		{"Honduras", "America/Tegucigalpa", -6},
		{"Tegucigalpa", "America/Tegucigalpa", -6},
		
		// Hungary
		{"Hungary", "Europe/Budapest", 1},
		{"Budapest", "Europe/Budapest", 1},
		
		// Iceland
		{"Iceland", "Atlantic/Reykjavik", 0},
		{"Reykjavik", "Atlantic/Reykjavik", 0},
		
		// India
		{"India", "Asia/Kolkata", 5},
		{"New Delhi", "Asia/Kolkata", 5},
		{"Delhi", "Asia/Kolkata", 5},
		{"Mumbai", "Asia/Kolkata", 5},
		{"Bangalore", "Asia/Kolkata", 5},
		{"Hyderabad", "Asia/Kolkata", 5},
		{"Chennai", "Asia/Kolkata", 5},
		{"Kolkata", "Asia/Kolkata", 5},
		{"Pune", "Asia/Kolkata", 5},
		{"Ahmedabad", "Asia/Kolkata", 5},
		
		// Indonesia
		{"Indonesia", "Asia/Jakarta", 7},
		{"Jakarta", "Asia/Jakarta", 7},
		{"Surabaya", "Asia/Jakarta", 7},
		{"Bandung", "Asia/Jakarta", 7},
		{"Medan", "Asia/Jakarta", 7},
		{"Bali", "Asia/Makassar", 8},
		{"Denpasar", "Asia/Makassar", 8},
		
		// Iran
		{"Iran", "Asia/Tehran", 3},
		{"Tehran", "Asia/Tehran", 3},
		{"Mashhad", "Asia/Tehran", 3},
		{"Isfahan", "Asia/Tehran", 3},
		
		// Iraq
		{"Iraq", "Asia/Baghdad", 3},
		{"Baghdad", "Asia/Baghdad", 3},
		{"Basra", "Asia/Baghdad", 3},
		
		// Ireland
		{"Ireland", "Europe/Dublin", 0},
		{"Dublin", "Europe/Dublin", 0},
		{"Cork", "Europe/Dublin", 0},
		
		// Israel
		{"Israel", "Asia/Jerusalem", 2},
		{"Jerusalem", "Asia/Jerusalem", 2},
		{"Tel Aviv", "Asia/Jerusalem", 2},
		{"Haifa", "Asia/Jerusalem", 2},
		
		// Italy
		{"Italy", "Europe/Rome", 1},
		{"Rome", "Europe/Rome", 1},
		{"Milan", "Europe/Rome", 1},
		{"Naples", "Europe/Rome", 1},
		{"Turin", "Europe/Rome", 1},
		{"Florence", "Europe/Rome", 1},
		{"Venice", "Europe/Rome", 1},
		
		// Jamaica
		{"Jamaica", "America/Jamaica", -5},
		{"Kingston", "America/Jamaica", -5},
		
		// Japan
		{"Japan", "Asia/Tokyo", 9},
		{"Tokyo", "Asia/Tokyo", 9},
		{"Osaka", "Asia/Tokyo", 9},
		{"Yokohama", "Asia/Tokyo", 9},
		{"Nagoya", "Asia/Tokyo", 9},
		{"Sapporo", "Asia/Tokyo", 9},
		{"Fukuoka", "Asia/Tokyo", 9},
		{"Kyoto", "Asia/Tokyo", 9},
		
		// Jordan
		{"Jordan", "Asia/Amman", 2},
		{"Amman", "Asia/Amman", 2},
		
		// Kazakhstan
		{"Kazakhstan", "Asia/Almaty", 6},
		{"Nur-Sultan", "Asia/Almaty", 6},
		{"Astana", "Asia/Almaty", 6},
		{"Almaty", "Asia/Almaty", 6},
		
		// Kenya
		{"Kenya", "Africa/Nairobi", 3},
		{"Nairobi", "Africa/Nairobi", 3},
		{"Mombasa", "Africa/Nairobi", 3},
		
		// Kiribati
		{"Kiribati", "Pacific/Tarawa", 12},
		{"Tarawa", "Pacific/Tarawa", 12},
		
		// North Korea
		{"North Korea", "Asia/Pyongyang", 9},
		{"Pyongyang", "Asia/Pyongyang", 9},
		
		// South Korea
		{"South Korea", "Asia/Seoul", 9},
		{"Seoul", "Asia/Seoul", 9},
		{"Busan", "Asia/Seoul", 9},
		{"Incheon", "Asia/Seoul", 9},
		
		// Kuwait
		{"Kuwait", "Asia/Kuwait", 3},
		{"Kuwait City", "Asia/Kuwait", 3},
		
		// Kyrgyzstan
		{"Kyrgyzstan", "Asia/Bishkek", 6},
		{"Bishkek", "Asia/Bishkek", 6},
		
		// Laos
		{"Laos", "Asia/Vientiane", 7},
		{"Vientiane", "Asia/Vientiane", 7},
		
		// Latvia
		{"Latvia", "Europe/Riga", 2},
		{"Riga", "Europe/Riga", 2},
		
		// Lebanon
		{"Lebanon", "Asia/Beirut", 2},
		{"Beirut", "Asia/Beirut", 2},
		
		// Lesotho
		{"Lesotho", "Africa/Maseru", 2},
		{"Maseru", "Africa/Maseru", 2},
		
		// Liberia
		{"Liberia", "Africa/Monrovia", 0},
		{"Monrovia", "Africa/Monrovia", 0},
		
		// Libya
		{"Libya", "Africa/Tripoli", 2},
		{"Tripoli", "Africa/Tripoli", 2},
		{"Benghazi", "Africa/Tripoli", 2},
		
		// Liechtenstein
		{"Liechtenstein", "Europe/Vaduz", 1},
		{"Vaduz", "Europe/Vaduz", 1},
		
		// Lithuania
		{"Lithuania", "Europe/Vilnius", 2},
		{"Vilnius", "Europe/Vilnius", 2},
		
		// Luxembourg
		{"Luxembourg", "Europe/Luxembourg", 1},
		{"Luxembourg City", "Europe/Luxembourg", 1},
		
		// Madagascar
		{"Madagascar", "Indian/Antananarivo", 3},
		{"Antananarivo", "Indian/Antananarivo", 3},
		
		// Malawi
		{"Malawi", "Africa/Blantyre", 2},
		{"Lilongwe", "Africa/Blantyre", 2},
		{"Blantyre", "Africa/Blantyre", 2},
		
		// Malaysia
		{"Malaysia", "Asia/Kuala_Lumpur", 8},
		{"Kuala Lumpur", "Asia/Kuala_Lumpur", 8},
		{"Penang", "Asia/Kuala_Lumpur", 8},
		{"Johor Bahru", "Asia/Kuala_Lumpur", 8},
		
		// Maldives
		{"Maldives", "Indian/Maldives", 5},
		{"Male", "Indian/Maldives", 5},
		
		// Mali
		{"Mali", "Africa/Bamako", 0},
		{"Bamako", "Africa/Bamako", 0},
		
		// Malta
		{"Malta", "Europe/Malta", 1},
		{"Valletta", "Europe/Malta", 1},
		
		// Marshall Islands
		{"Marshall Islands", "Pacific/Majuro", 12},
		{"Majuro", "Pacific/Majuro", 12},
		
		// Mauritania
		{"Mauritania", "Africa/Nouakchott", 0},
		{"Nouakchott", "Africa/Nouakchott", 0},
		
		// Mauritius
		{"Mauritius", "Indian/Mauritius", 4},
		{"Port Louis", "Indian/Mauritius", 4},
		
		// Mexico
		{"Mexico", "America/Mexico_City", -6},
		{"Mexico City", "America/Mexico_City", -6},
		{"Guadalajara", "America/Mexico_City", -6},
		{"Monterrey", "America/Monterrey", -6},
		{"Puebla", "America/Mexico_City", -6},
		{"Tijuana", "America/Tijuana", -8},
		{"Cancun", "America/Cancun", -5},
		
		// Micronesia
		{"Micronesia", "Pacific/Pohnpei", 11},
		{"Palikir", "Pacific/Pohnpei", 11},
		
		// Moldova
		{"Moldova", "Europe/Chisinau", 2},
		{"Chisinau", "Europe/Chisinau", 2},
		
		// Monaco
		{"Monaco", "Europe/Monaco", 1},
		{"Monaco-Ville", "Europe/Monaco", 1},
		
		// Mongolia
		{"Mongolia", "Asia/Ulaanbaatar", 8},
		{"Ulaanbaatar", "Asia/Ulaanbaatar", 8},
		
		// Montenegro
		{"Montenegro", "Europe/Podgorica", 1},
		{"Podgorica", "Europe/Podgorica", 1},
		
		// Morocco
		{"Morocco", "Africa/Casablanca", 0},
		{"Rabat", "Africa/Casablanca", 0},
		{"Casablanca", "Africa/Casablanca", 0},
		{"Marrakesh", "Africa/Casablanca", 0},
		
		// Mozambique
		{"Mozambique", "Africa/Maputo", 2},
		{"Maputo", "Africa/Maputo", 2},
		
		// Myanmar
		{"Myanmar", "Asia/Yangon", 6},
		{"Naypyidaw", "Asia/Yangon", 6},
		{"Yangon", "Asia/Yangon", 6},
		{"Rangoon", "Asia/Yangon", 6},
		
		// Namibia
		{"Namibia", "Africa/Windhoek", 2},
		{"Windhoek", "Africa/Windhoek", 2},
		
		// Nauru
		{"Nauru", "Pacific/Nauru", 12},
		{"Yaren", "Pacific/Nauru", 12},
		
		// Nepal
		{"Nepal", "Asia/Kathmandu", 5},
		{"Kathmandu", "Asia/Kathmandu", 5},
		
		// Netherlands
		{"Netherlands", "Europe/Amsterdam", 1},
		{"Amsterdam", "Europe/Amsterdam", 1},
		{"Rotterdam", "Europe/Amsterdam", 1},
		{"The Hague", "Europe/Amsterdam", 1},
		{"Utrecht", "Europe/Amsterdam", 1},
		
		// New Zealand
		{"New Zealand", "Pacific/Auckland", 12},
		{"Wellington", "Pacific/Auckland", 12},
		{"Auckland", "Pacific/Auckland", 12},
		{"Christchurch", "Pacific/Auckland", 12},
		
		// Nicaragua
		{"Nicaragua", "America/Managua", -6},
		{"Managua", "America/Managua", -6},
		
		// Niger
		{"Niger", "Africa/Niamey", 1},
		{"Niamey", "Africa/Niamey", 1},
		
		// Nigeria
		{"Nigeria", "Africa/Lagos", 1},
		{"Abuja", "Africa/Lagos", 1},
		{"Lagos", "Africa/Lagos", 1},
		{"Kano", "Africa/Lagos", 1},
		
		// North Macedonia
		{"North Macedonia", "Europe/Skopje", 1},
		{"Skopje", "Europe/Skopje", 1},
		
		// Norway
		{"Norway", "Europe/Oslo", 1},
		{"Oslo", "Europe/Oslo", 1},
		{"Bergen", "Europe/Oslo", 1},
		
		// Oman
		{"Oman", "Asia/Muscat", 4},
		{"Muscat", "Asia/Muscat", 4},
		
		// Pakistan
		{"Pakistan", "Asia/Karachi", 5},
		{"Islamabad", "Asia/Karachi", 5},
		{"Karachi", "Asia/Karachi", 5},
		{"Lahore", "Asia/Karachi", 5},
		
		// Palau
		{"Palau", "Pacific/Palau", 9},
		{"Ngerulmud", "Pacific/Palau", 9},
		
		// Panama
		{"Panama", "America/Panama", -5},
		{"Panama City", "America/Panama", -5},
		
		// Papua New Guinea
		{"Papua New Guinea", "Pacific/Port_Moresby", 10},
		{"Port Moresby", "Pacific/Port_Moresby", 10},
		
		// Paraguay
		{"Paraguay", "America/Asuncion", -4},
		{"Asuncion", "America/Asuncion", -4},
		
		// Peru
		{"Peru", "America/Lima", -5},
		{"Lima", "America/Lima", -5},
		{"Cusco", "America/Lima", -5},
		
		// Philippines
		{"Philippines", "Asia/Manila", 8},
		{"Manila", "Asia/Manila", 8},
		{"Quezon City", "Asia/Manila", 8},
		{"Davao", "Asia/Manila", 8},
		{"Cebu City", "Asia/Manila", 8},
		
		// Poland
		{"Poland", "Europe/Warsaw", 1},
		{"Warsaw", "Europe/Warsaw", 1},
		{"Krakow", "Europe/Warsaw", 1},
		{"Gdansk", "Europe/Warsaw", 1},
		
		// Portugal
		{"Portugal", "Europe/Lisbon", 0},
		{"Lisbon", "Europe/Lisbon", 0},
		{"Porto", "Europe/Lisbon", 0},
		
		// Qatar
		{"Qatar", "Asia/Qatar", 3},
		{"Doha", "Asia/Qatar", 3},
		
		// Romania
		{"Romania", "Europe/Bucharest", 2},
		{"Bucharest", "Europe/Bucharest", 2},
		{"Cluj-Napoca", "Europe/Bucharest", 2},
		
		// Russia
		{"Russia", "Europe/Moscow", 3},
		{"Moscow", "Europe/Moscow", 3},
		{"St Petersburg", "Europe/Moscow", 3},
		{"Novosibirsk", "Asia/Novosibirsk", 7},
		{"Yekaterinburg", "Asia/Yekaterinburg", 5},
		{"Vladivostok", "Asia/Vladivostok", 10},
		
		// Rwanda
		{"Rwanda", "Africa/Kigali", 2},
		{"Kigali", "Africa/Kigali", 2},
		
		// Saint Kitts and Nevis
		{"Saint Kitts and Nevis", "America/St_Kitts", -4},
		{"Basseterre", "America/St_Kitts", -4},
		
		// Saint Lucia
		{"Saint Lucia", "America/St_Lucia", -4},
		{"Castries", "America/St_Lucia", -4},
		
		// Saint Vincent and the Grenadines
		{"Saint Vincent and the Grenadines", "America/St_Vincent", -4},
		{"Kingstown", "America/St_Vincent", -4},
		
		// Samoa
		{"Samoa", "Pacific/Apia", 13},
		{"Apia", "Pacific/Apia", 13},
		
		// San Marino
		{"San Marino", "Europe/San_Marino", 1},
		{"City of San Marino", "Europe/San_Marino", 1},
		
		// Sao Tome and Principe
		{"Sao Tome and Principe", "Africa/Sao_Tome", 0},
		{"Sao Tome", "Africa/Sao_Tome", 0},
		
		// Saudi Arabia
		{"Saudi Arabia", "Asia/Riyadh", 3},
		{"Riyadh", "Asia/Riyadh", 3},
		{"Jeddah", "Asia/Riyadh", 3},
		{"Mecca", "Asia/Riyadh", 3},
		{"Medina", "Asia/Riyadh", 3},
		
		// Senegal
		{"Senegal", "Africa/Dakar", 0},
		{"Dakar", "Africa/Dakar", 0},
		
		// Serbia
		{"Serbia", "Europe/Belgrade", 1},
		{"Belgrade", "Europe/Belgrade", 1},
		
		// Seychelles
		{"Seychelles", "Indian/Mahe", 4},
		{"Victoria", "Indian/Mahe", 4},
		
		// Sierra Leone
		{"Sierra Leone", "Africa/Freetown", 0},
		{"Freetown", "Africa/Freetown", 0},
		
		// Singapore
		{"Singapore", "Asia/Singapore", 8},
		
		// Slovakia
		{"Slovakia", "Europe/Bratislava", 1},
		{"Bratislava", "Europe/Bratislava", 1},
		
		// Slovenia
		{"Slovenia", "Europe/Ljubljana", 1},
		{"Ljubljana", "Europe/Ljubljana", 1},
		
		// Solomon Islands
		{"Solomon Islands", "Pacific/Guadalcanal", 11},
		{"Honiara", "Pacific/Guadalcanal", 11},
		
		// Somalia
		{"Somalia", "Africa/Mogadishu", 3},
		{"Mogadishu", "Africa/Mogadishu", 3},
		
		// South Africa
		{"South Africa", "Africa/Johannesburg", 2},
		{"Pretoria", "Africa/Johannesburg", 2},
		{"Johannesburg", "Africa/Johannesburg", 2},
		{"Cape Town", "Africa/Johannesburg", 2},
		{"Durban", "Africa/Johannesburg", 2},
		
		// South Sudan
		{"South Sudan", "Africa/Juba", 2},
		{"Juba", "Africa/Juba", 2},
		
		// Spain
		{"Spain", "Europe/Madrid", 1},
		{"Madrid", "Europe/Madrid", 1},
		{"Barcelona", "Europe/Madrid", 1},
		{"Valencia", "Europe/Madrid", 1},
		{"Seville", "Europe/Madrid", 1},
		
		// Sri Lanka
		{"Sri Lanka", "Asia/Colombo", 5},
		{"Colombo", "Asia/Colombo", 5},
		
		// Sudan
		{"Sudan", "Africa/Khartoum", 2},
		{"Khartoum", "Africa/Khartoum", 2},
		
		// Suriname
		{"Suriname", "America/Paramaribo", -3},
		{"Paramaribo", "America/Paramaribo", -3},
		
		// Sweden
		{"Sweden", "Europe/Stockholm", 1},
		{"Stockholm", "Europe/Stockholm", 1},
		{"Gothenburg", "Europe/Stockholm", 1},
		{"Malmo", "Europe/Stockholm", 1},
		
		// Switzerland
		{"Switzerland", "Europe/Zurich", 1},
		{"Bern", "Europe/Zurich", 1},
		{"Zurich", "Europe/Zurich", 1},
		{"Geneva", "Europe/Zurich", 1},
		
		// Syria
		{"Syria", "Asia/Damascus", 2},
		{"Damascus", "Asia/Damascus", 2},
		{"Aleppo", "Asia/Damascus", 2},
		
		// Taiwan
		{"Taiwan", "Asia/Taipei", 8},
		{"Taipei", "Asia/Taipei", 8},
		{"Kaohsiung", "Asia/Taipei", 8},
		
		// Tajikistan
		{"Tajikistan", "Asia/Dushanbe", 5},
		{"Dushanbe", "Asia/Dushanbe", 5},
		
		// Tanzania
		{"Tanzania", "Africa/Dar_es_Salaam", 3},
		{"Dodoma", "Africa/Dar_es_Salaam", 3},
		{"Dar es Salaam", "Africa/Dar_es_Salaam", 3},
		
		// Thailand
		{"Thailand", "Asia/Bangkok", 7},
		{"Bangkok", "Asia/Bangkok", 7},
		{"Chiang Mai", "Asia/Bangkok", 7},
		{"Phuket", "Asia/Bangkok", 7},
		
		// Timor-Leste
		{"Timor-Leste", "Asia/Dili", 9},
		{"East Timor", "Asia/Dili", 9},
		{"Dili", "Asia/Dili", 9},
		
		// Togo
		{"Togo", "Africa/Lome", 0},
		{"Lome", "Africa/Lome", 0},
		
		// Tonga
		{"Tonga", "Pacific/Tongatapu", 13},
		{"Nuku'alofa", "Pacific/Tongatapu", 13},
		
		// Trinidad and Tobago
		{"Trinidad and Tobago", "America/Port_of_Spain", -4},
		{"Port of Spain", "America/Port_of_Spain", -4},
		
		// Tunisia
		{"Tunisia", "Africa/Tunis", 1},
		{"Tunis", "Africa/Tunis", 1},
		
		// Turkey
		{"Turkey", "Europe/Istanbul", 3},
		{"Ankara", "Europe/Istanbul", 3},
		{"Istanbul", "Europe/Istanbul", 3},
		{"Izmir", "Europe/Istanbul", 3},
		
		// Turkmenistan
		{"Turkmenistan", "Asia/Ashgabat", 5},
		{"Ashgabat", "Asia/Ashgabat", 5},
		
		// Tuvalu
		{"Tuvalu", "Pacific/Funafuti", 12},
		{"Funafuti", "Pacific/Funafuti", 12},
		
		// Uganda
		{"Uganda", "Africa/Kampala", 3},
		{"Kampala", "Africa/Kampala", 3},
		
		// Ukraine
		{"Ukraine", "Europe/Kyiv", 2},
		{"Kyiv", "Europe/Kyiv", 2},
		{"Kiev", "Europe/Kyiv", 2},
		{"Kharkiv", "Europe/Kyiv", 2},
		{"Odessa", "Europe/Kyiv", 2},
		
		// United Arab Emirates
		{"United Arab Emirates", "Asia/Dubai", 4},
		{"UAE", "Asia/Dubai", 4},
		{"Abu Dhabi", "Asia/Dubai", 4},
		{"Dubai", "Asia/Dubai", 4},
		{"Sharjah", "Asia/Dubai", 4},
		
		// United Kingdom
		{"United Kingdom", "Europe/London", 0},
		{"UK", "Europe/London", 0},
		{"London", "Europe/London", 0},
		{"Birmingham", "Europe/London", 0},
		{"Manchester", "Europe/London", 0},
		{"Edinburgh", "Europe/London", 0},
		{"Glasgow", "Europe/London", 0},
		{"Liverpool", "Europe/London", 0},
		{"Leeds", "Europe/London", 0},
		
		// United States
		{"United States", "America/New_York", -5},
		{"USA", "America/New_York", -5},
		{"Washington DC", "America/New_York", -5},
		{"New York", "America/New_York", -5},
		{"Los Angeles", "America/Los_Angeles", -8},
		{"Chicago", "America/Chicago", -6},
		{"Houston", "America/Chicago", -6},
		{"Phoenix", "America/Phoenix", -7},
		{"Philadelphia", "America/New_York", -5},
		{"San Antonio", "America/Chicago", -6},
		{"San Diego", "America/Los_Angeles", -8},
		{"Dallas", "America/Chicago", -6},
		{"San Francisco", "America/Los_Angeles", -8},
		{"Austin", "America/Chicago", -6},
		{"Seattle", "America/Los_Angeles", -8},
		{"Denver", "America/Denver", -7},
		{"Boston", "America/New_York", -5},
		{"Miami", "America/New_York", -5},
		{"Atlanta", "America/New_York", -5},
		{"Las Vegas", "America/Los_Angeles", -8},
		{"Portland", "America/Los_Angeles", -8},
		
		// Uruguay
		{"Uruguay", "America/Montevideo", -3},
		{"Montevideo", "America/Montevideo", -3},
		
		// Uzbekistan
		{"Uzbekistan", "Asia/Tashkent", 5},
		{"Tashkent", "Asia/Tashkent", 5},
		
		// Vanuatu
		{"Vanuatu", "Pacific/Efate", 11},
		{"Port Vila", "Pacific/Efate", 11},
		
		// Venezuela
		{"Venezuela", "America/Caracas", -4},
		{"Caracas", "America/Caracas", -4},
		{"Maracaibo", "America/Caracas", -4},
		
		// Vietnam
		{"Vietnam", "Asia/Ho_Chi_Minh", 7},
		{"Hanoi", "Asia/Bangkok", 7},
		{"Ho Chi Minh City", "Asia/Ho_Chi_Minh", 7},
		{"Saigon", "Asia/Ho_Chi_Minh", 7},
		
		// Yemen
		{"Yemen", "Asia/Aden", 3},
		{"Sana'a", "Asia/Aden", 3},
		{"Aden", "Asia/Aden", 3},
		
		// Zambia
		{"Zambia", "Africa/Lusaka", 2},
		{"Lusaka", "Africa/Lusaka", 2},
		
		// Zimbabwe
		{"Zimbabwe", "Africa/Harare", 2},
		{"Harare", "Africa/Harare", 2},
		
		// ===== Additional Territories & Disputed Areas =====
		
		// State of Palestine
		{"State of Palestine", "Asia/Hebron", 2},
		{"Palestine", "Asia/Hebron", 2},
		{"Ramallah", "Asia/Hebron", 2},
		
		// Western Sahara
		{"Western Sahara", "Africa/El_Aaiun", 0},
		
		// Kosovo
		{"Kosovo", "Europe/Belgrade", 1},
		{"Pristina", "Europe/Belgrade", 1},
		
		// ===== Additional Pacific Islands =====
		
		{"Pago Pago", "Pacific/Pago_Pago", -11},
		{"Tahiti", "Pacific/Tahiti", -10},
		{"Papeete", "Pacific/Tahiti", -10},
		{"Noumea", "Pacific/Noumea", 11},
		{"Guam", "Pacific/Guam", 10},
		{"Saipan", "Pacific/Saipan", 10},
		{"Honolulu", "Pacific/Honolulu", -10},
		{"Hong Kong", "Asia/Hong_Kong", 8},
		{"Macau", "Asia/Macau", 8},
	}
	
	for _, loc := range locations {
		key := strings.ToLower(loc.Name)
		s.locations[key] = &Location{
			Name:     loc.Name,
			IanaName: loc.IanaName,
			Offset:   loc.Offset,
		}
	}
}
// GetLocation retrieves a location by name.
func (s *System) GetLocation(name string) (*Location, error) {
	loc, ok := s.locations[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("unknown timezone: %s", name)
	}
	return loc, nil
}

// GetOffset returns the time difference between two locations in hours.
func (s *System) GetOffset(from, to string) (int, error) {
	fromLoc, err := s.GetLocation(from)
	if err != nil {
		return 0, err
	}
	
	toLoc, err := s.GetLocation(to)
	if err != nil {
		return 0, err
	}
	
	return toLoc.Offset - fromLoc.Offset, nil
}

// ConvertTime converts a time from one timezone to another.
func (s *System) ConvertTime(t time.Time, from, to string) (time.Time, error) {
	offset, err := s.GetOffset(from, to)
	if err != nil {
		return time.Time{}, err
	}
	
	return t.Add(time.Duration(offset) * time.Hour), nil
}

// ListLocations returns all available timezone locations.
func (s *System) ListLocations() []string {
	var names []string
	for _, loc := range s.locations {
		names = append(names, loc.Name)
	}
	return names
}

// ParseTimeString parses a time string like "10:00", "14:30", etc.
func ParseTimeString(s string) (time.Time, error) {
	// Try various time formats
	formats := []string{
		"15:04",
		"3:04pm",
		"3:04 pm",
		"3pm",
		"15:04:05",
	}
	
	now := time.Now()
	
	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			// Combine parsed time with today's date
			return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local), nil
		}
	}
	
	return time.Time{}, fmt.Errorf("unable to parse time: %s", s)
}
