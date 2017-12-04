package banner

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"

	"clickyab.com/crane/crane/builder"
	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/assert"
	"github.com/mitchellh/mapstructure"
)

var dynTemplate = map[int]*template.Template{}

const additional = `<style>
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
`

var dynamic1 = map[int]string{
	3: `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="robots" content="nofollow">
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type">
    <style>
    body{
    	margin:0;
    }
        @font-face {
font-family: 'behdad';
src:    url('//static.clickyab.com/font/Behdad-Regular-1-0.woff') format('woff'),
url('//static.clickyab.com/font/Behdad-Regular-1-0.otf') format('opentype') ,
url('//static.clickyab.com/font/Behdad-Regular-1-0.ttf') format('truetype');
font-weight: normal;
font-style: normal;
}
.minimal_template .banner_thumb {
direction: rtl;
display: block;
font-family: "behdad",tahoma;
font-size: 16px;
height: 100%;
width: 100%;
}
.minimal_template .banner_thumb h1 , .banner_thumb span {
font-weight: normal;
margin: 0;
}
.minimal_template .banner_thumb h1  {
color: {{.TitleBannerPickerSelector}};
}
.banner_thumb p{
color: {{.DescriptionBannerPickerSelector}};
}

.minimal_template a {
border: medium none;
}
.minimal_template .banner_thumb p {
font-size: 0.85em;
line-height: 1.3;
margin: 0;
}
.minimal_template .banner_thumb h1 {
font-family: "behdad",tahoma;
font-size: 1em;
font-weight: bold;
line-height: 1.4;
}
.minimal_template.banner_size_1 .banner_thumb_title {
left: 5px;
max-width: 155px;
min-width: 155px;
top: 135px;
}
.minimal_template.banner_size_1 {
width: 298px;
height: 248px;
overflow: hidden;
position: relative;
border: 1px solid #c8bfbf;
}
.minimal_template .section_thumb {
position: absolute;
}
.minimal_template.banner_size_1 .banner_thumb_logo {
right: 10px;
top: 22px;
}
.minimal_template.banner_size_1 .banner_thumb_logo img {
height: auto;
max-width: 110px;
}
.minimal_template.banner_size_1 .banner_thumb_product img {
height: auto;
max-width: 126px;
}
.minimal_template.banner_size_1 .banner_thumb_product {
left: 27px;
top: 5px;
}
.minimal_template.banner_size_1 .banner_thumb_text {
left: 5px;
max-width: 155px;
min-width: 155px;
top: 186px;
text-align:right;
}
.minimal_template.banner_size_1 .banner_thumb_cta {
bottom: 10px;
right: 10px;
}
.minimal_template.banner_size_1 .banner_thumb_price {
bottom: 52px;
min-width: 120px;
right: 10px;
}
.minimal_template .banner_thumb_price .off_price {
color:{{.OffPriceBannerPickerSelector}};
text-decoration: line-through;
}
.minimal_template .banner_thumb_price p {
font-size: 1em;
margin: 0;
text-align: center;
}
.minimal_template .banner_thumb_price .live_price {
color:{{.LivePriceBannerPickerSelector}}
}
.minimal_template .banner_thumb .banner_thumb_cta span {
background-color: {{.CtaBannerPickerSelector}};
border-radius: 6px;
color: {{.TextCtaBannerPickerSelector}};
display: block;
font-size: 0.9em;
min-width: 70px;
padding: 5px 15px 5px 35px;
position: relative;
text-align: center;
text-decoration: none;
}
.minimal_template .banner_thumb .banner_thumb_cta svg {
left: 7px;
max-width: 20px;
position: absolute;
top: 8px;
z-index: 1;
}
.minimal_template .banner_thumb .banner_thumb_cta svg g {
fill: {{.IconCtaBannerPickerSelector}};
}
    </style>
</head>
<body>
<div class="minimal_template banner_size_1">
<a title="{{.BannerTitleTextType}}" href="{{.Link}}" target="_blank" style="width:100%; height:100%;display: block;">
<div class="banner_thumb" style="background:{{.BackgroundBannerPickerSelector}}">
<div class="banner_thumb_logo section_thumb">
<img src="{{.Logo}}"/>
</div>
<div class="banner_thumb_product section_thumb">
<img src="{{.Product}}"/>
</div>
<div class="banner_thumb_title section_thumb">
<h1>{{.BannerTitleTextType}}</h1>
</div>
<div class="banner_thumb_text section_thumb">
<p>{{.BannerDescriptionTextType}} </p>
</div>
<div class="banner_thumb_price section_thumb">
<p class="off_price">{{.PriceTextType}}</p>
<p class="live_price">{{.OffPriceTextType}}</p>
</div>
<div class="banner_thumb_cta section_thumb">
<svg width="16px" height="16px" viewBox="0 0 16 16" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:sketch="http://www.bohemiancoding.com/sketch/ns">
<g id="call-to-action-cta" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" sketch:type="MSPage">
<g id="btn-icon" sketch:type="MSLayerGroup" transform="translate(1.000000, 0.000000)" fill="#000000">
<path d="M7.3,0.4 C3.2,0.4 -0.1,3.8 -0.1,8 C-0.1,12.2 3.2,15.6 7.3,15.6 C11.4,15.6 14.7,12.2 14.7,8 C14.7,3.8 11.4,0.4 7.3,0.4 L7.3,0.4 Z M9.5,12 L6.9,12 L3.4,8.4 L6.9,4.8 L9.5,4.8 L6,8.4 L9.5,12 L9.5,12 Z" id="Shape" sketch:type="MSShapeGroup"></path>
</g>
</g>
</svg>
<span>{{.CtaTitleTextType}}</span>
</div>
</div>
</a>
</div>
</body></html>`,
	4: `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="robots" content="nofollow">
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type">
    <style>
    body{
    	margin:0;
    }
 font-family: 'behdad';
	src:    url('//static.clickyab.com/font/Behdad-Regular-1-0.woff') format('woff'),
	url('//static.clickyab.com/font/Behdad-Regular-1-0.otf') format('opentype') ,
	url('//static.clickyab.com/font/Behdad-Regular-1-0.ttf') format('truetype');
	font-style: normal;
	}
	.minimal_template .banner_thumb {
		font-family: 'behdad';
		direction: rtl;
	}
	.minimal_template .banner_thumb h1 , .banner_thumb span {
		font-weight: normal;
		margin: 0;
	}
	.minimal_template a {
		border: medium none;
	}
	.minimal_template .banner_thumb h1 {
		color: {{.BackgroundBannerPickerSelector}};
	}
	.banner_thumb p{
		color: {{.DescriptionBannerPickerSelector}};
	}
	.minimal_template .banner_thumb h1 {
		font-family: "behdad",tahoma;
		font-size: 1em;
		font-weight: bold;
		line-height: 1.4;
	}
	.minimal_template .banner_thumb p {
		font-size: 0.95em;
		line-height: 1.3;
		margin: 0;
	}
	.minimal_template.banner_size_4 .banner_thumb_title {
		top: 167px;
		left: 5px;
		max-width: 180px;
		min-width: 180px;
	}
	.minimal_template.banner_size_4 {
		width: 334px;
		height: 278px;
		overflow: hidden;
		position: relative;
		border: 1px solid #c8bfbf;
	}
	.minimal_template .section_thumb {
		position: absolute;
	}
	.minimal_template .banner_thumb {
		direction: rtl;
		display: block;
		font-family: "behdad",tahoma;
		font-size: 16px;
		height: 100%;
		width: 100%;
	}
	.minimal_template.banner_size_4 .banner_thumb_logo {
		right: 5px;
		top: 22px;
	}
	.minimal_template.banner_size_4 .banner_thumb_logo img {
  height: auto;
  max-width: 132px;
}
	.minimal_template.banner_size_4 .banner_thumb_product img {
		height: auto;
		max-width: 150px;
	}
	.minimal_template.banner_size_4 .banner_thumb_product {
  left: 18px;
  top: 5px;
}
	.minimal_template.banner_size_4 .banner_thumb_text {
		top: 214px;
		left: 5px;
		max-width: 180px;
		min-width: 180px;
	}
	.minimal_template.banner_size_4 .banner_thumb_cta {
		bottom: 5px;
		right: 5px;
	}
	.minimal_template.banner_size_4 .banner_thumb_price {
		bottom: 55px;
		min-width: 131px;
		right: 5px;
	}
	.minimal_template .banner_thumb_price .off_price {
		color:{{.OffPriceBannerPickerSelector}};
		text-decoration: line-through;
	}
	.minimal_template .banner_thumb_price p {
		font-size: 1.1em;
		margin: 0;
		text-align: center;
	}
	.minimal_template .banner_thumb_price .live_price {
		color:{{.LivePriceBannerPickerSelector}}
	}
	.minimal_template .banner_thumb .banner_thumb_cta span {
		background-color: {{.CtaBannerPickerSelector}};
		border-radius: 6px;
		color: {{.TextCtaBannerPickerSelector}};
		display: block;
		font-size: 1.1em;
		min-width: 87px;
		padding: 5px 15px 5px 35px;
		position: relative;
		text-align: center;
		text-decoration: none;
	}
	.minimal_template .banner_thumb .banner_thumb_cta svg {
		left: 7px;
		max-width: 20px;
		position: absolute;
		top: 10px;
		z-index: 1;
	}
	.minimal_template .banner_thumb .banner_thumb_cta svg g {
	fill: {{.IconCtaBannerPickerSelector}};
	}
    </style>
</head>
<body>
<div class="minimal_template banner_size_4">
	<a title="" class="banner_link" href="{{.Link}}" target="_blank" style="width: 100%; height: 100%; display: block;">
		<div class="banner_thumb" style="background:{{.BackgroundBannerPickerSelector}}">
			<div class="banner_thumb_logo section_thumb">
				<img src="{{.Logo}}"/>
			</div>
			<div class="banner_thumb_product section_thumb">
				<img src="{{.Product}}"/>
			</div>
			<div class="banner_thumb_title section_thumb">
				<h1>{{.BannerTitleTextType}}</h1>
			</div>
			<div class="banner_thumb_text section_thumb">
				<p>{{.BannerDescriptionTextType}}</p>
			</div>
			<div class="banner_thumb_price section_thumb">
				<p class="off_price">{{.PriceTextType}} </p>
				<p class="live_price">{{.OffPriceTextType}} </p>
			</div>
			<div class="banner_thumb_cta section_thumb">
				<svg width="16px" height="16px" viewBox="0 0 16 16" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:sketch="http://www.bohemiancoding.com/sketch/ns">
					<g id="call-to-action-cta" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" sketch:type="MSPage">
						<g id="btn-icon" sketch:type="MSLayerGroup" transform="translate(1.000000, 0.000000)" fill="#000000">
							<path d="M7.3,0.4 C3.2,0.4 -0.1,3.8 -0.1,8 C-0.1,12.2 3.2,15.6 7.3,15.6 C11.4,15.6 14.7,12.2 14.7,8 C14.7,3.8 11.4,0.4 7.3,0.4 L7.3,0.4 Z M9.5,12 L6.9,12 L3.4,8.4 L6.9,4.8 L9.5,4.8 L6,8.4 L9.5,12 L9.5,12 Z" id="Shape" sketch:type="MSShapeGroup"></path>
						</g>
					</g>
				</svg>
				<span>{{.CtaTitleTextType}}</span>
			</div>

		</div>
	</a>
</div>
</body></html>`,
	1: `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="robots" content="nofollow">
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type">
    <style>
    body{
    	margin:0;
    }
  @font-face {
        font-family: 'behdad';
        src:    url('//static.clickyab.com/font/Behdad-Regular-1-0.woff') format('woff'),
	url('//static.clickyab.com/font/Behdad-Regular-1-0.otf') format('opentype') ,
	url('//static.clickyab.com/font/Behdad-Regular-1-0.ttf') format('truetype');
        font-weight: normal;
        font-style: normal;
    }
        .minimal_template .banner_thumb {
            font-family: 'behdad';
            direction: rtl;
        }
        .minimal_template .banner_thumb h1 , .banner_thumb span {
            font-weight: normal;
            margin: 0;
        }
        .minimal_template a {
            border: medium none;
        }
        .minimal_template .banner_thumb h1 {
            color: {{.TitleBannerPickerSelector}};
        }
         .banner_thumb p {
            color: {{.DescriptionBannerPickerSelector}};
         }
        .minimal_template .banner_thumb h1 {
            font-family: "behdad",tahoma;
            font-size: 1em;
            font-weight: bold;
            line-height: 1.4;
            text-align:center
        }
        .minimal_template .banner_thumb p {
            font-size: 0.95em;
            line-height: 1.3;
            margin: 0;
            text-align:center
        }

        .minimal_template.banner_size_5 .banner_thumb_title {
            left: 5px;
            max-width: 110px;
            min-width: 110px;
            top: 293px;
        }
        .minimal_template.banner_size_5 {
            width: 118px;
            height: 598px;
            overflow: hidden;
            position: relative;
            border: 1px solid #c8bfbf;
        }
        .minimal_template .section_thumb {
            position: absolute;
        }
        .minimal_template .banner_thumb {
            direction: rtl;
            display: block;
            font-family: "behdad",tahoma;
            font-size: 16px;
            height: 100%;
            width: 100%;
        }
        .minimal_template.banner_size_5 .banner_thumb_logo {
            right: 5px;
            top: 5px;
        }
        .minimal_template.banner_size_5 .banner_thumb_logo img {
            height: auto;
            max-width: 110px;
        }
        .minimal_template.banner_size_5 .banner_thumb_product img {
            height: auto;
            max-width: 110px;
        }
        .minimal_template.banner_size_5 .banner_thumb_product {
            left: 5px;
            top: 140px;
        }
        .minimal_template.banner_size_5 .banner_thumb_text {
            left: 5px;
            max-width: 110px;
            min-width: 110px;
            top: 365px;
        }
        .minimal_template.banner_size_5 .banner_thumb_cta {
            bottom: 10px;
            right: 5px;
        }
        .minimal_template.banner_size_5 .banner_thumb_price {
            min-width: 110px;
            right: 5px;
            top: 499px;
        }
        .minimal_template .banner_thumb_price .off_price {
            color:{{.OffPriceBannerPickerSelector}};
            text-decoration: line-through;
        }
        .minimal_template .banner_thumb_price p {
            font-size: 0.9em;
            margin: 0;
            text-align: center;
        }
        .minimal_template .banner_thumb_price .live_price {
            color:{{.LivePriceBannerPickerSelector}}
        }
         .minimal_template .banner_thumb .banner_thumb_cta span {
            background-color: {{.CtaBannerPickerSelector}};
            border-radius: 6px;
            color: {{.TextCtaBannerPickerSelector}};
            display: block;
            font-size: 0.8em;
            min-width: 60px;
            padding: 5px 15px 5px 35px;
            position: relative;
            text-align: center;
            text-decoration: none;
        }
        .minimal_template .banner_thumb .banner_thumb_cta svg {
            left: 7px;
            max-width: 20px;
            position: absolute;
            top: 6px;
            z-index: 1;
        }
        .minimal_template .banner_thumb .banner_thumb_cta svg g {
            fill: {{.IconCtaBannerPickerSelector}};
        }
    </style>
</head>
<body>
<div class="minimal_template banner_size_5">
                <a title="" class="banner_link" href="{{.Link}}" target="_blank" style="width: 100%; height: 100%; display: block;">
                    <div class="banner_thumb" style="background:{{.BackgroundBannerPickerSelector}}">
                            <div class="banner_thumb_logo section_thumb">
                                <img src="{{.Logo}}"/>
                            </div>
                            <div class="banner_thumb_product section_thumb">
                                <img src="{{.Product}}"/>
                            </div>
                            <div class="banner_thumb_title section_thumb">
                                <h1>{{.BannerTitleTextType}}</h1>
                            </div>
                            <div class="banner_thumb_text section_thumb">
                                <p>{{.BannerDescriptionTextType}}</p>
                            </div>
                            <div class="banner_thumb_price section_thumb">
                                <p class="off_price">{{.PriceTextType}} </p>
                                <p class="live_price">{{.OffPriceTextType}}</p>
                            </div>
                            <div class="banner_thumb_cta section_thumb">
                                <svg width="16px" height="16px" viewBox="0 0 16 16" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:sketch="http://www.bohemiancoding.com/sketch/ns">
                                    <g id="call-to-action-cta" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" sketch:type="MSPage">
                                        <g id="btn-icon" sketch:type="MSLayerGroup" transform="translate(1.000000, 0.000000)" fill="#000000">
                                            <path d="M7.3,0.4 C3.2,0.4 -0.1,3.8 -0.1,8 C-0.1,12.2 3.2,15.6 7.3,15.6 C11.4,15.6 14.7,12.2 14.7,8 C14.7,3.8 11.4,0.4 7.3,0.4 L7.3,0.4 Z M9.5,12 L6.9,12 L3.4,8.4 L6.9,4.8 L9.5,4.8 L6,8.4 L9.5,12 L9.5,12 Z" id="Shape" sketch:type="MSShapeGroup"></path>
                                        </g>
                                    </g>
                                </svg>
                                <span>{{.CtaTitleTextType}}</span>
                            </div>

                        </div>
                </a>
            </div>
</body></html>`,
	7: `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="robots" content="nofollow">
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type">
    <style>
    body{
    	margin:0;
    }
  @font-face {
	font-family: 'behdad';
	src:    url('//static.clickyab.com/font/Behdad-Regular-1-0.woff') format('woff'),
	url('//static.clickyab.com/font/Behdad-Regular-1-0.otf') format('opentype') ,
	url('//static.clickyab.com/font/Behdad-Regular-1-0.ttf') format('truetype');
	font-weight: normal;
	font-style: normal;
	}
	.minimal_template .banner_thumb {
	font-family: 'behdad';
	direction: rtl;
	}
	.minimal_template .banner_thumb h1 , .banner_thumb span {
	font-weight: normal;
	margin: 0;
	}
	.minimal_template a {
	border: medium none;
	}
	.minimal_template .banner_thumb h1 {
	color: {{.BackgroundBannerPickerSelector}};
	}
	.banner_thumb p{
	color: {{.DescriptionBannerPickerSelector}};
	}
	.minimal_template .banner_thumb h1 {
	font-family: "behdad",tahoma;
	font-size: 1em;
	font-weight: bold;
	line-height: 1.4;
	}
	.minimal_template .banner_thumb p {
	font-size: 0.95em;
	line-height: 1.3;
	margin: 0;
	}
	.minimal_template .banner_thumb h1 {
	font-size: 1em;
	font-weight: bold;
	line-height: 1.4;
	}
	.minimal_template.banner_size_6 {
	width: 118px;
	height: 238px;
	overflow: hidden;
	position: relative;
	border: 1px solid #c8bfbf;
	}
	.minimal_template .section_thumb {
	position: absolute;
	}
	.minimal_template .banner_thumb {
	direction: rtl;
	display: block;
	font-family: "behdad",tahoma;
	font-size: 16px;
	height: 100%;
	width: 100%;
	}
	.minimal_template.banner_size_6 .banner_thumb_logo {
	left: 50%;
	top: 5px;
	transform: translateX(-50%);
	-webkit-transform: translateX(-50%);
	-moz-transform: translateX(-50%);
	-o-transform: translateX(-50%);
	-ms-transform: translateX(-50%);
	}
	.minimal_template.banner_size_6 .banner_thumb_logo img {
	height: auto;
	max-width: 82px;
	}
	.minimal_template.banner_size_6 .banner_thumb_product img {
  height: auto;
  max-width: 107px;
}
	.minimal_template.banner_size_6 .banner_thumb_product {
	left: 50%;
	top: 39px;
	transform: translateX(-50%);
	-webkit-transform: translateX(-50%);
	-moz-transform: translateX(-50%);
	-o-transform: translateX(-50%);
	-ms-transform: translateX(-50%);
	}
	.minimal_template.banner_size_6 .banner_thumb_text {
	left: 5px;
	max-width: 110px;
	min-width: 110px;
	top: 365px;
	}
	.minimal_template.banner_size_6 .banner_thumb_cta {
  bottom: 5px;
  left: 50%;
  transform: translateX(-50%);
  -moz-transform: translateX(-50%);
  -webkit-transform: translateX(-50%);
  -ms-transform: translateX(-50%);
  -o-transform: translateX(-50%);
}
	.minimal_template.banner_size_6 .banner_thumb_price {
          min-width: 110px;
          right: 5px;
          top: 162px;
}
	.minimal_template .banner_thumb_price .off_price {
	color:{{.OffPriceBannerPickerSelector}};
	text-decoration: line-through;
	}
	.minimal_template .banner_thumb_price p {
	font-size: 0.8em;
	margin: 0;
	text-align: center;
	}
	.minimal_template .banner_thumb_price .live_price {
	color:{{.LivePriceBannerPickerSelector}}
	}
	.minimal_template .banner_thumb .banner_thumb_cta span {
	background-color: {{.CtaBannerPickerSelector}};
	border-radius: 6px;
	color: {{.TextCtaBannerPickerSelector}};
	display: block;
	font-size: 0.8em;
	min-width: 60px;
	padding: 3px 10px 3px 27px;
	position: relative;
	text-align: center;
	text-decoration: none;
	}
	.minimal_template .banner_thumb .banner_thumb_cta svg {
	left: 7px;
	max-width: 20px;
	position: absolute;
	top: 5px;
	z-index: 1;
	}
	.minimal_template .banner_thumb .banner_thumb_cta svg g {
	fill: {{.IconCtaBannerPickerSelector}};
	}
    </style>
</head>
<body>
<div class="minimal_template banner_size_6">
	<a title="" class="banner_link" href="{{.Link}}" target="_blank" style="width: 100%; height: 100%; display: block;">
		<div class="banner_thumb" style="background: {{.BackgroundBannerPickerSelector}}">
			<div class="banner_thumb_logo section_thumb">
				<img src="{{.Logo}}"/>
			</div>
			<div class="banner_thumb_product section_thumb">
				<img src="{{.Product}}"/>
			</div>
			<div class="banner_thumb_price section_thumb">
                                <p class="off_price">{{.PriceTextType}} </p>
                                <p class="live_price">{{.OffPriceTextType}}</p>
            </div>
			<div class="banner_thumb_cta section_thumb">
				<svg width="16px" height="16px" viewBox="0 0 16 16" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:sketch="http://www.bohemiancoding.com/sketch/ns">
					<g id="call-to-action-cta" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" sketch:type="MSPage">
						<g id="btn-icon" sketch:type="MSLayerGroup" transform="translate(1.000000, 0.000000)" fill="#000000">
							<path d="M7.3,0.4 C3.2,0.4 -0.1,3.8 -0.1,8 C-0.1,12.2 3.2,15.6 7.3,15.6 C11.4,15.6 14.7,12.2 14.7,8 C14.7,3.8 11.4,0.4 7.3,0.4 L7.3,0.4 Z M9.5,12 L6.9,12 L3.4,8.4 L6.9,4.8 L9.5,4.8 L6,8.4 L9.5,12 L9.5,12 Z" id="Shape" sketch:type="MSShapeGroup"></path>
						</g>
					</g>
				</svg>
				<span>{{.CtaTitleTextType}}</span>
			</div>

		</div>
	</a>
</div>
</body></html>`,
	2: `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="robots" content="nofollow">
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type">
    <style>
    body{
    	margin:0;
    }
  @font-face {
        font-family: 'behdad';
        src:    url('//static.clickyab.com/font/Behdad-Regular-1-0.woff') format('woff'),
	url('//static.clickyab.com/font/Behdad-Regular-1-0.otf') format('opentype') ,
	url('//static.clickyab.com/font/Behdad-Regular-1-0.ttf') format('truetype');
        font-weight: normal;
        font-style: normal;
    }
        .minimal_template .banner_thumb {
            font-family: 'behdad';
            direction: rtl;
        }
        .minimal_template .banner_thumb h1 , .banner_thumb span {
            font-weight: normal;
            margin: 0;
        }
        .minimal_template a {
            border: medium none;
        }
        .minimal_template .banner_thumb h1  {
            color: {{.TitleBannerPickerSelector}};
        }
        .banner_thumb p {
            color: {{.DescriptionBannerPickerSelector}};
        }
        .minimal_template .banner_thumb h1 {
            font-family: "behdad",tahoma;
            font-size: 1em;
            font-weight: bold;
            line-height: 1.4;
            text-align:center
        }
        .minimal_template .banner_thumb p {
            font-size: 0.95em;
            line-height: 1.3;
            margin: 0;
            text-align:center
        }
        .minimal_template.banner_size_7 .banner_thumb_title {
            max-width: 150px;
            min-width: 150px;
            right: 5px;
            top: 296px;
        }
        .minimal_template.banner_size_7 {
            width: 158px;
            height: 598px;
            overflow: hidden;
            position: relative;
            border: 1px solid #c8bfbf;
        }
        .minimal_template .section_thumb {
            position: absolute;
        }
        .minimal_template .banner_thumb {
            direction: rtl;
            display: block;
            font-family: "behdad",tahoma;
            font-size: 16px;
            height: 100%;
            width: 100%;
        }
        .minimal_template.banner_size_7 .banner_thumb_logo {
            left: 50%;
            top: 5px;
            transform: translateX(-50%);
            -webkit-transform: translateX(-50%);
            -moz-transform: translateX(-50%);
            -o-transform: translateX(-50%);
            -ms-transform: translateX(-50%);
        }
        .minimal_template.banner_size_7 .banner_thumb_logo img {
            height: auto;
            max-width: 110px;
        }
        .minimal_template.banner_size_7 .banner_thumb_product img {
            height: auto;
            max-width: 150px;
        }
        .minimal_template.banner_size_7 .banner_thumb_product {
            left: 50%;
            top: 126px;
            transform: translateX(-50%);
            -webkit-transform: translateX(-50%);
            -moz-transform: translateX(-50%);
            -o-transform: translateX(-50%);
            -ms-transform: translateX(-50%);
        }
        .minimal_template.banner_size_7 .banner_thumb_text {
            max-width: 150px;
            min-width: 150px;
            right: 5px;
            top: 365px;
        }
        .minimal_template.banner_size_7 .banner_thumb_cta {
  bottom: 13px;
  left: 50%;
  transform: translateX(-50%);
  -webkit-transform: translateX(-50%);
  -moz-transform: translateX(-50%);
  -ms-transform: translateX(-50%);
  -o-transform: translateX(-50%);
}
        .minimal_template.banner_size_7 .banner_thumb_price {
            min-width: 139px;
            right: 5px;
            top: 500px;
        }
        .minimal_template .banner_thumb_price .off_price {
            color:{{.OffPriceBannerPickerSelector}};
            text-decoration: line-through;
        }
        .minimal_template .banner_thumb_price p {
            font-size: 0.9em;
            margin: 0;
            text-align: center;
        }
        .minimal_template .banner_thumb_price .live_price {
            color:{{.LivePriceBannerPickerSelector}}
        }
        .minimal_template .banner_thumb .banner_thumb_cta span {
            background-color: {{.CtaBannerPickerSelector}};
            border-radius: 6px;
            color: {{.TextCtaBannerPickerSelector}};
            display: block;
            font-size: 0.8em;
            padding: 5px 5px 5px 28px;
            position: relative;
            text-align: center;
            text-decoration: none;
            min-width:100px
        }
        .minimal_template .banner_thumb .banner_thumb_cta svg {
            left: 7px;
            max-width: 20px;
            position: absolute;
            top: 7px;
            z-index: 1;
        }
        .minimal_template .banner_thumb .banner_thumb_cta svg g {
            fill: {{.IconCtaBannerPickerSelector}};
        }
    </style>
</head>
<body>
<div class="minimal_template banner_size_7">
                <a title="" class="banner_link" href="{{.Link}}" target="_blank" style="width: 100%; height: 100%; display: block;">
                    <div class="banner_thumb" style="background:{{.BackgroundBannerPickerSelector}}">
                            <div class="banner_thumb_logo section_thumb">
                                <img src="{{.Logo}}"/>
                            </div>
                            <div class="banner_thumb_product section_thumb">
                                <img src="{{.Product}}"/>
                            </div>
                            <div class="banner_thumb_title section_thumb">
                                <h1>{{.BannerTitleTextType}}</h1>
                            </div>
                            <div class="banner_thumb_text section_thumb">
                                <p>{{.BannerDescriptionTextType}}</p>
                            </div>
                            <div class="banner_thumb_price section_thumb">
                                <p class="off_price">{{.PriceTextType}} </p>
                                <p class="live_price">{{.OffPriceTextType}} </p>
                            </div>
                            <div class="banner_thumb_cta section_thumb">
                                <svg width="16px" height="16px" viewBox="0 0 16 16" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:sketch="http://www.bohemiancoding.com/sketch/ns">
                                    <g id="call-to-action-cta" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" sketch:type="MSPage">
                                        <g id="btn-icon" sketch:type="MSLayerGroup" transform="translate(1.000000, 0.000000)" fill="#000000">
                                            <path d="M7.3,0.4 C3.2,0.4 -0.1,3.8 -0.1,8 C-0.1,12.2 3.2,15.6 7.3,15.6 C11.4,15.6 14.7,12.2 14.7,8 C14.7,3.8 11.4,0.4 7.3,0.4 L7.3,0.4 Z M9.5,12 L6.9,12 L3.4,8.4 L6.9,4.8 L9.5,4.8 L6,8.4 L9.5,12 L9.5,12 Z" id="Shape" sketch:type="MSShapeGroup"></path>
                                        </g>
                                    </g>
                                </svg>
                                <span>{{.CtaTitleTextType}}</span>
                            </div>

                        </div>
                </a>
            </div>
</body></html>`,
	5: `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="robots" content="nofollow">
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type">
    <style>
    body{
    	margin:0;
    }
  @font-face {
            font-family: 'behdad';
            src:    url('//static.clickyab.com/font/Behdad-Regular-1-0.woff') format('woff'),
	url('//static.clickyab.com/font/Behdad-Regular-1-0.otf') format('opentype') ,
	url('//static.clickyab.com/font/Behdad-Regular-1-0.ttf') format('truetype');
            font-weight: normal;
            font-style: normal;
        }
        .minimal_template .banner_thumb {
            font-family: 'behdad';
            direction: rtl;
        }
        .minimal_template .banner_thumb h1 , .banner_thumb span {
            font-weight: normal;
            margin: 0;
        }
        .minimal_template a {
            border: medium none;
        }
        .minimal_template .banner_thumb h1 {
            color: {{.TitleBannerPickerSelector}};
        }
        .banner_thumb p  {
            color: {{.DescriptionBannerPickerSelector}};
        }
        .minimal_template .banner_thumb h1 {
            font-family: "behdad",tahoma;
            font-size: 1em;
            font-weight: bold;
            line-height: 1.4;
        }
        .minimal_template .banner_thumb p {
            font-size: 0.95em;
            line-height: 1.3;
            margin: 0;
        }
        .minimal_template .banner_thumb h1 {
            font-size: 0.8em;
            font-weight: bold;
            line-height: 1.4;
        }
.minimal_template.banner_size_8 .banner_thumb_title {
  max-width: 149px;
  min-width: 149px;
  right: 117px;
  top: 50%;
  transform: translateY(-50%);
  -webkit-transform: translateY(-50%);
  -moz-transform: translateY(-50%);
  -o-transform: translateY(-50%);
  -ms-transform: translateY(-50%);
}
        .minimal_template.banner_size_8 {
            width: 466px;
            height: 58px;
            overflow: hidden;
            position: relative;
            border: 1px solid #c8bfbf;
        }
        .minimal_template .section_thumb {
            position: absolute;
        }
        .minimal_template .banner_thumb {
            direction: rtl;
            display: block;
            font-family: "behdad",tahoma;
            font-size: 16px;
            height: 100%;
            width: 100%;
        }
        .minimal_template.banner_size_8 .banner_thumb_logo {
            right: 5px;
              transform: translateY(-50%);
  -webkit-transform: translateY(-50%);
  -moz-transform: translateY(-50%);
  -o-transform: translateY(-50%);
  -ms-transform: translateY(-50%);
            top: 50%;
        }
        .minimal_template.banner_size_8 .banner_thumb_logo img {
            height: auto;
            max-width: 100px;
        }
        .minimal_template.banner_size_8 .banner_thumb_product img {
  height: auto;
  max-width: 50px;
}
  .minimal_template.banner_size_8 .banner_thumb_product {
  right: 271px;
  top: 2px;
}
        .minimal_template.banner_size_8 .banner_thumb_text {
            max-width: 150px;
            min-width: 150px;
            right: 10px;
            top: 365px;
        }
  .minimal_template.banner_size_8 .banner_thumb_cta {
  bottom: 5px;
  left: 5px;
}
  .minimal_template.banner_size_8 .banner_thumb_price {
  left: 4px;
  min-width: 115px;
  text-align: center;
  top: 2px;
}
        .minimal_template .banner_thumb_price .off_price {
            color:{{.OffPriceBannerPickerSelector}};
            text-decoration: line-through;
            display:none
        }
       .minimal_template .banner_thumb_price p {
  display: inline-block;
  font-size: 0.9em;
  margin: 0;
  text-align: center;
}
        .minimal_template .banner_thumb_price .live_price {
            color:{{.LivePriceBannerPickerSelector}}
        }
        .minimal_template .banner_thumb .banner_thumb_cta span {
            background-color: {{.CtaBannerPickerSelector}};
            border-radius: 6px;
            color: {{.TextCtaBannerPickerSelector}};
            display: block;
            font-size: 0.9em;
            min-width: 70px;
            padding: 2px 15px 2px 35px;
            position: relative;
            text-align: center;
            text-decoration: none;
        }
        .minimal_template .banner_thumb .banner_thumb_cta svg {
  left: 7px;
  max-width: 20px;
  position: absolute;
  top: 5px;
  z-index: 1;
}
        .minimal_template .banner_thumb .banner_thumb_cta svg g {
            fill: {{.IconCtaBannerPickerSelector}};
        }
    </style>
</head>
<body>
<div class="minimal_template banner_size_8">
                <a title="" class="banner_link" href="{{.Link}}" target="_blank" style="width: 100%; height: 100%; display: block;">
                    <div class="banner_thumb" style="background:{{.BackgroundBannerPickerSelector}}">
                            <div class="banner_thumb_logo section_thumb">
                                <img src="{{.Logo}}"/>
                            </div>
                            <div class="banner_thumb_product section_thumb">
                                <img src="{{.Product}}"/>
                            </div>
                            <div class="banner_thumb_title section_thumb">
                                <h1>{{.BannerTitleTextType}}</h1>
                            </div>
                            <div class="banner_thumb_price section_thumb">
                                <p class="off_price">{{.PriceTextType}} </p>
                                <p class="live_price">{{.OffPriceTextType}} </p>
                            </div>

                            <div class="banner_thumb_cta section_thumb">
                                <svg width="16px" height="16px" viewBox="0 0 16 16" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:sketch="http://www.bohemiancoding.com/sketch/ns">
                                    <g id="call-to-action-cta" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" sketch:type="MSPage">
                                        <g id="btn-icon" sketch:type="MSLayerGroup" transform="translate(1.000000, 0.000000)" fill="#000000">
                                            <path d="M7.3,0.4 C3.2,0.4 -0.1,3.8 -0.1,8 C-0.1,12.2 3.2,15.6 7.3,15.6 C11.4,15.6 14.7,12.2 14.7,8 C14.7,3.8 11.4,0.4 7.3,0.4 L7.3,0.4 Z M9.5,12 L6.9,12 L3.4,8.4 L6.9,4.8 L9.5,4.8 L6,8.4 L9.5,12 L9.5,12 Z" id="Shape" sketch:type="MSShapeGroup"></path>
                                        </g>
                                    </g>
                                </svg>
                                <span>{{.CtaTitleTextType}}</span>
                            </div>

                        </div>
                </a>
            </div>
</body></html>`,
	6: `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="robots" content="nofollow">
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type">
    <style>
    body{
    	margin:0;
    }
 @font-face {
            font-family: 'behdad';
            src:    url('//static.clickyab.com/font/Behdad-Regular-1-0.woff') format('woff'),
	url('//static.clickyab.com/font/Behdad-Regular-1-0.otf') format('opentype') ,
	url('//static.clickyab.com/font/Behdad-Regular-1-0.ttf') format('truetype');
            font-weight: normal;
            font-style: normal;
        }
        .minimal_template .banner_thumb {
            font-family: 'behdad';
            direction: rtl;
        }
        .minimal_template .banner_thumb h1 , .banner_thumb span {
            font-weight: normal;
            margin: 0;
        }
        .minimal_template a {
            border: medium none;
        }
        .minimal_template .banner_thumb h1  {
            color: {{.TitleBannerPickerSelector}};
        }
        .banner_thumb p {
            color: {{.DescriptionBannerPickerSelector}};
        }
        .minimal_template .banner_thumb h1 {

        }
        .minimal_template .banner_thumb p {
            font-size: 0.75em;
            line-height: 1.3;
            margin: 0;
        }
        .minimal_template .banner_thumb h1 {
            font-family: "behdad",tahoma;
            font-size: 1em;
            font-weight: bold;
            line-height: 1.4;
        }
        .minimal_template.banner_size_10 .banner_thumb_title {
  max-width: 225px;
  min-width: 225px;
  right: 150px;
  top: 5px;
}
        .minimal_template.banner_size_10 {
            width: 726px;
            height: 88px;
            overflow: hidden;
            position: relative;
            border: 1px solid #c8bfbf;
        }
        .minimal_template .section_thumb {
            position: absolute;
        }
        .minimal_template .banner_thumb {
            direction: rtl;
            display: block;
            font-family: "behdad",tahoma;
            font-size: 16px;
            height: 100%;
            width: 100%;
        }
        .minimal_template.banner_size_10 .banner_thumb_logo {
                  right: 5px;
                  top: 50%;
                  transform: translateY(-50%);
                  -webkit-transform: translateY(-50%);
                  -moz-transform: translateY(-50%);
                  -ms-transform: translateY(-50%);
                  -o-transform: translateY(-50%);
}
        .minimal_template.banner_size_10 .banner_thumb_logo img {
            height: auto;
            max-width: 105px;
        }
        .minimal_template.banner_size_10 .banner_thumb_product img {
            height: auto;
            max-width: 85px;
        }
        .minimal_template.banner_size_10 .banner_thumb_product {
  left: 196px;
  top: 2px;
}
        .minimal_template.banner_size_10 .banner_thumb_text {
  max-width: 235px;
  min-width: 235px;
  right: 150px;
  top: 49px;
}
        .minimal_template.banner_size_10 .banner_thumb_cta {
          bottom: 7px;
          left: 20px;
}
.minimal_template.banner_size_10 .banner_thumb_price {
  left: 19px;
  min-width: 109px;
  top: 5px;
}
        .minimal_template .banner_thumb_price .off_price {
            color:{{.OffPriceBannerPickerSelector}};
            text-decoration: line-through;
        }
        .minimal_template .banner_thumb_price p {
            font-size: 0.9em;
            margin: 0;
            text-align: center;
        }
        .minimal_template .banner_thumb_price .live_price {
            color:{{.LivePriceBannerPickerSelector}}
        }
        .minimal_template .banner_thumb .banner_thumb_cta span {
            background-color: {{.CtaBannerPickerSelector}};
            border-radius: 6px;
            color: {{.TextCtaBannerPickerSelector}};
            display: block;
            font-size: 0.8em;
            min-width: 85px;
            padding: 5px 15px 5px 35px;
            position: relative;
            text-align: center;
            text-decoration: none;
        }
        .minimal_template .banner_thumb .banner_thumb_cta svg {
            left: 7px;
            max-width: 20px;
            position: absolute;
            top: 6px;
            z-index: 1;
        }
        .minimal_template .banner_thumb .banner_thumb_cta svg g {
            fill: {{.IconCtaBannerPickerSelector}};
        }
    </style>
</head>
<body>
<div class="minimal_template banner_size_10">
                <a title="" class="banner_link" href="{{.Link}}" target="_blank" style="width: 100%; height: 100%; display: block;">
                    <div class="banner_thumb" style="background:{{.BackgroundBannerPickerSelector}}">
                            <div class="banner_thumb_logo section_thumb">
                                <img src="{{.Logo}}"/>
                            </div>
                            <div class="banner_thumb_product section_thumb">
                                <img src="{{.Product}}"/>
                            </div>
                            <div class="banner_thumb_title section_thumb">
                                <h1>{{.BannerTitleTextType}}</h1>
                            </div>
                            <div class="banner_thumb_text section_thumb">
                                <p>{{.BannerDescriptionTextType}}</p>
                            </div>
                            <div class="banner_thumb_price section_thumb">
                                <p class="off_price">{{.PriceTextType}} </p>
                                <p class="live_price">{{.OffPriceTextType}} </p>
                            </div>

                            <div class="banner_thumb_cta section_thumb">
                                <svg width="16px" height="16px" viewBox="0 0 16 16" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:sketch="http://www.bohemiancoding.com/sketch/ns">
                                    <g id="call-to-action-cta" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" sketch:type="MSPage">
                                        <g id="btn-icon" sketch:type="MSLayerGroup" transform="translate(1.000000, 0.000000)" fill="#000000">
                                            <path d="M7.3,0.4 C3.2,0.4 -0.1,3.8 -0.1,8 C-0.1,12.2 3.2,15.6 7.3,15.6 C11.4,15.6 14.7,12.2 14.7,8 C14.7,3.8 11.4,0.4 7.3,0.4 L7.3,0.4 Z M9.5,12 L6.9,12 L3.4,8.4 L6.9,4.8 L9.5,4.8 L6,8.4 L9.5,12 L9.5,12 Z" id="Shape" sketch:type="MSShapeGroup"></path>
                                        </g>
                                    </g>
                                </svg>
                                <span>{{.CtaTitleTextType}}</span>
                            </div>
                    </div>
                </a>
            </div>
</body></html>`,
	8: `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="robots" content="nofollow">
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type">
    <style>
    body{
    	margin:0;
    }
 @font-face {
            font-family: 'behdad';
           src:    url('//static.clickyab.com/font/Behdad-Regular-1-0.woff') format('woff'),
	url('//static.clickyab.com/font/Behdad-Regular-1-0.otf') format('opentype') ,
	url('//static.clickyab.com/font/Behdad-Regular-1-0.ttf') format('truetype');
            font-weight: normal;
            font-style: normal;
        }
        .minimal_template .banner_thumb {
            font-family: 'behdad';
            direction: rtl;
        }
        .minimal_template .banner_thumb h1 , .banner_thumb span {
            font-weight: normal;
            margin: 0;
        }
        .minimal_template a {
            border: medium none;
        }
        .minimal_template .banner_thumb h1 {
            color: {{.TitleBannerPickerSelector}};
        }
        .banner_thumb p{
            color: {{.DescriptionBannerPickerSelector}};
        }
        .minimal_template .banner_thumb h1 {
            font-family: "behdad",tahoma;
            font-size: 1em;
            font-weight: bold;
            line-height: 1.4;
        }
        .minimal_template .banner_thumb p {
            font-size: 0.75em;
            line-height: 1.3;
            margin: 0;
        }
        .minimal_template.banner_size_11 .banner_thumb h1 {
            font-size: 0.7em !important;
            font-weight: bold;
            line-height: 1.4;
        }
        .minimal_template.banner_size_11 .banner_thumb_title {
  max-width: 190px;
  min-width: 190px;
  right: 91px;
  top: 3px;
}
.live_price {
        margin-right: 10px !important;
    }
        .minimal_template.banner_size_11 {
            width: 318px;
            height: 48px;
            overflow: hidden;
            position: relative;
            border: 1px solid #c8bfbf;
        }
        .minimal_template .section_thumb {
            position: absolute;
        }
        .minimal_template .banner_thumb {
            direction: rtl;
            display: block;
            font-family: "behdad",tahoma;
            font-size: 16px;
            height: 100%;
            width: 100%;
        }
 .minimal_template.banner_size_11 .banner_thumb_logo {
        right: 5px;
        top: 50%;
        transform: translateY(-50%);
        -moz-transform: translateY(-50%);
        -webkit-transform: translateY(-50%);
        -ms-transform: translateY(-50%);
        -o-transform: translateY(-50%);
    }
        .minimal_template.banner_size_11 .banner_thumb_logo img {
            height: auto;
            max-width: 80px;
        }
        .minimal_template.banner_size_11 .banner_thumb_product img {
            height: auto;
            max-width: 45px;
        }
        .minimal_template.banner_size_11 .banner_thumb_product {
            right: 110px;
            top: 2px;
        }
        .minimal_template.banner_size_11 .banner_thumb_text {
            max-width: 190px;
            min-width: 190px;
            right: 300px;
            top: 42px;
        }
        .minimal_template.banner_size_11 .banner_thumb_cta {
            bottom: 9px;
            left: 6px;
        }
        .minimal_template.banner_size_11 .banner_thumb_price {
  min-width: 150px;
  right: 92px;
  top: 22px;
}
        .minimal_template .banner_thumb_price .off_price {
            color:{{.OffPriceBannerPickerSelector}};
            text-decoration: line-through;
        }
        .minimal_template .banner_thumb_price p {
  display: inline-block;
  font-size: 0.8em;
  margin: 0;
  text-align: center;
}
        .minimal_template .banner_thumb_price .live_price {
            color:{{.LivePriceBannerPickerSelector}};
        }
        .minimal_template.banner_size_11 .banner_thumb .banner_thumb_cta span {
              border-radius: 6px;
              color: transparent !important;
              display: block;
              font-size: 0.8em;
              height: 30px !important;
              min-width: 10px !important;
              padding: 5px 5px 5px 0;
              position: relative;
              text-align: center;
              text-decoration: none;
              width: 30px !important;
              background-color: {{.CtaBannerPickerSelector}};
}
        .minimal_template .banner_thumb .banner_thumb_cta svg {
            left: 7px;
            max-width: 20px;
            position: absolute;
            top: 6px;
            z-index: 1;
        }
        .minimal_template .banner_thumb .banner_thumb_cta svg g {
            fill: {{.IconCtaBannerPickerSelector}};
        }
    </style>
</head>
<body>
<div class="minimal_template banner_size_11">
                <a title="" class="banner_link" href="{{.Link}}" target="_blank" style="width: 100%; height: 100%; display: block;">
                    <div class="banner_thumb" style="background:{{.BackgroundBannerPickerSelector}}">
                            <div class="banner_thumb_logo section_thumb">
                                <img src="{{.Logo}}"/>
                            </div>
                            <div class="banner_thumb_title section_thumb">
                                <h1>{{.BannerTitleTextType}}</h1>
                            </div>
                            <div class="banner_thumb_price section_thumb">
                                <p class="off_price">{{.PriceTextType}} </p>
                                <p class="live_price">{{.OffPriceTextType}} </p>
                            </div>

                            <div class="banner_thumb_cta section_thumb">
                                <svg width="16px" height="16px" viewBox="0 0 16 16" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:sketch="http://www.bohemiancoding.com/sketch/ns">
                                    <g id="call-to-action-cta" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" sketch:type="MSPage">
                                        <g id="btn-icon" sketch:type="MSLayerGroup" transform="translate(1.000000, 0.000000)" fill="#000000">
                                            <path d="M7.3,0.4 C3.2,0.4 -0.1,3.8 -0.1,8 C-0.1,12.2 3.2,15.6 7.3,15.6 C11.4,15.6 14.7,12.2 14.7,8 C14.7,3.8 11.4,0.4 7.3,0.4 L7.3,0.4 Z M9.5,12 L6.9,12 L3.4,8.4 L6.9,4.8 L9.5,4.8 L6,8.4 L9.5,12 L9.5,12 Z" id="Shape" sketch:type="MSShapeGroup"></path>
                                        </g>
                                    </g>
                                </svg>
                                <span>{{.CtaTitleTextType}}</span>
                            </div>
                    </div>
                </a>
            </div>
</body></html>`,
	11: `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="robots" content="nofollow">
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type">
    <style>
    body{
    	margin:0;
    }
 @font-face {
            font-family: 'behdad';
            src:    url('//static.clickyab.com/font/Behdad-Regular-1-0.woff') format('woff'),
	url('//static.clickyab.com/font/Behdad-Regular-1-0.otf') format('opentype') ,
	url('//static.clickyab.com/font/Behdad-Regular-1-0.ttf') format('truetype');
            font-weight: normal;
            font-style: normal;
        }
        .minimal_template .banner_thumb {
            font-family: 'behdad';
            direction: rtl;
        }
        .minimal_template .banner_thumb h1 , .banner_thumb span {
            font-weight: normal;
            margin: 0;
        }
        .minimal_template a {
            border: medium none;
        }
        .minimal_template .banner_thumb h1  {
            color: {{.TitleBannerPickerSelector}};
        }
        .banner_thumb p {
            color: {{.DescriptionBannerPickerSelector}};
        }
        .minimal_template .banner_thumb p {
            font-size: 0.95em;
            line-height: 1.3;
            margin: 0;
        }
        .minimal_template .banner_thumb h1 {
            font-family: "behdad",tahoma;
            font-size: 1em;
            font-weight: bold;
            line-height: 1.4;
        }
        .minimal_template.banner_size_12 .banner_thumb_title {
            max-width: 210px;
            min-width: 210px;
            left: 50%;
            text-align: center;
            top: 369px;
            transform: translateX(-50%);
            -webkit-transform: translateX(-50%);
            -moz-transform: translateX(-50%);
            -o-transform: translateX(-50%);
            -ms-transform: translateX(-50%);
        }
        .minimal_template.banner_size_12 {
            width: 298px;
            height: 598px;
            overflow: hidden;
            position: relative;
            border: 1px solid #c8bfbf;
        }
        .minimal_template .section_thumb {
            position: absolute;
        }
        .minimal_template .banner_thumb {
            direction: rtl;
            display: block;
            font-family: "behdad",tahoma;
            font-size: 16px;
            height: 100%;
            width: 100%;
        }
        .minimal_template.banner_size_12 .banner_thumb_logo {
            left: 50%;
            top: 5px;
            transform: translateX(-50%);
            -webkit-transform: translateX(-50%);
            -moz-transform: translateX(-50%);
            -o-transform: translateX(-50%);
            -ms-transform: translateX(-50%);
        }
        .minimal_template.banner_size_12 .banner_thumb_logo img {
            height: auto;
            max-width: 150px;
        }
        .minimal_template.banner_size_12 .banner_thumb_product img {
            height: auto;
            max-width: 245px;
        }
        .minimal_template.banner_size_12 .banner_thumb_product {
            left: 50%;
            top: 96px;
            transform: translateX(-50%);
            -webkit-transform: translateX(-50%);
            -moz-transform: translateX(-50%);
            -o-transform: translateX(-50%);
            -ms-transform: translateX(-50%);
        }
        .minimal_template.banner_size_12 .banner_thumb_text {
            left: 50%;
            max-width: 190px;
            min-width: 190px;
            text-align: center;
            top: 425px;
            transform: translateX(-50%);
            -webkit-transform: translateX(-50%);
            -moz-transform: translateX(-50%);
            -o-transform: translateX(-50%);
            -ms-transform: translateX(-50%);
        }
        .minimal_template.banner_size_12 .banner_thumb_cta {
            bottom: 11px;
            left: 50%;
            transform: translateX(-50%);
            -webkit-transform: translateX(-50%);
            -moz-transform: translateX(-50%);
            -o-transform: translateX(-50%);
            -ms-transform: translateX(-50%);
        }
        .minimal_template.banner_size_12 .banner_thumb_price {
            min-width: 150px;
            left: 50%;
            top: 500px;
            transform: translateX(-50%);
            -webkit-transform: translateX(-50%);
            -moz-transform: translateX(-50%);
            -o-transform: translateX(-50%);
            -ms-transform: translateX(-50%);
        }
        .minimal_template .banner_thumb_price .off_price {
            color:{{.OffPriceBannerPickerSelector}};
            text-decoration: line-through;
        }
        .minimal_template .banner_thumb_price p {
            font-size: 1.1em;
            margin: 0;
            text-align: center;
        }
        .minimal_template .banner_thumb_price .live_price {
            color:{{.LivePriceBannerPickerSelector}}
        }
        .minimal_template .banner_thumb .banner_thumb_cta span {
            background-color: {{.CtaBannerPickerSelector}};
            border-radius: 6px;
            color: {{.TextCtaBannerPickerSelector}};
            display: block;
            font-size: 0.9em;
            min-width: 82px;
            padding: 5px 15px 5px 35px;
            position: relative;
            text-align: center;
            text-decoration: none;
        }
        .minimal_template .banner_thumb .banner_thumb_cta svg {
            left: 7px;
            max-width: 20px;
            position: absolute;
            top: 7px;
            z-index: 1;
        }
        .minimal_template .banner_thumb .banner_thumb_cta svg g {
            fill: {{.IconCtaBannerPickerSelector}};
        }
    </style>
</head>
<body>
<div class="minimal_template banner_size_12">
                <a title="" class="banner_link" href="{{.Link}}" target="_blank" style="width: 100%; height: 100%; display: block;">
                    <div class="banner_thumb" style="background:{{.BackgroundBannerPickerSelector}}">
                            <div class="banner_thumb_logo section_thumb">
                                <img src="{{.Logo}}"/>
                            </div>
                            <div class="banner_thumb_product section_thumb">
                                <img src="{{.Product}}"/>
                            </div>
                            <div class="banner_thumb_title section_thumb">
                                <h1>{{.BannerTitleTextType}}</h1>
                            </div>
                            <div class="banner_thumb_text section_thumb">
                                <p>{{.BannerDescriptionTextType}}</p>
                            </div>
                            <div class="banner_thumb_price section_thumb">
                                <p class="off_price">{{.PriceTextType}} </p>
                                <p class="live_price">{{.OffPriceTextType}} </p>
                            </div>

                            <div class="banner_thumb_cta section_thumb">
                                <svg width="16px" height="16px" viewBox="0 0 16 16" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:sketch="http://www.bohemiancoding.com/sketch/ns">
                                    <g id="call-to-action-cta" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" sketch:type="MSPage">
                                        <g id="btn-icon" sketch:type="MSLayerGroup" transform="translate(1.000000, 0.000000)" fill="#000000">
                                            <path d="M7.3,0.4 C3.2,0.4 -0.1,3.8 -0.1,8 C-0.1,12.2 3.2,15.6 7.3,15.6 C11.4,15.6 14.7,12.2 14.7,8 C14.7,3.8 11.4,0.4 7.3,0.4 L7.3,0.4 Z M9.5,12 L6.9,12 L3.4,8.4 L6.9,4.8 L9.5,4.8 L6,8.4 L9.5,12 L9.5,12 Z" id="Shape" sketch:type="MSShapeGroup"></path>
                                        </g>
                                    </g>
                                </svg>
                                <span>{{.CtaTitleTextType}}</span>
                            </div>
                    </div>
                </a>
            </div>
</body></html>`,
	14: `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="robots" content="nofollow">
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type">
    <style>
    body{
    	margin:0;
    }
@font-face {
		font-family: 'behdad';
		src:    url('//static.clickyab.com/font/Behdad-Regular-1-0.woff') format('woff'),
	url('//static.clickyab.com/font/Behdad-Regular-1-0.otf') format('opentype') ,
	url('//static.clickyab.com/font/Behdad-Regular-1-0.ttf') format('truetype');
		font-weight: normal;
		font-style: normal;
	}
	.minimal_template .banner_thumb {
		font-family: 'behdad';
		direction: rtl;
	}
	.minimal_template .banner_thumb h1 , .banner_thumb span {
		font-weight: normal;
		margin: 0;
	}
	.minimal_template a {
		border: medium none;
	}
	.minimal_template .banner_thumb h1 {
		color: {{.TitleBannerPickerSelector}};
	}
	.banner_thumb p{
		color: {{.DescriptionBannerPickerSelector}};
	}
	.minimal_template .banner_thumb p {
		font-size: 0.85em;
		line-height: 1.3;
		margin: 0;
	}
	.minimal_template .banner_thumb h1 {
		font-family: "behdad",tahoma;
		font-size: 0.9em;
		font-weight: bold;
		line-height: 1.4;
	}
	.minimal_template.banner_size_13 .banner_thumb_title {
    max-width: 120px;
    min-width: 120px;
    right: 125px;
    top: 135px;
	}
	.minimal_template.banner_size_13 {
		width: 248px;
		height: 248px;
		overflow: hidden;
		position: relative;
		border: 1px solid #c8bfbf;
	}
	.minimal_template .section_thumb {
		position: absolute;
	}
	.minimal_template .banner_thumb {
		direction: rtl;
		display: block;
		font-family: "behdad",tahoma;
		font-size: 16px;
		height: 100%;
		width: 100%;
	}
	.minimal_template.banner_size_13 .banner_thumb_logo {
		right: 10px;
		top: 22px;

	}
	.minimal_template.banner_size_13 .banner_thumb_logo img {
		height: auto;
		max-width: 100px;
	}
	.minimal_template.banner_size_13 .banner_thumb_product img {
		height: auto;
		max-width: 120px;
	}
	.minimal_template.banner_size_13 .banner_thumb_product {
		left: 5px;
		top: 5px;

	}
	.minimal_template.banner_size_13 .banner_thumb_text {
 max-width: 120px;
    min-width: 120px;
    right: 125px;
    top: 186px;
	}
	.minimal_template.banner_size_13 .banner_thumb_cta {
        bottom: 10px;
          right: 10px;
	}
.minimal_template.banner_size_13 .banner_thumb_price {
  bottom: 52px;
  min-width: 110px;
  right: 9px;
  text-align: center;
}
	.minimal_template .banner_thumb_price .off_price {
		color:{{.OffPriceBannerPickerSelector}};
		text-decoration: line-through;
	}
	.minimal_template .banner_thumb_price p {
		font-size: 0.85em;
    line-height: 1.3;
    margin: 0;
	}
	.minimal_template .banner_thumb_price .live_price {
		color:{{.LivePriceBannerPickerSelector}}
	}
	.minimal_template .banner_thumb .banner_thumb_cta span {
		background-color: {{.CtaBannerPickerSelector}};
		border-radius: 6px;
		color: {{.TextCtaBannerPickerSelector}};
		display: block;
		font-size: 0.8em;
		min-width: 60px;
		padding: 5px 15px 5px 35px;
		position: relative;
		text-align: center;
		text-decoration: none;
	}
	.minimal_template .banner_thumb .banner_thumb_cta svg {
		left: 7px;
		max-width: 20px;
		position: absolute;
		top: 6px;
		z-index: 1;
	}
	.minimal_template .banner_thumb .banner_thumb_cta svg g {
		fill: {{.IconCtaBannerPickerSelector}};
	}

    </style>
</head>
<body>
<div class="minimal_template banner_size_13">
	<a title="" class="banner_link" href="{{.Link}}" target="_blank" style="width: 100%; height: 100%; display: block;">
		<div class="banner_thumb" style="background:{{.BackgroundBannerPickerSelector}}">
			<div class="banner_thumb_logo section_thumb">
                                <img src="{{.Logo}}"/>
                            </div>
                            <div class="banner_thumb_product section_thumb">
                                <img src="{{.Product}}"/>
                            </div>
                            <div class="banner_thumb_title section_thumb">
                                <h1>{{.BannerTitleTextType}}</h1>
                            </div>
                            <div class="banner_thumb_text section_thumb">
                                <p>{{.BannerDescriptionTextType}}</p>
                            </div>
                            <div class="banner_thumb_price section_thumb">
                                <p class="off_price">{{.PriceTextType}} </p>
                                <p class="live_price">{{.OffPriceTextType}} </p>
                            </div>

			<div class="banner_thumb_cta section_thumb">
				<svg width="16px" height="16px" viewBox="0 0 16 16" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:sketch="http://www.bohemiancoding.com/sketch/ns">
					<g id="call-to-action-cta" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" sketch:type="MSPage">
						<g id="btn-icon" sketch:type="MSLayerGroup" transform="translate(1.000000, 0.000000)" fill="#000000">
							<path d="M7.3,0.4 C3.2,0.4 -0.1,3.8 -0.1,8 C-0.1,12.2 3.2,15.6 7.3,15.6 C11.4,15.6 14.7,12.2 14.7,8 C14.7,3.8 11.4,0.4 7.3,0.4 L7.3,0.4 Z M9.5,12 L6.9,12 L3.4,8.4 L6.9,4.8 L9.5,4.8 L6,8.4 L9.5,12 L9.5,12 Z" id="Shape" sketch:type="MSShapeGroup"></path>
						</g>
					</g>
				</svg>
				<span>{{.CtaTitleTextType}}</span>
			</div>
		</div>
	</a>
</div>
</body></html>`,
	16: `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="robots" content="nofollow">
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type">
    <style>
    body{
    	margin:0;
    }
  @font-face {
            font-family: 'behdad';
           src:    url('//static.clickyab.com/font/Behdad-Regular-1-0.woff') format('woff'),
	url('//static.clickyab.com/font/Behdad-Regular-1-0.otf') format('opentype') ,
	url('//static.clickyab.com/font/Behdad-Regular-1-0.ttf') format('truetype');
            font-weight: normal;
            font-style: normal;
        }
        .minimal_template .banner_thumb {
            font-family: 'behdad';
            direction: rtl;
        }
        .minimal_template .banner_thumb h1 , .banner_thumb span {
            font-weight: normal;
            margin: 0;
        }
        .minimal_template a {
            border: medium none;
        }
        .minimal_template .banner_thumb h1  {
            color: {{.TitleBannerPickerSelector}};
        }
        .banner_thumb p {
             color: {{.DescriptionBannerPickerSelector}};
        }
        .minimal_template .banner_thumb p {
            font-size: 0.95em;
            line-height: 1.3;
            margin: 0;
        }
        .minimal_template .banner_thumb h1 {
            font-family: "behdad",tahoma;
            font-size: 1em;
            font-weight: bold;
            line-height: 1.4;
        }
        .minimal_template.banner_size_14 .banner_thumb_title {
            max-width: 210px;
            min-width: 210px;
            left: 50%;
            text-align: center;
            top: 265px;
            transform: translateX(-50%);
            -webkit-transform: translateX(-50%);
            -moz-transform: translateX(-50%);
            -o-transform: translateX(-50%);
            -ms-transform: translateX(-50%);
        }
        .minimal_template.banner_size_14 {
            width: 318px;
            height: 478px;
            overflow: hidden;
            position: relative;
            border: 1px solid #c8bfbf;
        }
        .minimal_template .section_thumb {
            position: absolute;
        }
        .minimal_template .banner_thumb {
            direction: rtl;
            display: block;
            font-family: "behdad",tahoma;
            font-size: 16px;
            height: 100%;
            width: 100%;
        }
        .minimal_template.banner_size_14 .banner_thumb_logo {
            left: 50%;
            top: 2px;
            transform: translateX(-50%);
            -webkit-transform: translateX(-50%);
            -moz-transform: translateX(-50%);
            -o-transform: translateX(-50%);
            -ms-transform: translateX(-50%);
        }
        .minimal_template.banner_size_14 .banner_thumb_logo img {
            height: auto;
            max-width: 100px;
        }
        .minimal_template.banner_size_14 .banner_thumb_product img {
            height: auto;
            max-width: 173px;
        }
        .minimal_template.banner_size_14 .banner_thumb_product {
            left: 50%;
            top: 74px;
            transform: translateX(-50%);
            -webkit-transform: translateX(-50%);
            -moz-transform: translateX(-50%);
            -o-transform: translateX(-50%);
            -ms-transform: translateX(-50%);
        }
        .minimal_template.banner_size_14 .banner_thumb_text {
            left: 50%;
            max-width: 190px;
            min-width: 190px;
            text-align: center;
            top: 319px;
            transform: translateX(-50%);
            -webkit-transform: translateX(-50%);
            -moz-transform: translateX(-50%);
            -o-transform: translateX(-50%);
            -ms-transform: translateX(-50%);
        }
        .minimal_template.banner_size_14 .banner_thumb_cta {
            bottom: 11px;
            left: 50%;
            transform: translateX(-50%);
            -webkit-transform: translateX(-50%);
            -moz-transform: translateX(-50%);
            -o-transform: translateX(-50%);
            -ms-transform: translateX(-50%);
        }
        .minimal_template.banner_size_14 .banner_thumb_price {
            min-width: 150px;
            left: 50%;
            top: 384px;
            transform: translateX(-50%);
            -webkit-transform: translateX(-50%);
            -moz-transform: translateX(-50%);
            -o-transform: translateX(-50%);
            -ms-transform: translateX(-50%);
        }
        .minimal_template .banner_thumb_price .off_price {
            color:{{.OffPriceBannerPickerSelector}};
            text-decoration: line-through;
        }
        .minimal_template .banner_thumb_price p {
            font-size: 1.1em;
            margin: 0;
            text-align: center;
        }
        .minimal_template .banner_thumb_price .live_price {
            color:{{.LivePriceBannerPickerSelector}}
        }
        .minimal_template .banner_thumb .banner_thumb_cta span {
            background-color: {{.CtaBannerPickerSelector}};
            border-radius: 6px;
            color: {{.TextCtaBannerPickerSelector}};
            display: block;
            font-size: 1em;
            min-width: 110px;
            padding: 5px 15px 5px 35px;
            position: relative;
            text-align: center;
            text-decoration: none;
        }
        .minimal_template .banner_thumb .banner_thumb_cta svg {
            left: 7px;
            max-width: 20px;
            position: absolute;
            top: 8px;
            z-index: 1;
        }
        .minimal_template .banner_thumb .banner_thumb_cta svg g {
            fill: {{.IconCtaBannerPickerSelector}};
        }

    </style>
</head>
<body>
<div class="minimal_template banner_size_14">
                <a title="" class="banner_link" href="{{.Link}}" target="_blank" style="width: 100%; height: 100%; display: block;">
                    <div class="banner_thumb" style="background:{{.BackgroundBannerPickerSelector}}">
                            <div class="banner_thumb_logo section_thumb">
                                <img src="{{.Logo}}"/>
                            </div>
                            <div class="banner_thumb_product section_thumb">
                                <img src="{{.Product}}"/>
                            </div>
                            <div class="banner_thumb_title section_thumb">
                                <h1>{{.BannerTitleTextType}}</h1>
                            </div>
                            <div class="banner_thumb_text section_thumb">
                                <p>{{.BannerDescriptionTextType}}</p>
                            </div>
                            <div class="banner_thumb_price section_thumb">
                                <p class="off_price">{{.PriceTextType}} </p>
                                <p class="live_price">{{.OffPriceTextType}} </p>
                            </div>

                            <div class="banner_thumb_cta section_thumb">
                                <svg width="16px" height="16px" viewBox="0 0 16 16" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:sketch="http://www.bohemiancoding.com/sketch/ns">
                                    <g id="call-to-action-cta" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" sketch:type="MSPage">
                                        <g id="btn-icon" sketch:type="MSLayerGroup" transform="translate(1.000000, 0.000000)" fill="#000000">
                                            <path d="M7.3,0.4 C3.2,0.4 -0.1,3.8 -0.1,8 C-0.1,12.2 3.2,15.6 7.3,15.6 C11.4,15.6 14.7,12.2 14.7,8 C14.7,3.8 11.4,0.4 7.3,0.4 L7.3,0.4 Z M9.5,12 L6.9,12 L3.4,8.4 L6.9,4.8 L9.5,4.8 L6,8.4 L9.5,12 L9.5,12 Z" id="Shape" sketch:type="MSShapeGroup"></path>
                                        </g>
                                    </g>
                                </svg>
                                <span>{{.CtaTitleTextType}}</span>
                            </div>
                    </div>
                </a>
            </div>
            </body></html>`,
}
var r = &sync.RWMutex{}

// dynamicAttribute ad struct
type dynamicAttribute struct {
	Link                            string `mapstructure:"-"`
	BannerTitleTextType             string `mapstructure:"banner_title_text_type"`
	TemplateID                      string `mapstructure:"template_id"`
	CtaTitleTextType                string `mapstructure:"cta_title_text_type"`
	Logo                            string `mapstructure:"logo"`
	Product                         string `mapstructure:"product"`
	BannerDescriptionTextType       string `mapstructure:"banner_description_text_type"`
	PriceTextType                   string `mapstructure:"price_text_type"`
	OffPriceTextType                string `mapstructure:"off_price_text_type"`
	BackgroundBannerPickerSelector  string `mapstructure:"background_banner_picker_selector"`
	CtaBannerPickerSelector         string `mapstructure:"cta_banner_picker_selector"`
	TitleBannerPickerSelector       string `mapstructure:"title_banner_picker_selector"`
	DescriptionBannerPickerSelector string `mapstructure:"description_banner_picker_selector"`
	IconCtaBannerPickerSelector     string `mapstructure:"icon_cta_banner_picker_selector"`
	TextCtaBannerPickerSelector     string `mapstructure:"text_cta_banner_picker_selector"`
	OffPriceBannerPickerSelector    string `mapstructure:"off_price_banner_picker_selector"`
	LivePriceBannerPickerSelector   string `mapstructure:"live_price_banner_picker_selector"`
}

func getTemplate(size int) *template.Template {
	r.RLock()
	defer r.RUnlock()
	res, ok := dynTemplate[size]
	assert.True(ok, "[BUG] invalid size")

	return res
}

func renderDynamicBanner(w http.ResponseWriter, ctx *builder.Context, slot entity.Seat, ad entity.Advertise) error {
	attr := &dynamicAttribute{}
	err := mapstructure.Decode(attr, ad.Attributes())
	assert.Nil(err)
	if ctx.GetCommon().Scheme == "https" {
		attr.Product = strings.Replace(attr.Product, "http://", "https://", -1)
		attr.Logo = strings.Replace(attr.Logo, "http://", "https://", -1)
	}
	res := getTemplate(slot.Size())
	attr.Link = slot.ClickURL()
	return res.Execute(w, attr)
}

func init() {
	r.Lock()
	defer r.Unlock()
	for i := range dynamic1 {
		dynTemplate[i] = template.Must(template.New(fmt.Sprintf("template_%d", i)).Parse(strings.Replace(dynamic1[i], "<style>", additional, 1)))
	}
}
