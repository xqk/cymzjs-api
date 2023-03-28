package xqiniu

type ImageStyle struct {
	Url            string `json:"url"`
	ShareUrl       string `json:"share_url"`
	DownloadUrl    string `json:"download_url"`
	OriginalUrl    string `json:"original_url"`
	ImageSearchUrl string `json:"image_search_url"`
}

var BucketImageStyle = map[int64]*ImageStyle{
	BucketImolaStatic.Int64(): {
		Url:            "!wxa.pic",
		ShareUrl:       "!wxa.pic",
		DownloadUrl:    "",
		OriginalUrl:    "!wxa.pic",
		ImageSearchUrl: "!wxa.pic",
	},
	BucketNestPhoto.Int64(): {
		Url:            "~500x.jpg",
		ShareUrl:       "~500x.jpg",
		DownloadUrl:    "~jpg",
		OriginalUrl:    "~1200x1200.jpg",
		ImageSearchUrl: "~736x.jpg",
	},
}
