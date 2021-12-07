package http

type shortURLDataIn struct {
	URL string `json:"url"`
}

type shortURLDataOut struct {
	URL string `json:"url"`
}

type csvDataOut [][]string

type loadBalancerURLDataIn struct {
	URLs []string `json:"urls"`
}

type loadBalancerURLDataOut struct {
	URL string `json:"url"`
}
