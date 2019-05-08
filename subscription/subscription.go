package subscription

type Currency struct {
	IsSubscribed bool
	SubsChan     chan bool
}

type Nodes struct {
	Eth Currency
	Etc Currency
	Btc Currency
	Ltc Currency
	Bch Currency
	Xlm Currency
}

type Subscription struct {
	Info map[int]Nodes
}

func SubNew() Subscription {
	return Subscription{
		Info: make(map[int]Nodes),
	}
}

func (ds *Subscription) Set(key int, value chan bool, currency string) {
	switch currency {
	case "eth":
		entries := ds.Info[key]
		entries.Eth.IsSubscribed = true
		entries.Eth.SubsChan = value
		ds.Info[key] = entries

	case "etc":
		entries := ds.Info[key]
		entries.Etc.IsSubscribed = true
		entries.Etc.SubsChan = value
		ds.Info[key] = entries

	case "btc":
		entries := ds.Info[key]
		entries.Btc.IsSubscribed = true
		entries.Btc.SubsChan = value
		ds.Info[key] = entries

	case "ltc":
		entries := ds.Info[key]
		entries.Ltc.IsSubscribed = true
		entries.Ltc.SubsChan = value
		ds.Info[key] = entries

	case "bch":
		entries := ds.Info[key]
		entries.Bch.IsSubscribed = true
		entries.Bch.SubsChan = value
		ds.Info[key] = entries

	case "xlm":
		entries := ds.Info[key]
		entries.Xlm.IsSubscribed = true
		entries.Xlm.SubsChan = value
		ds.Info[key] = entries
	}
}

func (ds *Subscription) Remove(key int, currency string) {
	switch currency {

	case "eth":
		entries := ds.Info[key]
		entries.Eth.IsSubscribed = false
		entries.Eth.SubsChan = nil
		ds.Info[key] = entries

	case "etc":
		entries := ds.Info[key]
		entries.Etc.IsSubscribed = false
		entries.Etc.SubsChan = nil
		ds.Info[key] = entries

	case "btc":
		entries := ds.Info[key]
		entries.Btc.IsSubscribed = false
		entries.Btc.SubsChan = nil
		ds.Info[key] = entries

	case "ltc":
		entries := ds.Info[key]
		entries.Ltc.IsSubscribed = false
		entries.Ltc.SubsChan = nil
		ds.Info[key] = entries

	case "bch":
		entries := ds.Info[key]
		entries.Bch.IsSubscribed = false
		entries.Bch.SubsChan = nil
		ds.Info[key] = entries

	case "xlm":
		entries := ds.Info[key]
		entries.Xlm.IsSubscribed = false
		entries.Xlm.SubsChan = nil
		ds.Info[key] = entries
	}
}
