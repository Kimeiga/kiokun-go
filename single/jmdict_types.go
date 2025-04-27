// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    jmdictTypes, err := UnmarshalJmdictTypes(bytes)
//    bytes, err = jmdictTypes.Marshal()

package main

import "bytes"
import "errors"

import "encoding/json"

func UnmarshalJmdictTypes(data []byte) (JmdictTypes, error) {
	var r JmdictTypes
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *JmdictTypes) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type JmdictTypes struct {
	Version       string            `json:"version"`
	Languages     []Lang            `json:"languages"`
	CommonOnly    bool              `json:"commonOnly"`
	DictDate      string            `json:"dictDate"`
	DictRevisions []string          `json:"dictRevisions"`
	Tags          map[string]string `json:"tags"`
	Words         []Word            `json:"words"`
}

type Word struct {
	ID    string  `json:"id"`
	Kanji []Kan   `json:"kanji"`
	Kana  []Kan   `json:"kana"`
	Sense []Sense `json:"sense"`
}

type Kan struct {
	Common         bool     `json:"common"`
	Text           string   `json:"text"`
	Tags           []Tag    `json:"tags"`
	AppliesToKanji []string `json:"appliesToKanji,omitempty"`
}

type Sense struct {
	PartOfSpeech   []PartOfSpeech   `json:"partOfSpeech"`
	AppliesToKanji []string         `json:"appliesToKanji"`
	AppliesToKana  []string         `json:"appliesToKana"`
	Related        [][]Antonym      `json:"related"`
	Antonym        [][]Antonym      `json:"antonym"`
	Field          []Field          `json:"field"`
	Dialect        []Dialect        `json:"dialect"`
	Misc           []Misc           `json:"misc"`
	Info           []string         `json:"info"`
	LanguageSource []LanguageSource `json:"languageSource"`
	Gloss          []Gloss          `json:"gloss"`
}

type Gloss struct {
	Lang   Lang        `json:"lang"`
	Gender interface{} `json:"gender"`
	Type   *Type       `json:"type"`
	Text   string      `json:"text"`
}

type LanguageSource struct {
	Lang  Lang    `json:"lang"`
	Full  bool    `json:"full"`
	Wasei bool    `json:"wasei"`
	Text  *string `json:"text"`
}

type Lang string

const (
	Afr     Lang = "afr"
	Ain     Lang = "ain"
	Alg     Lang = "alg"
	Amh     Lang = "amh"
	Ara     Lang = "ara"
	Arn     Lang = "arn"
	Bnt     Lang = "bnt"
	Bre     Lang = "bre"
	Bul     Lang = "bul"
	Bur     Lang = "bur"
	Chi     Lang = "chi"
	Cze     Lang = "cze"
	Dan     Lang = "dan"
	Dut     Lang = "dut"
	Eng     Lang = "eng"
	Epo     Lang = "epo"
	Est     Lang = "est"
	Fil     Lang = "fil"
	Fin     Lang = "fin"
	Fre     Lang = "fre"
	Geo     Lang = "geo"
	Ger     Lang = "ger"
	Glg     Lang = "glg"
	Grc     Lang = "grc"
	Gre     Lang = "gre"
	Haw     Lang = "haw"
	Heb     Lang = "heb"
	Hin     Lang = "hin"
	Hun     Lang = "hun"
	Ice     Lang = "ice"
	Ind     Lang = "ind"
	Ita     Lang = "ita"
	Khm     Lang = "khm"
	Kor     Lang = "kor"
	Kur     Lang = "kur"
	LangChn Lang = "chn"
	LangPol Lang = "pol"
	Lat     Lang = "lat"
	Lit     Lang = "lit"
	Mal     Lang = "mal"
	Mao     Lang = "mao"
	May     Lang = "may"
	Mnc     Lang = "mnc"
	Mol     Lang = "mol"
	Mon     Lang = "mon"
	Nor     Lang = "nor"
	Per     Lang = "per"
	Por     Lang = "por"
	Rum     Lang = "rum"
	Rus     Lang = "rus"
	SAN     Lang = "san"
	SPA     Lang = "spa"
	Scr     Lang = "scr"
	Slo     Lang = "slo"
	Slv     Lang = "slv"
	Som     Lang = "som"
	Swa     Lang = "swa"
	Swe     Lang = "swe"
	Tah     Lang = "tah"
	Tam     Lang = "tam"
	Tgl     Lang = "tgl"
	Tha     Lang = "tha"
	Tib     Lang = "tib"
	Tur     Lang = "tur"
	Ukr     Lang = "ukr"
	Urd     Lang = "urd"
	Uzb     Lang = "uzb"
	Vie     Lang = "vie"
	Yid     Lang = "yid"
)

type Tag string

const (
	Ateji Tag = "ateji"
	Gikun Tag = "gikun"
	IK    Tag = "iK"
	Ik    Tag = "ik"
	Io    Tag = "io"
	OK    Tag = "oK"
	Ok    Tag = "ok"
	RK    Tag = "rK"
	Rk    Tag = "rk"
	SK    Tag = "sK"
	Sk    Tag = "sk"
)

type Dialect string

const (
	Bra  Dialect = "bra"
	Hob  Dialect = "hob"
	Ksb  Dialect = "ksb"
	Ktb  Dialect = "ktb"
	Kyb  Dialect = "kyb"
	Kyu  Dialect = "kyu"
	Nab  Dialect = "nab"
	Osb  Dialect = "osb"
	Rkb  Dialect = "rkb"
	Thb  Dialect = "thb"
	Tsb  Dialect = "tsb"
	Tsug Dialect = "tsug"
)

type Field string

const (
	Agric    Field = "agric"
	Anat     Field = "anat"
	Archeol  Field = "archeol"
	Archit   Field = "archit"
	Art      Field = "art"
	Astron   Field = "astron"
	Audvid   Field = "audvid"
	Aviat    Field = "aviat"
	Baseb    Field = "baseb"
	Biochem  Field = "biochem"
	Biol     Field = "biol"
	Bot      Field = "bot"
	Boxing   Field = "boxing"
	Buddh    Field = "Buddh"
	Bus      Field = "bus"
	Cards    Field = "cards"
	Chem     Field = "chem"
	Chmyth   Field = "chmyth"
	Christn  Field = "Christn"
	Civeng   Field = "civeng"
	Cloth    Field = "cloth"
	Comp     Field = "comp"
	Cryst    Field = "cryst"
	Dent     Field = "dent"
	Ecol     Field = "ecol"
	Econ     Field = "econ"
	Elec     Field = "elec"
	Electr   Field = "electr"
	Embryo   Field = "embryo"
	Engr     Field = "engr"
	Ent      Field = "ent"
	Figskt   Field = "figskt"
	Film     Field = "film"
	Finc     Field = "finc"
	Fish     Field = "fish"
	Food     Field = "food"
	Gardn    Field = "gardn"
	Genet    Field = "genet"
	Geogr    Field = "geogr"
	Geol     Field = "geol"
	Geom     Field = "geom"
	Go       Field = "go"
	Golf     Field = "golf"
	Gramm    Field = "gramm"
	Grmyth   Field = "grmyth"
	Hanaf    Field = "hanaf"
	Horse    Field = "horse"
	Internet Field = "internet"
	Jpmyth   Field = "jpmyth"
	Kabuki   Field = "kabuki"
	Law      Field = "law"
	Ling     Field = "ling"
	Logic    Field = "logic"
	Ma       Field = "MA"
	Mahj     Field = "mahj"
	Manga    Field = "manga"
	Math     Field = "math"
	Mech     Field = "mech"
	Med      Field = "med"
	Met      Field = "met"
	Mil      Field = "mil"
	Min      Field = "min"
	Mining   Field = "mining"
	Motor    Field = "motor"
	Music    Field = "music"
	Noh      Field = "noh"
	Ornith   Field = "ornith"
	Paleo    Field = "paleo"
	Pathol   Field = "pathol"
	Pharm    Field = "pharm"
	Phil     Field = "phil"
	Photo    Field = "photo"
	Physics  Field = "physics"
	Physiol  Field = "physiol"
	Politics Field = "politics"
	Print    Field = "print"
	Prowres  Field = "prowres"
	Psy      Field = "psy"
	Psyanal  Field = "psyanal"
	Psych    Field = "psych"
	Rail     Field = "rail"
	Rommyth  Field = "rommyth"
	Shinto   Field = "Shinto"
	Shogi    Field = "shogi"
	Ski      Field = "ski"
	Sports   Field = "sports"
	Stat     Field = "stat"
	Stockm   Field = "stockm"
	Sumo     Field = "sumo"
	Surg     Field = "surg"
	Telec    Field = "telec"
	Tradem   Field = "tradem"
	Tv       Field = "tv"
	Vet      Field = "vet"
	Vidg     Field = "vidg"
	Zool     Field = "zool"
)

type Type string

const (
	Explanation Type = "explanation"
	Figurative  Type = "figurative"
	Literal     Type = "literal"
	Trademark   Type = "trademark"
)

type Misc string

const (
	Abbr         Misc = "abbr"
	Arch         Misc = "arch"
	Char         Misc = "char"
	Col          Misc = "col"
	Company      Misc = "company"
	Creat        Misc = "creat"
	Dated        Misc = "dated"
	Dei          Misc = "dei"
	Derog        Misc = "derog"
	Doc          Misc = "doc"
	Euph         Misc = "euph"
	Ev           Misc = "ev"
	Fam          Misc = "fam"
	Fem          Misc = "fem"
	Fict         Misc = "fict"
	Form         Misc = "form"
	Given        Misc = "given"
	Group        Misc = "group"
	Hist         Misc = "hist"
	Hon          Misc = "hon"
	Hum          Misc = "hum"
	ID           Misc = "id"
	Joc          Misc = "joc"
	Leg          Misc = "leg"
	MSl          Misc = "m-sl"
	Male         Misc = "male"
	MiscChn      Misc = "chn"
	MiscPol      Misc = "pol"
	Myth         Misc = "myth"
	NetSl        Misc = "net-sl"
	Obj          Misc = "obj"
	Obs          Misc = "obs"
	OnMim        Misc = "on-mim"
	Organization Misc = "organization"
	Person       Misc = "person"
	Place        Misc = "place"
	Poet         Misc = "poet"
	Product      Misc = "product"
	Proverb      Misc = "proverb"
	Quote        Misc = "quote"
	Rare         Misc = "rare"
	Sens         Misc = "sens"
	Serv         Misc = "serv"
	Ship         Misc = "ship"
	Sl           Misc = "sl"
	Surname      Misc = "surname"
	Uk           Misc = "uk"
	Unclass      Misc = "unclass"
	Vulg         Misc = "vulg"
	Work         Misc = "work"
	Yoji         Misc = "yoji"
)

type PartOfSpeech string

const (
	AdjF     PartOfSpeech = "adj-f"
	AdjI     PartOfSpeech = "adj-i"
	AdjIx    PartOfSpeech = "adj-ix"
	AdjKu    PartOfSpeech = "adj-ku"
	AdjNa    PartOfSpeech = "adj-na"
	AdjNari  PartOfSpeech = "adj-nari"
	AdjNo    PartOfSpeech = "adj-no"
	AdjPn    PartOfSpeech = "adj-pn"
	AdjShiku PartOfSpeech = "adj-shiku"
	AdjT     PartOfSpeech = "adj-t"
	Adv      PartOfSpeech = "adv"
	AdvTo    PartOfSpeech = "adv-to"
	Aux      PartOfSpeech = "aux"
	AuxAdj   PartOfSpeech = "aux-adj"
	AuxV     PartOfSpeech = "aux-v"
	Conj     PartOfSpeech = "conj"
	Cop      PartOfSpeech = "cop"
	Ctr      PartOfSpeech = "ctr"
	Exp      PartOfSpeech = "exp"
	Int      PartOfSpeech = "int"
	N        PartOfSpeech = "n"
	NPref    PartOfSpeech = "n-pref"
	NSuf     PartOfSpeech = "n-suf"
	Num      PartOfSpeech = "num"
	Pn       PartOfSpeech = "pn"
	Pref     PartOfSpeech = "pref"
	Prt      PartOfSpeech = "prt"
	Suf      PartOfSpeech = "suf"
	UNC      PartOfSpeech = "unc"
	V1       PartOfSpeech = "v1"
	V1S      PartOfSpeech = "v1-s"
	V2AS     PartOfSpeech = "v2a-s"
	V2BK     PartOfSpeech = "v2b-k"
	V2DS     PartOfSpeech = "v2d-s"
	V2GK     PartOfSpeech = "v2g-k"
	V2GS     PartOfSpeech = "v2g-s"
	V2HK     PartOfSpeech = "v2h-k"
	V2HS     PartOfSpeech = "v2h-s"
	V2KK     PartOfSpeech = "v2k-k"
	V2KS     PartOfSpeech = "v2k-s"
	V2MS     PartOfSpeech = "v2m-s"
	V2NS     PartOfSpeech = "v2n-s"
	V2RK     PartOfSpeech = "v2r-k"
	V2RS     PartOfSpeech = "v2r-s"
	V2SS     PartOfSpeech = "v2s-s"
	V2TK     PartOfSpeech = "v2t-k"
	V2TS     PartOfSpeech = "v2t-s"
	V2WS     PartOfSpeech = "v2w-s"
	V2YK     PartOfSpeech = "v2y-k"
	V2YS     PartOfSpeech = "v2y-s"
	V2ZS     PartOfSpeech = "v2z-s"
	V4B      PartOfSpeech = "v4b"
	V4G      PartOfSpeech = "v4g"
	V4H      PartOfSpeech = "v4h"
	V4K      PartOfSpeech = "v4k"
	V4M      PartOfSpeech = "v4m"
	V4R      PartOfSpeech = "v4r"
	V4S      PartOfSpeech = "v4s"
	V4T      PartOfSpeech = "v4t"
	V5Aru    PartOfSpeech = "v5aru"
	V5B      PartOfSpeech = "v5b"
	V5G      PartOfSpeech = "v5g"
	V5K      PartOfSpeech = "v5k"
	V5KS     PartOfSpeech = "v5k-s"
	V5M      PartOfSpeech = "v5m"
	V5N      PartOfSpeech = "v5n"
	V5R      PartOfSpeech = "v5r"
	V5RI     PartOfSpeech = "v5r-i"
	V5S      PartOfSpeech = "v5s"
	V5T      PartOfSpeech = "v5t"
	V5U      PartOfSpeech = "v5u"
	V5US     PartOfSpeech = "v5u-s"
	VR       PartOfSpeech = "vr"
	VT       PartOfSpeech = "vt"
	Vi       PartOfSpeech = "vi"
	Vk       PartOfSpeech = "vk"
	Vn       PartOfSpeech = "vn"
	Vs       PartOfSpeech = "vs"
	VsC      PartOfSpeech = "vs-c"
	VsI      PartOfSpeech = "vs-i"
	VsS      PartOfSpeech = "vs-s"
	Vz       PartOfSpeech = "vz"
)

type Antonym struct {
	Integer *int64
	String  *string
}

func (x *Antonym) UnmarshalJSON(data []byte) error {
	object, err := unmarshalUnion(data, &x.Integer, nil, nil, &x.String, false, nil, false, nil, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
	}
	return nil
}

func (x *Antonym) MarshalJSON() ([]byte, error) {
	return marshalUnion(x.Integer, nil, nil, x.String, false, nil, false, nil, false, nil, false, nil, false)
}

func unmarshalUnion(data []byte, pi **int64, pf **float64, pb **bool, ps **string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) (bool, error) {
	if pi != nil {
			*pi = nil
	}
	if pf != nil {
			*pf = nil
	}
	if pb != nil {
			*pb = nil
	}
	if ps != nil {
			*ps = nil
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	tok, err := dec.Token()
	if err != nil {
			return false, err
	}

	switch v := tok.(type) {
	case json.Number:
			if pi != nil {
					i, err := v.Int64()
					if err == nil {
							*pi = &i
							return false, nil
					}
			}
			if pf != nil {
					f, err := v.Float64()
					if err == nil {
							*pf = &f
							return false, nil
					}
					return false, errors.New("Unparsable number")
			}
			return false, errors.New("Union does not contain number")
	case float64:
			return false, errors.New("Decoder should not return float64")
	case bool:
			if pb != nil {
					*pb = &v
					return false, nil
			}
			return false, errors.New("Union does not contain bool")
	case string:
			if haveEnum {
					return false, json.Unmarshal(data, pe)
			}
			if ps != nil {
					*ps = &v
					return false, nil
			}
			return false, errors.New("Union does not contain string")
	case nil:
			if nullable {
					return false, nil
			}
			return false, errors.New("Union does not contain null")
	case json.Delim:
			if v == '{' {
					if haveObject {
							return true, json.Unmarshal(data, pc)
					}
					if haveMap {
							return false, json.Unmarshal(data, pm)
					}
					return false, errors.New("Union does not contain object")
			}
			if v == '[' {
					if haveArray {
							return false, json.Unmarshal(data, pa)
					}
					return false, errors.New("Union does not contain array")
			}
			return false, errors.New("Cannot handle delimiter")
	}
	return false, errors.New("Cannot unmarshal union")

}

func marshalUnion(pi *int64, pf *float64, pb *bool, ps *string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) ([]byte, error) {
	if pi != nil {
			return json.Marshal(*pi)
	}
	if pf != nil {
			return json.Marshal(*pf)
	}
	if pb != nil {
			return json.Marshal(*pb)
	}
	if ps != nil {
			return json.Marshal(*ps)
	}
	if haveArray {
			return json.Marshal(pa)
	}
	if haveObject {
			return json.Marshal(pc)
	}
	if haveMap {
			return json.Marshal(pm)
	}
	if haveEnum {
			return json.Marshal(pe)
	}
	if nullable {
			return json.Marshal(nil)
	}
	return nil, errors.New("Union must not be null")
}
