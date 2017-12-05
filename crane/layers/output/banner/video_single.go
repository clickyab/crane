package banner

import (
	"html/template"
	"io"
	"math/rand"

	"strings"

	"fmt"

	"clickyab.com/crane/crane/entity"
)

const videoADTemplateText = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="robots" content="nofollow">
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type">
    <script type="text/javascript">function showb(){ document.getElementById("showb").style.display = "block"; }</script>
    <style>
html, body, div, span, applet, object, iframe,
h1, h2, h3, h4, h5, h6, p, blockquote, pre,
a, abbr, acronym, address, big, cite, code,
del, dfn, em, font, img, ins, kbd, q, s, samp,
small, strike, strong, sub, sup, tt, var,
dl, dt, dd, ol, ul, li,
fieldset, form, label, legend,
table, caption, tbody, tfoot, thead, tr, th, td {
  margin: 0;
  padding: 0;
  border: 0;
  outline: 0;
  font-weight: inherit;
  font-style: inherit;
  font-size: 100%;
  font-family: inherit;
  vertical-align: baseline;
}
/* remember to define focus styles! */
:focus {
  outline: 0;
}
body {
  line-height: 1;
  color: black;
  background: white;
}
ol, ul {
  list-style: none;
}
/* tables still need 'cellspacing="0"' in the markup */
table {
  border-collapse: separate;
  border-spacing: 0;
}
caption, th, td {
  text-align: left;
  font-weight: normal;
}
blockquote:before, blockquote:after,
q:before, q:after {
  content: "";
}
blockquote, q {
  quotes: "" "";
}
        body{ margin: 0; padding: 0; text-align: center; }
        .o{ position:absolute; top:0; left:0; border:0; height:{{ .Height }}px; width:{{ .Width }}px; z-index: 99; }
        #showb{ position:absolute; top:0; left:0; border:0; line-height:{{ .Height }}px; height:{{ .Height }}px; width:{{ .Width }}px; z-index: 100; background: rgba(0, 0, 0, 0.60); text-align: center; }
        .tiny2{ height: 18px; width: 19px; position: absolute; bottom: 0px; right: 0; z-index: 100; background: url("//static.clickyab.com/img/clickyab-tiny.png") right top no-repeat; border-top-left-radius:4px; -moz-border-radius-topleft:4px  }
        .tiny2:hover{ width: 66px;  }
        .tiny{ height: 18px; width: 19px; position: absolute; top: 0px; right: 0; z-index: 100; background: url("//static.clickyab.com/img/clickyab-tiny.png") right top no-repeat; border-bottom-left-radius:4px; -moz-border-radius-bottomleft:4px  }
        .tiny:hover{ width: 66px;  }
        .tiny3{ position: absolute; top: 0px; right: 0; z-index: 100; }
        .butl {background: #4474CB;color: #FFF;padding: 10px;text-decoration: none;border: 2px solid #FFFFFF;font-family: tahoma;font-size: 13px;}
        img.adhere {max-width:100%;height:auto;}
        video {background: #232323 none repeat scroll 0 0;}
    </style>
</head>
<body>
	<div id="video_advertise_{{ .Rand }}">
	    <video id="e_video_{{ .Rand }}" width="{{ .Width }}" height="{{ .Height }}" autoplay muted loop>
		<source src="{{ .Src }}" type="video/mp4">
	    </video>
	    <div class="call-to-action-holder" style="position: absolute;height: 62px;background: rgba(0, 0, 0,0.6);bottom: 0;width: {{ .Width }}px;">
        <a href="{{ .Link }}" target="_blank" style="background-color: rgb(226, 62, 67);width: 100px;text-decoration: none;display: block;text-align: center;padding: 5px;margin: 16px auto 0;color: #fff;font-size: 12px;font-family: tahoma;">مشاهده آگهی</a>
    </div>
	</div>
	<script>

	var clickyab_video = document.getElementById('e_video_{{ .Rand }}');

   document.getElementById("video_advertise_{{ .Rand }}").addEventListener("mouseover", function(t) {
        t.target.muted = false;
    });
    document.getElementById("video_advertise_{{ .Rand }}").addEventListener("mouseout", function(t) {
        t.target.muted = true;
    });

	    function unwrap(wrapper) {
		// place childNodes in document fragment
		var docFrag = document.createDocumentFragment();
		while (wrapper.firstChild) {
		    var child = wrapper.removeChild(wrapper.firstChild);
		    docFrag.appendChild(child);
		}

		// replace wrapper with document fragment
		wrapper.parentNode.replaceChild(docFrag, wrapper);
	    }
	    var link = "{{ .Link }}";
	    org_html = document.getElementById('video_advertise_{{ .Rand }}').innerHTML;
	    appendHtmlLink = "<a id='a_advertise' target='_blank' href='"+ link +"'>" + org_html + "</a>";
	    var FinalElementHtml = document.getElementById('video_advertise_{{ .Rand }}').innerHTML = appendHtmlLink;
	    document.getElementById('video_advertise_{{ .Rand }}').addEventListener("click", function () {
		var linkElement = document.getElementById('a_advertise');
		if (typeof(linkElement) != 'undefined' && linkElement != null)
		{
		    unwrap(document.getElementById('a_advertise'));
		}

	    });
	</script>
	</div>
	</body></html>`

var videoAdTemplate = template.Must(template.New("video_ad").Parse(videoADTemplateText))

type videoData struct {
	Link   string
	Src    string
	Tiny   bool
	Width  string
	Height string
	Rand   int
}

func renderVideoBanner(w io.Writer, ctx entity.Context, s entity.Seat) error {
	src := s.WinnerAdvertise().Media()
	if ctx.Protocol() == entity.HTTPS {
		src = strings.Replace(src, "http://", "https://", -1)
	}
	sa := videoData{
		Link:   s.ClickURL(),
		Height: fmt.Sprint(s.Height()),
		Width:  fmt.Sprint(s.Width()),
		Src:    src,
		Tiny:   true,
		Rand:   rand.Intn(100),
	}

	return videoAdTemplate.Execute(w, sa)
}
