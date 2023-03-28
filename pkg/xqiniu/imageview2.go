package xqiniu

import "fmt"

type (
	ImageView2Params struct {
		Mode      int64
		Width     int64
		Height    int64
		Format    string
		Interlace string
		Colors    int64
		Quality   int64
		IgnoreErr int64
		Watermark string
	}

	ImageView2Option func(params *ImageView2Params)
)

// ImageView2 七牛图片基本处理
// 官方文档：https://developer.qiniu.com/dora/1279/basic-processing-images-imageview2
func ImageView2(scheme string, bucket Bucket, key string, options ...ImageView2Option) string {
	host := bucket.Host(scheme)
	url := fmt.Sprintf("%s/%s", host, key)

	params := new(ImageView2Params)
	for _, option := range options {
		option(params)
	}

	// 组装url参数
	if params.Mode != 0 {
		url += fmt.Sprintf("?imageView2/%d", params.Mode)
		if params.Width > 0 {
			url += fmt.Sprintf("/w/%d", params.Width)
		}
		if params.Height > 0 {
			url += fmt.Sprintf("/h/%d", params.Height)
		}
		if params.Format != "" {
			url += fmt.Sprintf("/format/%s", params.Format)
		}
		if params.Interlace != "" {
			url += fmt.Sprintf("/interlace/%s", params.Interlace)
		}
		if params.Colors > 0 {
			url += fmt.Sprintf("/colors/%d", params.Colors)
		}
		if params.Quality > 0 {
			url += fmt.Sprintf("/q/%d", params.Quality)
		}
		if params.IgnoreErr > 0 {
			url += fmt.Sprintf("/ignore-error/%d", params.IgnoreErr)
		}
		if params.Watermark != "" {
			url += fmt.Sprintf("|%s", params.Watermark)
		}
	} else if params.Watermark != "" {
		url += fmt.Sprintf("?%s", params.Watermark)
	}

	return url
}

func ImageView2Mode(value int64) ImageView2Option {
	return func(params *ImageView2Params) {
		params.Mode = value
	}
}

func ImageView2Width(value int64) ImageView2Option {
	return func(params *ImageView2Params) {
		params.Width = value
	}
}

func ImageView2Height(value int64) ImageView2Option {
	return func(params *ImageView2Params) {
		params.Height = value
	}
}

func ImageView2Format(value string) ImageView2Option {
	return func(params *ImageView2Params) {
		params.Format = value
	}
}

func ImageView2Interlace(value string) ImageView2Option {
	return func(params *ImageView2Params) {
		params.Interlace = value
	}
}

func ImageView2Colors(value int64) ImageView2Option {
	return func(params *ImageView2Params) {
		params.Colors = value
	}
}

func ImageView2Quality(value int64) ImageView2Option {
	return func(params *ImageView2Params) {
		params.Quality = value
	}
}

func ImageView2IgnoreErr(value int64) ImageView2Option {
	return func(params *ImageView2Params) {
		params.IgnoreErr = value
	}
}

func ImageView2Watermark(value string) ImageView2Option {
	return func(params *ImageView2Params) {
		params.Watermark = value
	}
}
