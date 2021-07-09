package mojo

// Types is a mapping of format names (or file extensions) to MIME types
// suitable for the HTTP Content-Type and Accept headers.
//
// This list is taken from https://docs.mojolicious.org/Mojolicious/Types#DESCRIPTION.
// The first item in the array is the canonical one to use in
// Content-Type headers. The other items are aliases that we may find in
// requests.
//
// The most common types are already defined.
//
//   appcache -> text/cache-manifest
//   atom     -> application/atom+xml
//   bin      -> application/octet-stream
//   css      -> text/css
//   gif      -> image/gif
//   gz       -> application/x-gzip
//   htm      -> text/html
//   html     -> text/html;charset=UTF-8
//   ico      -> image/x-icon
//   jpeg     -> image/jpeg
//   jpg      -> image/jpeg
//   js       -> application/javascript
//   json     -> application/json;charset=UTF-8
//   mp3      -> audio/mpeg
//   mp4      -> video/mp4
//   ogg      -> audio/ogg
//   ogv      -> video/ogg
//   pdf      -> application/pdf
//   png      -> image/png
//   rss      -> application/rss+xml
//   svg      -> image/svg+xml
//   ttf      -> font/ttf
//   txt      -> text/plain;charset=UTF-8
//   webm     -> video/webm
//   woff     -> font/woff
//   woff2    -> font/woff2
//   xml      -> application/xml,text/xml
//   zip      -> application/zip
var Types = map[string][]string{
	"appcache": []string{"text/cache-manifest"},
	"atom":     []string{"application/atom+xml"},
	"bin":      []string{"application/octet-stream"},
	"css":      []string{"text/css"},
	"gif":      []string{"image/gif"},
	"gz":       []string{"application/x-gzip"},
	"htm":      []string{"text/html"},
	"html":     []string{"text/html;charset=UTF-8"},
	"ico":      []string{"image/x-icon"},
	"jpeg":     []string{"image/jpeg"},
	"jpg":      []string{"image/jpeg"},
	"js":       []string{"application/javascript"},
	"json":     []string{"application/json;charset=UTF-8"},
	"mp3":      []string{"audio/mpeg"},
	"mp4":      []string{"video/mp4"},
	"ogg":      []string{"audio/ogg"},
	"ogv":      []string{"video/ogg"},
	"pdf":      []string{"application/pdf"},
	"png":      []string{"image/png"},
	"rss":      []string{"application/rss+xml"},
	"svg":      []string{"image/svg+xml"},
	"ttf":      []string{"font/ttf"},
	"txt":      []string{"text/plain;charset=UTF-8"},
	"webm":     []string{"video/webm"},
	"woff":     []string{"font/woff"},
	"woff2":    []string{"font/woff2"},
	"xml":      []string{"application/xml", "text/xml"},
	"zip":      []string{"application/zip"},
}
