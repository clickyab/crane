package output

import (
	"context"
	"html/template"
	"io"
	"math/rand"

	"clickyab.com/crane/openrtb/v2.5"

	"github.com/clickyab/services/config"
)

const inappTemplateText string = `
<!doctype html>
<html>
<head>
    <meta charset="UTF-8">
    <meta http-equiv="content-type" content="text/html; charset=UTF-8"/>
    <meta name="viewport" content="initial-scale=1, maximum-scale=1, user-scalable=no">
    <title>Clickyab</title>
    <style>
    {{ if .FullScreen }}
*{padding:0;margin:0;}
html,body{width:100%;height:100%;background:#fff;}
a,img{display:block;float:left;cursor:pointer;text-decoration:none}
.portrait a,.portrait img{width:320px;height:480px;}
.landscape a,.landscape img{width:480px;height:320px;}
a.close,a.largeclose{width:24px;height:24px;line-height:24px;font-size:18px;background:rgba(62,73,90,0.92);color:#FFF;text-align:center;position:absolute;display:inline-block;font-family:verdana,helvetica,arial,sans-serif;font-weight:bold;
}
a.largeclose{width:32px;height:32px;line-height:30px;font-size:24px;}
.portrait a.close{left:296px;top:0;}
.portrait a.largeclose{left:288px;}
.landscape a.close{left:456px;top:0;}
.landscape a.largeclose{left:448px;}
    {{ else }}
*{padding: 0;margin: 0;}
html,body{width: 100%;height: 100%;background: #fff;overflow:hidden;}
a{display: block;float: left;width: 100%;height: 100%;cursor: pointer; text-decoration: none}
a.close,a.largeclose{
    width: 24px;
    height: 24px;
    line-height: 24px;
    font-size: 18px;
    background: rgba(62, 73, 90, 0.92);
    color: #FFF;
    text-align: center;
    position: absolute;
    left: 0;
    bottom: 0;
    display: inline-block;
    font-family: verdana,helvetica,arial,sans-serif;
    font-weight: bold;
	z-index: 9999 !important;
}
a.largeclose{ width: 32px; height: 32px; line-height: 32px; font-size: 24px; }
    {{ end }}
    {{ .ExtraStyle }}
    #showb {
        position: relative;border: 0;line-height: 1;
        height: 100%;
        width: 100%;
        z-index: 9999 !important;
        text-align: center;
    }
    #showb a {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%,-50%);
        -moz-transform: translate(-50%,-50%);
        -webkit-transform: translate(-50%,-50%);
        -o-transform: translate(-50%,-50%);
        -ms-transform: translate(-50%,-50%);
        background:  #000;
        border:1px solid #eee;
        color:#fff;
        font-family: tahoma,Arial,"Times New Roman";
        display: inline-block;
        cursor: pointer;
        width: 100%;
        padding: 10px 0;
        display: none;
        z-index: 1000;
        height: auto !important;
    }
	iframe {
		width: 100vw;
		height: 100vh;
	}
    </style>
</head>
<body class="{{ .BodyClass }}">
    {{ .AdMarkup }}
<a onclick="{{ if .FullScreen }}AndroidHide(){{ else }}AndroidClose(){{ end }};" class="{{ .CloseClass }}">x</a>

<script type="text/javascript">
    {{ if lt .SdkVersion 5 }}
	if (window.addEventListener) {
		window.addEventListener("message", onMessage, false);
	}
	else if (window.attachEvent) {
		window.attachEvent("onmessage", onMessage, false);
	}

	function onMessage(event) {
		var data = event.data;
		if (data.url !== undefined) { 
			onClickyabClicked(data.url,null,null)
		}
	}
	{{ end }}
    function showHitted() {
        document.getElementById("hitted").style.display = 'block';
    }

    var ads = document.getElementById('ads');
    if (ads != null) document.getElementById("ads").onclick = showHitted;

    //function
    var impId  = {{ .ImpID }};
    function AndroidSetImpId(impId) {
        clickyab.setImpId(impId);
    }
    function AndroidHide() {
        clickyab.hide();
    }
    function AndroidClose() {
        clickyab.closeFullAd();
    }
    function AndroidShow() {
        clickyab.show();
    }
    function AndroidRefresh() {
        clickyab.refresh();
    }
    function AndroidRefreshIfActive() {
        clickyab.refreshIfActive();
    }
    function AndroidOpenIntent(url) {
        clickyab.openIntent(url);
    }
    function AndroidOpenIntentWithin(url,packagename) {
        clickyab.openIntentWithin(url,packagename);
    }
    function AndroidOnClick() {
        clickyab.onClicked();
    }
    function AndroidHit(url,impId) {
        clickyab.hit(url,impId);
    }
    function AndroidHasNoAds() {
        clickyab.hasNoAds();
    }
    function AndroidSetHasAds(p) {
        clickyab.setHasAds(p);
    }
    //onAdsClick
    function onClickyabClicked(url,packagename,hitUrl) {
        AndroidOnClick();
        //tell android to submit hit to server
        AndroidHit(hitUrl,impId);
        if (packagename == '' || packagename == null){
            AndroidOpenIntent(url);
        }else{
            AndroidOpenIntentWithin(url,packagename);
        }
        AndroidRefreshIfActive();
    }

    document.addEventListener('DOMContentLoaded', function() {
        if (clickyab === undefined) {
    		return;
		}

        AndroidSetImpId(impId);
        {{ if .NoAd }}
            {{ if eq .SdkVersion 3 }}
            setTimeout(function () {
                AndroidHasNoAds();
            }, 100);// 0.1 sec
            {{ else if ge .SdkVersion 4 }}
            setTimeout(function () {
                AndroidSetHasAds(false);
            }, 100);// 0.1 sec
            {{ else }}
            setTimeout(function () {
                {{ if .FullScreen }}AndroidHide(){{ else }}AndroidClose(){{ end }};
            }, 100);// 0.1 sec
            {{ end }}
        {{ else }}
            {{ if ge .SdkVersion 4 }}
            setTimeout(function () {
                AndroidSetHasAds(true);
            }, 100);// 0.1 sec
            {{ end }}
            setInterval(function(){
                AndroidRefreshIfActive();
            },{{ .RefreshMinute }}*60*1000);// 60 sec
        {{ end }}
    });
</script>

</body>
</html>`

var (
	inappTemplate  = template.Must(template.New("inapp-template").Parse(inappTemplateText))
	refreshAppTime = config.RegisterInt("crane.app.refresh", 2, "app refresh time")
)

type inappContext struct {
	FullScreen    bool
	ExtraStyle    string
	BodyClass     string
	CloseClass    string
	ImpID         int
	SdkVersion    int64
	RefreshMinute int
	NoAd          bool
	AdMarkup      template.HTML
}

// RenderApp will render single ad for app
func RenderApp(ctx context.Context, w io.Writer, res *openrtb.BidResponse, full string, sdk int64, size int32) error {
	closeClass := "largeclose"
	if size == 8 {
		closeClass = "close"
	}
	var noAd bool
	var adMarkup string

	if len(res.Seatbid) == 0 || len(res.Seatbid[0].Bid) == 0 {
		noAd = true
	} else {
		adMarkup = res.Seatbid[0].Bid[0].GetAdm()
	}
	return inappTemplate.Execute(w, inappContext{
		ExtraStyle:    "",
		ImpID:         rand.Int(),
		CloseClass:    closeClass,
		RefreshMinute: refreshAppTime.Int(),
		FullScreen:    full != "",
		BodyClass:     full,
		SdkVersion:    sdk,
		NoAd:          noAd,
		AdMarkup:      template.HTML(adMarkup),
	})
}
