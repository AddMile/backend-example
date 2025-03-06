package shared

type Platform string

const (
	IOSPlatform     Platform = "ios"
	AndroidPlatform Platform = "android"
	WebPlatform     Platform = "web"
)

func (p Platform) String() string {
	return string(p)
}

type Currency string

const (
	USDCurrency Currency = "USD"
	EURCurrency Currency = "EUR"
	GBPCurrency Currency = "GBP"
)

func (c Currency) String() string {
	return string(c)
}

type Provider string

const (
	SolidProvider   Provider = "solid"
	AppleProvider   Provider = "apple"
	AddMileProvider Provider = "addmile"
)

func (p Provider) String() string {
	return string(p)
}

func (p Provider) IsSolid() bool {
	return p == SolidProvider
}

func (p Provider) IsApple() bool {
	return p == AppleProvider
}

func (p Provider) IsAddMile() bool {
	return p == AddMileProvider
}

type Language string

const (
	EnglishLanguage Language = "en"
	SpanishLanguage Language = "es"
)

func (l Language) String() string {
	return string(l)
}

const (
	MunuteUnit Unit = "minute"
	DayUnit    Unit = "day"
	WeekUnit   Unit = "week"
	MonthUnit  Unit = "month"
	YearUnit   Unit = "year"
)

type Unit string

func (u Unit) String() string {
	return string(u)
}

type Environment string

const (
	LocalEnvironment Environment = "local"
	DevEnvironment   Environment = "dev"
	ProdEnvironment  Environment = "prod"
)

func (e Environment) String() string {
	return string(e)
}

func (e Environment) IsLocal() bool {
	return e == LocalEnvironment
}

func (e Environment) IsDev() bool {
	return e == DevEnvironment
}

func (e Environment) IsProd() bool {
	return e == ProdEnvironment
}
