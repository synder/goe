package middleware

import (
	"synder.me/goe"
	"strings"
	"compress/flate"
	"io"
	"compress/gzip"
)

var encodeMapping = map[string]io.WriteCloser{

}

func parseEncoding(context *goe.Context) string{
	acceptEncoding := strings.ToLower(context.Request.Get("Accept-Encoding"))

	if acceptEncoding == "" {
		return ""
	}

	return ""
}

func Compress() goe.Handler {

	return func(context *goe.Context){
		acceptEncoding := strings.ToLower(context.Request.Get("Accept-Encoding"))

		var compressor io.WriteCloser
		var err error

		if strings.Index(acceptEncoding, "deflate") >=0 {
			compressor, err = flate.NewWriter(context.Response, flate.BestCompression)

			if err != nil {
				context.Next(err)
				return
			}
		}else if(strings.Index(acceptEncoding, "gzip") >= 0){
			compressor = gzip.NewWriter(context.Response)
		}else{
			context.Next(nil)
			return
		}

		if compressor == nil {
			context.Next(nil)
			return
		}

		defer compressor.Close()


	}

}
